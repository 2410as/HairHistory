package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

type CreateUserRequest struct{}

type CreateUserResponse struct {
	UserID string `json:"userId"`
}

type CreateUserUsecase struct {
	userRepo domain.UserRepository
}

func NewCreateUserUsecase(userRepo domain.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{userRepo: userRepo}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	user, err := u.userRepo.Create(ctx)
	if err != nil {
		return CreateUserResponse{}, err
	}
	return CreateUserResponse{UserID: user.ID}, nil
}

