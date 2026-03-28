package repository

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

// HairHistory defines the persistence contract for HairHistory.
type HairHistory interface {
	ListByUserID(ctx context.Context, userID string) ([]*domain.HairHistory, error)
	Create(ctx context.Context, userID string, req domain.CreateHairHistoryParams) (*domain.HairHistory, error)
	Update(ctx context.Context, historyID string, req domain.UpdateHairHistoryParams) (*domain.HairHistory, error)
	Delete(ctx context.Context, historyID string) error
}
