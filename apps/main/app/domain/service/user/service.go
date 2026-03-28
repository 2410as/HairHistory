package user

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/repository"
)

type Service interface {
	// create.go
	Create(ctx context.Context) (*domain.User, error)

	// finder_user.go
	GetByID(ctx context.Context, userID string) (*domain.User, error)
}

type service struct {
	userRepo repository.User
}

func NewService(userRepo repository.User) Service {
	return &service{
		userRepo: userRepo,
	}
}
