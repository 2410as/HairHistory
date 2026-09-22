package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

var reference = time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)

func fixedClock(now time.Time) Clock { return func() time.Time { return now } }

type stubRepository struct {
	upserted   UpsertUserParams
	user       User
	userErr    error
	session    Session
	sessionErr error
	created    Session
	deleted    string
}

func (s *stubRepository) UpsertUser(_ context.Context, params UpsertUserParams) (User, error) {
	s.upserted = params
	return User{ID: params.ID, GoogleSub: params.GoogleSub, Email: params.Email, Name: params.Name}, nil
}

func (s *stubRepository) FindUserByID(context.Context, uuid.UUID) (User, error) {
	return s.user, s.userErr
}

func (s *stubRepository) CreateSession(_ context.Context, session Session) error {
	s.created = session
	return nil
}

func (s *stubRepository) FindSessionByTokenHash(context.Context, string) (Session, error) {
	return s.session, s.sessionErr
}

func (s *stubRepository) DeleteSessionByTokenHash(_ context.Context, tokenHash string) error {
	s.deleted = tokenHash
	return nil
}

func (s *stubRepository) DeleteExpiredSessions(context.Context, time.Time) error { return nil }

type stubVerifier struct {
	claims GoogleClaims
	err    error
}

func (s stubVerifier) Verify(context.Context, string) (GoogleClaims, error) {
	return s.claims, s.err
}

func TestLoginWithGoogle(t *testing.T) {
	tests := []struct {
		name      string
		idToken   string
		verifier  TokenVerifier
		wantCode  httpx.ErrorCode
		wantName  string
		wantEmail string
	}{
		{
			name:      "valid token issues a session",
			idToken:   "google-id-token",
			verifier:  stubVerifier{claims: GoogleClaims{Subject: "sub-1", Email: "user@example.com", Name: "Anna"}},
			wantName:  "Anna",
			wantEmail: "user@example.com",
		},
		{
			name:      "missing name falls back to email",
			idToken:   "google-id-token",
			verifier:  stubVerifier{claims: GoogleClaims{Subject: "sub-1", Email: "user@example.com"}},
			wantName:  "user@example.com",
			wantEmail: "user@example.com",
		},
		{
			name:     "empty token is rejected",
			idToken:  "  ",
			verifier: stubVerifier{claims: GoogleClaims{Subject: "sub-1"}},
			wantCode: httpx.CodeInvalidArgument,
		},
		{
			name:     "verification failure is unauthenticated",
			idToken:  "bad-token",
			verifier: stubVerifier{err: errors.New("signature mismatch")},
			wantCode: httpx.CodeUnauthenticated,
		},
		{
			name:     "missing subject is unauthenticated",
			idToken:  "google-id-token",
			verifier: stubVerifier{claims: GoogleClaims{Email: "user@example.com"}},
			wantCode: httpx.CodeUnauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubRepository{}
			uc := NewUsecase(repo, tt.verifier, 720*time.Hour, fixedClock(reference))

			result, err := uc.LoginWithGoogle(context.Background(), tt.idToken)

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
			if result.User.Name != tt.wantName {
				t.Errorf("name = %q, want %q", result.User.Name, tt.wantName)
			}
			if result.User.Email != tt.wantEmail {
				t.Errorf("email = %q, want %q", result.User.Email, tt.wantEmail)
			}
			if !result.SessionExpiresAt.Equal(reference.Add(720 * time.Hour)) {
				t.Errorf("expiresAt = %v, want %v", result.SessionExpiresAt, reference.Add(720*time.Hour))
			}
			if repo.created.TokenHash != HashToken(result.SessionToken) {
				t.Error("stored hash must match the issued token")
			}
			if strings.Contains(repo.created.TokenHash, result.SessionToken) {
				t.Error("raw session token must not be stored")
			}
		})
	}
}

func TestLoginForDevelopment(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		wantSub  string
		wantName string
		wantCode httpx.ErrorCode
	}{
		{
			name:     "creates a namespaced user",
			email:    "dev@example.com",
			userName: "Dev User",
			wantSub:  "dev|dev@example.com",
			wantName: "Dev User",
		},
		{
			name:     "name defaults to email",
			email:    "dev@example.com",
			wantSub:  "dev|dev@example.com",
			wantName: "dev@example.com",
		},
		{
			name:     "email is required",
			email:    "   ",
			wantCode: httpx.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubRepository{}
			uc := NewUsecase(repo, nil, time.Hour, fixedClock(reference))

			result, err := uc.LoginForDevelopment(context.Background(), tt.email, tt.userName)

			if tt.wantCode != "" {
				var apiErr *httpx.APIError
				if !errors.As(err, &apiErr) || apiErr.Code != tt.wantCode {
					t.Fatalf("want %q, got %v", tt.wantCode, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.upserted.GoogleSub != tt.wantSub {
				t.Errorf("googleSub = %q, want %q", repo.upserted.GoogleSub, tt.wantSub)
			}
			if result.User.Name != tt.wantName {
				t.Errorf("name = %q, want %q", result.User.Name, tt.wantName)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		token    string
		repo     *stubRepository
		wantCode httpx.ErrorCode
	}{
		{
			name:  "unexpired session resolves the user",
			token: "session-token",
			repo: &stubRepository{
				session: Session{UserID: userID, ExpiresAt: reference.Add(time.Hour)},
				user:    User{ID: userID},
			},
		},
		{
			name:  "session expiring exactly now is rejected",
			token: "session-token",
			repo: &stubRepository{
				session: Session{UserID: userID, ExpiresAt: reference},
				user:    User{ID: userID},
			},
			wantCode: httpx.CodeUnauthenticated,
		},
		{
			name:  "expired session is rejected",
			token: "session-token",
			repo: &stubRepository{
				session: Session{UserID: userID, ExpiresAt: reference.Add(-time.Second)},
				user:    User{ID: userID},
			},
			wantCode: httpx.CodeUnauthenticated,
		},
		{
			name:     "unknown session is rejected",
			token:    "session-token",
			repo:     &stubRepository{sessionErr: ErrNotFound},
			wantCode: httpx.CodeUnauthenticated,
		},
		{
			name:     "empty token is rejected",
			token:    "",
			repo:     &stubRepository{},
			wantCode: httpx.CodeUnauthenticated,
		},
		{
			name:     "repository failure is internal",
			token:    "session-token",
			repo:     &stubRepository{sessionErr: errors.New("connection reset")},
			wantCode: httpx.CodeInternal,
		},
		{
			name:  "deleted user invalidates the session",
			token: "session-token",
			repo: &stubRepository{
				session: Session{UserID: userID, ExpiresAt: reference.Add(time.Hour)},
				userErr: ErrNotFound,
			},
			wantCode: httpx.CodeUnauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(tt.repo, nil, time.Hour, fixedClock(reference))

			user, err := uc.Authenticate(context.Background(), tt.token)

			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if user.ID != userID {
					t.Errorf("user id = %v, want %v", user.ID, userID)
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

func TestLogoutDeletesHashedToken(t *testing.T) {
	repo := &stubRepository{}
	uc := NewUsecase(repo, nil, time.Hour, fixedClock(reference))

	if err := uc.Logout(context.Background(), "session-token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleted != HashToken("session-token") {
		t.Errorf("deleted = %q, want the hashed token", repo.deleted)
	}
}

func TestGenerateTokenIsUnguessable(t *testing.T) {
	seen := map[string]struct{}{}
	for i := 0; i < 100; i++ {
		token, err := GenerateToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(token) < 40 {
			t.Fatalf("token length = %d, want at least 40", len(token))
		}
		if _, duplicate := seen[token]; duplicate {
			t.Fatal("tokens must not repeat")
		}
		seen[token] = struct{}{}
	}
}
