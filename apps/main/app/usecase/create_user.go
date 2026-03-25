package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type CreateUserResponse struct {
	UserID string `json:"userId"`
}

type CreateUserUsecase struct {
	userRepo domain.UserRepository
}

func NewCreateUserUsecase(userRepo domain.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{userRepo: userRepo}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, _ *request.CreateUser) (CreateUserResponse, error) {
	user, err := u.userRepo.Create(ctx)
	if err != nil {
		return CreateUserResponse{}, err
	}
	return CreateUserResponse{UserID: user.ID}, nil
}

