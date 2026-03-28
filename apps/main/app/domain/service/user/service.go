package user

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

type Service interface {
	// create.go
	Create(ctx context.Context) (*entity.User, error)

	// finder_user.go
	GetByID(ctx context.Context, userID string) (*entity.User, error)
}

type service struct {
	userRepo domain.UserRepository
}

func NewService(userRepo domain.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}
