package share

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/treatment"
)

var ErrNotFound = errors.New("share: record not found")

const (
	DefaultExpiresInHours = 168
	MaxExpiresInHours     = 24 * 365
	tokenBytes            = 32
)

type Link struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// IsActive reports whether the link may still be redeemed at the given moment.
func (l Link) IsActive(now time.Time) bool {
	if l.RevokedAt != nil {
		return false
	}
	return now.Before(l.ExpiresAt)
}

type Repository interface {
	Create(ctx context.Context, link Link) (Link, error)
	ListActiveByUserID(ctx context.Context, userID uuid.UUID, now time.Time) ([]Link, error)
	FindByToken(ctx context.Context, token string) (Link, error)
	Revoke(ctx context.Context, userID uuid.UUID, token string, revokedAt time.Time) error
}

type TreatmentReader interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]treatment.Treatment, error)
}

func GenerateToken() (string, error) {
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
