package treatment

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

type stubRepository struct {
	created    CreateParams
	createErr  error
	getResult  Treatment
	getErr     error
	updated    UpdateParams
	updateErr  error
	deletedIDs []uuid.UUID
	deleteErr  error
	listResult []Treatment
	listErr    error
}

func (s *stubRepository) Create(_ context.Context, params CreateParams) (Treatment, error) {
	s.created = params
	if s.createErr != nil {
		return Treatment{}, s.createErr
	}
	return Treatment{
		ID:        params.ID,
		UserID:    params.UserID,
		TreatedOn: params.TreatedOn,
		Services:  params.Services,
		SalonName: params.SalonName,
		Memo:      params.Memo,
		Cost:      params.Cost,
	}, nil
}

func (s *stubRepository) ListByUserID(context.Context, uuid.UUID) ([]Treatment, error) {
	return s.listResult, s.listErr
}

func (s *stubRepository) GetByID(context.Context, uuid.UUID) (Treatment, error) {
	return s.getResult, s.getErr
}

func (s *stubRepository) Update(_ context.Context, params UpdateParams) (Treatment, error) {
	s.updated = params
	if s.updateErr != nil {
		return Treatment{}, s.updateErr
	}
	return Treatment{ID: params.ID, TreatedOn: params.TreatedOn, Services: params.Services}, nil
}

func (s *stubRepository) Delete(_ context.Context, id uuid.UUID) error {
	s.deletedIDs = append(s.deletedIDs, id)
	return s.deleteErr
}

func ptr[T any](v T) *T { return &v }

func longString(n int) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = 'あ'
	}
	return string(runes)
}

func TestUsecaseCreateValidation(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name         string
		input        Input
		wantCode     httpx.ErrorCode
		wantDetail   string
		wantServices []string
	}{
		{
			name:         "minimal valid input",
			input:        Input{TreatedOn: "2026-09-10", Services: []string{"カラー"}},
			wantServices: []string{"カラー"},
		},
		{
			name: "full valid input",
			input: Input{
				TreatedOn: "2026-09-10",
				Services:  []string{"カット", "カラー"},
				SalonName: ptr("Hair Salon ABC"),
				Memo:      ptr("アッシュグレー"),
				Cost:      ptr(8000),
			},
			wantServices: []string{"カット", "カラー"},
		},
		{
			name:         "future date is allowed",
			input:        Input{TreatedOn: "2099-12-31", Services: []string{"パーマ"}},
			wantServices: []string{"パーマ"},
		},
		{
			name:         "services are trimmed",
			input:        Input{TreatedOn: "2026-09-10", Services: []string{"  カット  "}},
			wantServices: []string{"カット"},
		},
		{
			name:         "cost zero is allowed",
			input:        Input{TreatedOn: "2026-09-10", Services: []string{"カット"}, Cost: ptr(0)},
			wantServices: []string{"カット"},
		},
		{
			name:       "missing treatedOn",
			input:      Input{Services: []string{"カット"}},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "treatedOn",
		},
		{
			name:       "treatedOn with wrong format",
			input:      Input{TreatedOn: "2026/09/10", Services: []string{"カット"}},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "treatedOn",
		},
		{
			name:       "treatedOn carrying a time component",
			input:      Input{TreatedOn: "2026-09-10T00:00:00Z", Services: []string{"カット"}},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "treatedOn",
		},
		{
			name:       "empty services",
			input:      Input{TreatedOn: "2026-09-10", Services: []string{}},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "services",
		},
		{
			name:       "nil services",
			input:      Input{TreatedOn: "2026-09-10"},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "services",
		},
		{
			name:       "blank service item",
			input:      Input{TreatedOn: "2026-09-10", Services: []string{"カット", "   "}},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "services",
		},
		{
			name:       "too many services",
			input:      Input{TreatedOn: "2026-09-10", Services: make([]string, 11)},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "services",
		},
		{
			name:       "salonName too long",
			input:      Input{TreatedOn: "2026-09-10", Services: []string{"カット"}, SalonName: ptr(longString(101))},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "salonName",
		},
		{
			name:       "memo too long",
			input:      Input{TreatedOn: "2026-09-10", Services: []string{"カット"}, Memo: ptr(longString(1001))},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "memo",
		},
		{
			name:       "negative cost",
			input:      Input{TreatedOn: "2026-09-10", Services: []string{"カット"}, Cost: ptr(-1)},
			wantCode:   httpx.CodeInvalidArgument,
			wantDetail: "cost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubRepository{}
			uc := NewUsecase(repo)

			created, err := uc.Create(context.Background(), userID, tt.input)

			if tt.wantCode != "" {
				var apiErr *httpx.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("want APIError, got %v", err)
				}
				if apiErr.Code != tt.wantCode {
					t.Errorf("code = %q, want %q", apiErr.Code, tt.wantCode)
				}
				if _, ok := apiErr.Details[tt.wantDetail]; !ok {
					t.Errorf("details = %v, want key %q", apiErr.Details, tt.wantDetail)
				}
				if repo.created.ID != uuid.Nil {
					t.Error("repository must not be called for invalid input")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if created.UserID != userID {
				t.Errorf("userID = %v, want %v", created.UserID, userID)
			}
			if got, want := len(created.Services), len(tt.wantServices); got != want {
				t.Fatalf("len(services) = %d, want %d", got, want)
			}
			for i, service := range tt.wantServices {
				if created.Services[i] != service {
					t.Errorf("services[%d] = %q, want %q", i, created.Services[i], service)
				}
			}
			if created.ID == uuid.Nil {
				t.Error("id must be generated")
			}
		})
	}
}

func TestUsecaseCreateNormalizesBlankOptionalFieldsToNil(t *testing.T) {
	repo := &stubRepository{}
	uc := NewUsecase(repo)

	if _, err := uc.Create(context.Background(), uuid.New(), Input{
		TreatedOn: "2026-09-10",
		Services:  []string{"カット"},
		SalonName: ptr("   "),
		Memo:      ptr(""),
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.created.SalonName != nil {
		t.Errorf("salonName = %v, want nil", *repo.created.SalonName)
	}
	if repo.created.Memo != nil {
		t.Errorf("memo = %v, want nil", *repo.created.Memo)
	}
}

func TestUsecaseOwnershipChecks(t *testing.T) {
	owner := uuid.New()
	stranger := uuid.New()
	treatmentID := uuid.New()

	tests := []struct {
		name      string
		repo      *stubRepository
		requester uuid.UUID
		wantCode  httpx.ErrorCode
	}{
		{
			name:      "owner may read",
			repo:      &stubRepository{getResult: Treatment{ID: treatmentID, UserID: owner}},
			requester: owner,
		},
		{
			name:      "stranger is denied",
			repo:      &stubRepository{getResult: Treatment{ID: treatmentID, UserID: owner}},
			requester: stranger,
			wantCode:  httpx.CodePermissionDenied,
		},
		{
			name:      "missing row is not found",
			repo:      &stubRepository{getErr: ErrNotFound},
			requester: owner,
			wantCode:  httpx.CodeNotFound,
		},
		{
			name:      "repository failure is internal",
			repo:      &stubRepository{getErr: errors.New("connection reset")},
			requester: owner,
			wantCode:  httpx.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(tt.repo)

			_, err := uc.Get(context.Background(), tt.requester, treatmentID)

			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
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
			if apiErr.Code == httpx.CodeInternal && apiErr.Message != "internal server error" {
				t.Errorf("internal errors must not leak details, got %q", apiErr.Message)
			}
		})
	}
}

func TestUsecaseDeleteRejectsStranger(t *testing.T) {
	owner := uuid.New()
	repo := &stubRepository{getResult: Treatment{ID: uuid.New(), UserID: owner}}
	uc := NewUsecase(repo)

	err := uc.Delete(context.Background(), uuid.New(), repo.getResult.ID)

	var apiErr *httpx.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != httpx.CodePermissionDenied {
		t.Fatalf("want permission_denied, got %v", err)
	}
	if len(repo.deletedIDs) != 0 {
		t.Error("delete must not reach the repository")
	}
}
