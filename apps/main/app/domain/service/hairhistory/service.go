package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/repository"
)

type Service interface {
	// finder_history.go
	ListByUserID(ctx context.Context, userID string) ([]*domain.HairHistory, error)

	// create.go
	Create(ctx context.Context, userID string, req domain.CreateHairHistoryParams) (*domain.HairHistory, error)

	// update.go
	Update(ctx context.Context, historyID string, req domain.UpdateHairHistoryParams) (*domain.HairHistory, error)

	// delete.go
	Delete(ctx context.Context, historyID string) error
}

type service struct {
	hairHistoryRepo repository.HairHistory
}

func NewService(hairHistoryRepo repository.HairHistory) Service {
	return &service{
		hairHistoryRepo: hairHistoryRepo,
	}
}
