package domain

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

// HairHistoryRepository is the persistence contract for HairHistory (implemented in infra).
type HairHistoryRepository interface {
	ListByUserID(ctx context.Context, userID string) ([]*entity.HairHistory, error)
	Create(ctx context.Context, userID string, req entity.CreateHairHistoryParams) (*entity.HairHistory, error)
	Update(ctx context.Context, historyID string, req entity.UpdateHairHistoryParams) (*entity.HairHistory, error)
	Delete(ctx context.Context, historyID string) error
}
