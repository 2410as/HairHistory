package treatment

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Create(ctx context.Context, userID uuid.UUID, in Input) (Treatment, error) {
	valid, err := validate(in)
	if err != nil {
		return Treatment{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Treatment{}, httpx.Internal(err)
	}

	created, err := u.repo.Create(ctx, CreateParams{
		ID:        id,
		UserID:    userID,
		TreatedOn: valid.treatedOn,
		Services:  valid.services,
		SalonName: valid.salonName,
		Memo:      valid.memo,
		Cost:      valid.cost,
	})
	if err != nil {
		return Treatment{}, httpx.Internal(err)
	}
	return created, nil
}

func (u *Usecase) List(ctx context.Context, userID uuid.UUID) ([]Treatment, error) {
	treatments, err := u.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, httpx.Internal(err)
	}
	return treatments, nil
}

func (u *Usecase) Get(ctx context.Context, userID, id uuid.UUID) (Treatment, error) {
	return u.findOwned(ctx, userID, id)
}

func (u *Usecase) Update(ctx context.Context, userID, id uuid.UUID, in Input) (Treatment, error) {
	valid, err := validate(in)
	if err != nil {
		return Treatment{}, err
	}

	if _, err := u.findOwned(ctx, userID, id); err != nil {
		return Treatment{}, err
	}

	updated, err := u.repo.Update(ctx, UpdateParams{
		ID:        id,
		TreatedOn: valid.treatedOn,
		Services:  valid.services,
		SalonName: valid.salonName,
		Memo:      valid.memo,
		Cost:      valid.cost,
	})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Treatment{}, httpx.NotFound("treatment not found")
		}
		return Treatment{}, httpx.Internal(err)
	}
	return updated, nil
}

func (u *Usecase) Delete(ctx context.Context, userID, id uuid.UUID) error {
	if _, err := u.findOwned(ctx, userID, id); err != nil {
		return err
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return httpx.NotFound("treatment not found")
		}
		return httpx.Internal(err)
	}
	return nil
}

func (u *Usecase) findOwned(ctx context.Context, userID, id uuid.UUID) (Treatment, error) {
	found, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Treatment{}, httpx.NotFound("treatment not found")
		}
		return Treatment{}, httpx.Internal(err)
	}
	if found.UserID != userID {
		return Treatment{}, httpx.PermissionDenied("you do not have access to this treatment")
	}
	return found, nil
}
