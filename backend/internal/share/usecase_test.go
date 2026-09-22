package share

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
	"github.com/2410as/HairHistory/backend/internal/treatment"
)

var reference = time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)

func fixedClock(now time.Time) Clock { return func() time.Time { return now } }

type stubRepository struct {
	created   Link
	findLink  Link
	findErr   error
	revokeErr error
	revoked   bool
	listLinks []Link
	listErr   error
	listNow   time.Time
}

func (s *stubRepository) Create(_ context.Context, link Link) (Link, error) {
	s.created = link
	return link, nil
}

func (s *stubRepository) ListActiveByUserID(_ context.Context, _ uuid.UUID, now time.Time) ([]Link, error) {
	s.listNow = now
	return s.listLinks, s.listErr
}

func (s *stubRepository) FindByToken(context.Context, string) (Link, error) {
	return s.findLink, s.findErr
}

func (s *stubRepository) Revoke(context.Context, uuid.UUID, string, time.Time) error {
	if s.revokeErr != nil {
		return s.revokeErr
	}
	s.revoked = true
	return nil
}

type stubTreatments struct {
	items      []treatment.Treatment
	err        error
	calledWith uuid.UUID
}

func (s *stubTreatments) ListByUserID(_ context.Context, userID uuid.UUID) ([]treatment.Treatment, error) {
	s.calledWith = userID
	return s.items, s.err
}

func ptrTime(v time.Time) *time.Time { return &v }

func TestLinkIsActive(t *testing.T) {
	tests := []struct {
		name string
		link Link
		now  time.Time
		want bool
	}{
		{
			name: "expires in the future and not revoked",
			link: Link{ExpiresAt: reference.Add(time.Hour)},
			now:  reference,
			want: true,
		},
		{
			name: "one second before expiry",
			link: Link{ExpiresAt: reference.Add(time.Second)},
			now:  reference,
			want: true,
		},
		{
			name: "exactly at expiry is inactive",
			link: Link{ExpiresAt: reference},
			now:  reference,
			want: false,
		},
		{
			name: "expired one second ago",
			link: Link{ExpiresAt: reference.Add(-time.Second)},
			now:  reference,
			want: false,
		},
		{
			name: "revoked before expiry",
			link: Link{ExpiresAt: reference.Add(time.Hour), RevokedAt: ptrTime(reference.Add(-time.Minute))},
			now:  reference,
			want: false,
		},
		{
			name: "revoked and expired",
			link: Link{ExpiresAt: reference.Add(-time.Hour), RevokedAt: ptrTime(reference.Add(-time.Hour))},
			now:  reference,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.link.IsActive(tt.now); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseCreateExpiry(t *testing.T) {
	tests := []struct {
		name           string
		expiresInHours *int
		wantExpiresAt  time.Time
		wantCode       httpx.ErrorCode
	}{
		{
			name:          "defaults to seven days",
			wantExpiresAt: reference.Add(168 * time.Hour),
		},
		{
			name:           "explicit 24 hours",
			expiresInHours: intPtr(24),
			wantExpiresAt:  reference.Add(24 * time.Hour),
		},
		{
			name:           "maximum is accepted",
			expiresInHours: intPtr(MaxExpiresInHours),
			wantExpiresAt:  reference.Add(MaxExpiresInHours * time.Hour),
		},
		{
			name:           "zero is rejected",
			expiresInHours: intPtr(0),
			wantCode:       httpx.CodeInvalidArgument,
		},
		{
			name:           "negative is rejected",
			expiresInHours: intPtr(-1),
			wantCode:       httpx.CodeInvalidArgument,
		},
		{
			name:           "above the maximum is rejected",
			expiresInHours: intPtr(MaxExpiresInHours + 1),
			wantCode:       httpx.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubRepository{}
			uc := NewUsecase(repo, &stubTreatments{}, fixedClock(reference))

			created, err := uc.Create(context.Background(), uuid.New(), tt.expiresInHours)

			if tt.wantCode != "" {
				var apiErr *httpx.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("want APIError, got %v", err)
				}
				if apiErr.Code != tt.wantCode {
					t.Errorf("code = %q, want %q", apiErr.Code, tt.wantCode)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !created.ExpiresAt.Equal(tt.wantExpiresAt) {
				t.Errorf("expiresAt = %v, want %v", created.ExpiresAt, tt.wantExpiresAt)
			}
			if len(created.Token) < 40 {
				t.Errorf("token length = %d, want an unguessable token", len(created.Token))
			}
			if created.Token != repo.created.Token {
				t.Error("stored token must match the returned token")
			}
		})
	}
}

func TestUsecaseCreateGeneratesDistinctTokens(t *testing.T) {
	uc := NewUsecase(&stubRepository{}, &stubTreatments{}, fixedClock(reference))
	userID := uuid.New()

	seen := map[string]struct{}{}
	for i := 0; i < 100; i++ {
		link, err := uc.Create(context.Background(), userID, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, duplicate := seen[link.Token]; duplicate {
			t.Fatal("share tokens must not repeat")
		}
		seen[link.Token] = struct{}{}
	}
}

func TestUsecasePublicView(t *testing.T) {
	ownerID := uuid.New()
	items := []treatment.Treatment{{ID: uuid.New(), UserID: ownerID, Services: []string{"カラー"}}}

	tests := []struct {
		name      string
		token     string
		repo      *stubRepository
		wantCode  httpx.ErrorCode
		wantItems int
	}{
		{
			name:      "active link returns treatments",
			token:     "valid-token",
			repo:      &stubRepository{findLink: Link{UserID: ownerID, ExpiresAt: reference.Add(time.Hour)}},
			wantItems: 1,
		},
		{
			name:     "expired link is hidden",
			token:    "expired-token",
			repo:     &stubRepository{findLink: Link{UserID: ownerID, ExpiresAt: reference.Add(-time.Second)}},
			wantCode: httpx.CodeNotFound,
		},
		{
			name:  "revoked link is hidden",
			token: "revoked-token",
			repo: &stubRepository{findLink: Link{
				UserID:    ownerID,
				ExpiresAt: reference.Add(time.Hour),
				RevokedAt: ptrTime(reference.Add(-time.Minute)),
			}},
			wantCode: httpx.CodeNotFound,
		},
		{
			name:     "unknown token is not found",
			token:    "missing-token",
			repo:     &stubRepository{findErr: ErrNotFound},
			wantCode: httpx.CodeNotFound,
		},
		{
			name:     "empty token is not found",
			token:    "",
			repo:     &stubRepository{findLink: Link{UserID: ownerID, ExpiresAt: reference.Add(time.Hour)}},
			wantCode: httpx.CodeNotFound,
		},
		{
			name:     "repository failure is internal",
			token:    "valid-token",
			repo:     &stubRepository{findErr: errors.New("connection reset")},
			wantCode: httpx.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treatments := &stubTreatments{items: items}
			uc := NewUsecase(tt.repo, treatments, fixedClock(reference))

			got, err := uc.PublicView(context.Background(), tt.token)

			if tt.wantCode != "" {
				var apiErr *httpx.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("want APIError, got %v", err)
				}
				if apiErr.Code != tt.wantCode {
					t.Errorf("code = %q, want %q", apiErr.Code, tt.wantCode)
				}
				if treatments.calledWith != uuid.Nil {
					t.Error("treatments must not be read for an unusable link")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantItems {
				t.Errorf("len(treatments) = %d, want %d", len(got), tt.wantItems)
			}
			if treatments.calledWith != ownerID {
				t.Errorf("treatments read for %v, want %v", treatments.calledWith, ownerID)
			}
		})
	}
}

func TestUsecaseRevoke(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		repo     *stubRepository
		wantCode httpx.ErrorCode
	}{
		{
			name:  "revokes an owned link",
			token: "valid-token",
			repo:  &stubRepository{},
		},
		{
			name:     "unknown or foreign link is not found",
			token:    "valid-token",
			repo:     &stubRepository{revokeErr: ErrNotFound},
			wantCode: httpx.CodeNotFound,
		},
		{
			name:     "blank token is rejected",
			token:    "  ",
			repo:     &stubRepository{},
			wantCode: httpx.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(tt.repo, &stubTreatments{}, fixedClock(reference))

			err := uc.Revoke(context.Background(), uuid.New(), tt.token)

			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !tt.repo.revoked {
					t.Error("repository revoke was not called")
				}
				return
			}

			var apiErr *httpx.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("want APIError, got %v", err)
			}
			if apiErr.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", apiErr.Code, tt.wantCode)
			}
		})
	}
}

func TestUsecaseListUsesCurrentTime(t *testing.T) {
	repo := &stubRepository{}
	uc := NewUsecase(repo, &stubTreatments{}, fixedClock(reference))

	if _, err := uc.List(context.Background(), uuid.New()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.listNow.Equal(reference) {
		t.Errorf("listNow = %v, want %v", repo.listNow, reference)
	}
}

func intPtr(v int) *int { return &v }
