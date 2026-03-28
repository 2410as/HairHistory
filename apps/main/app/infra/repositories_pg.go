package infra

import (
	"context"
	"errors"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/repository"
)

// UserRepositoryPG is the PostgreSQL implementation.
// NOTE: Stub only for scaffolding.
type UserRepositoryPG struct{}

var _ repository.User = (*UserRepositoryPG)(nil)

func (r *UserRepositoryPG) Create(ctx context.Context) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepositoryPG) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

// HairHistoryRepositoryPG is the PostgreSQL implementation.
// NOTE: Stub only for scaffolding.
type HairHistoryRepositoryPG struct{}

var _ repository.HairHistory = (*HairHistoryRepositoryPG)(nil)

func (r *HairHistoryRepositoryPG) ListByUserID(ctx context.Context, userID string) ([]*domain.HairHistory, error) {
	return nil, errors.New("not implemented")
}

func (r *HairHistoryRepositoryPG) Create(ctx context.Context, userID string, req domain.CreateHairHistoryParams) (*domain.HairHistory, error) {
	return nil, errors.New("not implemented")
}

func (r *HairHistoryRepositoryPG) Update(ctx context.Context, historyID string, req domain.UpdateHairHistoryParams) (*domain.HairHistory, error) {
	return nil, errors.New("not implemented")
}

func (r *HairHistoryRepositoryPG) Delete(ctx context.Context, historyID string) error {
	return errors.New("not implemented")
}
