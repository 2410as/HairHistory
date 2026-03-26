package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

type User struct {
	userRepo domain.UserRepository
}

func NewUser(userRepo domain.UserRepository) *User {
	return &User{userRepo: userRepo}
}

func (u User) Create(ctx context.Context, _ *request.CreateUser) (*response.CreateUser, error) {
	user, err := u.userRepo.Create(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewCreateUser(user.ID), nil
}

