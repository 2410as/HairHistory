package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const tokenBytes = 32

type User struct {
	ID        uuid.UUID
	GoogleSub string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UpsertUserParams struct {
	ID        uuid.UUID
	GoogleSub string
	Email     string
	Name      string
}

type GoogleClaims struct {
	Subject string
	Email   string
	Name    string
}

type TokenVerifier interface {
	Verify(ctx context.Context, idToken string) (GoogleClaims, error)
}

type Repository interface {
	UpsertUser(ctx context.Context, params UpsertUserParams) (User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (User, error)
	CreateSession(ctx context.Context, session Session) error
	FindSessionByTokenHash(ctx context.Context, tokenHash string) (Session, error)
	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
	DeleteExpiredSessions(ctx context.Context, now time.Time) error
}

// HashToken derives the value stored in the database. Raw session tokens are
// never persisted.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func GenerateToken() (string, error) {
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
