package share

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
	"github.com/2410as/HairHistory/backend/internal/treatment"
)

type Clock func() time.Time

type Usecase struct {
	repo       Repository
	treatments TreatmentReader
	now        Clock
}

func NewUsecase(repo Repository, treatments TreatmentReader, now Clock) *Usecase {
	if now == nil {
		now = time.Now
	}
	return &Usecase{repo: repo, treatments: treatments, now: now}
}

func (u *Usecase) Create(ctx context.Context, userID uuid.UUID, expiresInHours *int) (Link, error) {
	hours := DefaultExpiresInHours
	if expiresInHours != nil {
		hours = *expiresInHours
	}
	if hours <= 0 || hours > MaxExpiresInHours {
		return Link{}, httpx.InvalidArgument("request contains invalid fields", map[string]string{
			"expiresInHours": "must be between 1 and " + strconv.Itoa(MaxExpiresInHours),
		})
	}

	token, err := GenerateToken()
	if err != nil {
		return Link{}, httpx.Internal(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Link{}, httpx.Internal(err)
	}

	now := u.now()
	created, err := u.repo.Create(ctx, Link{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(time.Duration(hours) * time.Hour),
		CreatedAt: now,
	})
	if err != nil {
		return Link{}, httpx.Internal(err)
	}
	return created, nil
}

func (u *Usecase) List(ctx context.Context, userID uuid.UUID) ([]Link, error) {
	links, err := u.repo.ListActiveByUserID(ctx, userID, u.now())
	if err != nil {
		return nil, httpx.Internal(err)
	}
	return links, nil
}

func (u *Usecase) Revoke(ctx context.Context, userID uuid.UUID, token string) error {
	if strings.TrimSpace(token) == "" {
		return httpx.InvalidArgument("token is required", map[string]string{"token": "required"})
	}

	err := u.repo.Revoke(ctx, userID, token, u.now())
	if errors.Is(err, ErrNotFound) {
		return httpx.NotFound("share link not found")
	}
	if err != nil {
		return httpx.Internal(err)
	}
	return nil
}

// PublicView resolves a share token for an anonymous visitor. It returns only
// treatments; the owner's identity is never part of the result.
func (u *Usecase) PublicView(ctx context.Context, token string) ([]treatment.Treatment, error) {
	if strings.TrimSpace(token) == "" {
		return nil, httpx.NotFound("share link not found")
	}

	link, err := u.repo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, httpx.NotFound("share link not found")
		}
		return nil, httpx.Internal(err)
	}

	if !link.IsActive(u.now()) {
		return nil, httpx.NotFound("share link not found")
	}

	treatments, err := u.treatments.ListByUserID(ctx, link.UserID)
	if err != nil {
		return nil, httpx.Internal(err)
	}
	return treatments, nil
}
