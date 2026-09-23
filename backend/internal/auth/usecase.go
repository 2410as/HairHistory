package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

var ErrNotFound = errors.New("auth: record not found")

const devGoogleSubPrefix = "dev|"

type Clock func() time.Time

type Usecase struct {
	repo       Repository
	verifier   TokenVerifier
	sessionTTL time.Duration
	now        Clock
}

func NewUsecase(repo Repository, verifier TokenVerifier, sessionTTL time.Duration, now Clock) *Usecase {
	if now == nil {
		now = time.Now
	}
	return &Usecase{repo: repo, verifier: verifier, sessionTTL: sessionTTL, now: now}
}

type LoginResult struct {
	User             User
	SessionToken     string
	SessionExpiresAt time.Time
}

func (u *Usecase) LoginWithGoogle(ctx context.Context, idToken string) (LoginResult, error) {
	if strings.TrimSpace(idToken) == "" {
		return LoginResult{}, httpx.InvalidArgument("idToken is required", map[string]string{"idToken": "required"})
	}
	if u.verifier == nil {
		return LoginResult{}, httpx.Internal(errors.New("google token verifier is not configured"))
	}

	claims, err := u.verifier.Verify(ctx, idToken)
	if err != nil {
		return LoginResult{}, httpx.Unauthenticated("google id token is invalid")
	}
	if claims.Subject == "" {
		return LoginResult{}, httpx.Unauthenticated("google id token is missing a subject")
	}

	name := claims.Name
	if strings.TrimSpace(name) == "" {
		name = claims.Email
	}

	return u.issueSession(ctx, UpsertUserParams{GoogleSub: claims.Subject, Email: claims.Email, Name: name})
}

func (u *Usecase) LoginForDevelopment(ctx context.Context, email, name string) (LoginResult, error) {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)

	if email == "" {
		return LoginResult{}, httpx.InvalidArgument("email is required", map[string]string{"email": "required"})
	}
	if name == "" {
		name = email
	}

	return u.issueSession(ctx, UpsertUserParams{GoogleSub: devGoogleSubPrefix + email, Email: email, Name: name})
}

func (u *Usecase) issueSession(ctx context.Context, params UpsertUserParams) (LoginResult, error) {
	userID, err := uuid.NewV7()
	if err != nil {
		return LoginResult{}, httpx.Internal(err)
	}
	params.ID = userID

	user, err := u.repo.UpsertUser(ctx, params)
	if err != nil {
		return LoginResult{}, httpx.Internal(err)
	}

	token, err := GenerateToken()
	if err != nil {
		return LoginResult{}, httpx.Internal(err)
	}

	sessionID, err := uuid.NewV7()
	if err != nil {
		return LoginResult{}, httpx.Internal(err)
	}

	now := u.now()
	session := Session{
		ID:        sessionID,
		UserID:    user.ID,
		TokenHash: HashToken(token),
		ExpiresAt: now.Add(u.sessionTTL),
		CreatedAt: now,
	}
	if err := u.repo.CreateSession(ctx, session); err != nil {
		return LoginResult{}, httpx.Internal(err)
	}

	return LoginResult{User: user, SessionToken: token, SessionExpiresAt: session.ExpiresAt}, nil
}

func (u *Usecase) Authenticate(ctx context.Context, token string) (User, error) {
	if strings.TrimSpace(token) == "" {
		return User{}, httpx.Unauthenticated("authentication required")
	}

	session, err := u.repo.FindSessionByTokenHash(ctx, HashToken(token))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, httpx.Unauthenticated("session is invalid or expired")
		}
		return User{}, httpx.Internal(err)
	}

	if !u.now().Before(session.ExpiresAt) {
		return User{}, httpx.Unauthenticated("session is invalid or expired")
	}

	user, err := u.repo.FindUserByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, httpx.Unauthenticated("session is invalid or expired")
		}
		return User{}, httpx.Internal(err)
	}

	return user, nil
}

func (u *Usecase) Logout(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return nil
	}
	if err := u.repo.DeleteSessionByTokenHash(ctx, HashToken(token)); err != nil && !errors.Is(err, ErrNotFound) {
		return httpx.Internal(err)
	}
	return nil
}

func (u *Usecase) PurgeExpiredSessions(ctx context.Context) error {
	if err := u.repo.DeleteExpiredSessions(ctx, u.now()); err != nil {
		return httpx.Internal(err)
	}
	return nil
}
