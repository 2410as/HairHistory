package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

type Service interface {
	// finder_history.go
	ListByUserID(ctx context.Context, userID string) ([]*entity.HairHistory, error)

	// create.go
	Create(ctx context.Context, userID string, req entity.CreateHairHistoryParams) (*entity.HairHistory, error)

	// update.go
	Update(ctx context.Context, historyID string, req entity.UpdateHairHistoryParams) (*entity.HairHistory, error)

	// delete.go
	Delete(ctx context.Context, historyID string) error
}

type service struct {
	hairHistoryRepo domain.HairHistoryRepository
	userRepo        domain.UserRepository
}

func NewService(hairHistoryRepo domain.HairHistoryRepository, userRepo domain.UserRepository) Service {
	return &service{
		hairHistoryRepo: hairHistoryRepo,
		userRepo:        userRepo,
	}
}
