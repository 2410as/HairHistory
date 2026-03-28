package infra

import (
	"context"
	"errors"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

// HairHistoryRepositoryPG is the PostgreSQL implementation.
// NOTE: Stub only for scaffolding.
type HairHistoryRepositoryPG struct{}

var _ domain.HairHistoryRepository = (*HairHistoryRepositoryPG)(nil)

func (r *HairHistoryRepositoryPG) ListByUserID(ctx context.Context, userID string) ([]*entity.HairHistory, error) {
	return nil, errors.New("not implemented")
}

func (r *HairHistoryRepositoryPG) Create(ctx context.Context, userID string, req entity.CreateHairHistoryParams) (*entity.HairHistory, error) {
	return nil, errors.New("not implemented")
}

func (r *HairHistoryRepositoryPG) Update(ctx context.Context, historyID string, req entity.UpdateHairHistoryParams) (*entity.HairHistory, error) {
	return nil, errors.New("not implemented")
}

func (r *HairHistoryRepositoryPG) Delete(ctx context.Context, historyID string) error {
	return errors.New("not implemented")
}
