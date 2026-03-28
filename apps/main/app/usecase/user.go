package usecase

import (
	"context"

	usersvc "github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/user"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

type User struct {
	userSvc usersvc.Service
}

func NewUser(userSvc usersvc.Service) *User {
	return &User{userSvc: userSvc}
}

func (u User) Create(ctx context.Context, _ *request.CreateUser) (*response.CreateUser, error) {
	user, err := u.userSvc.Create(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewCreateUser(user.ID), nil
}
