package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

func (u User) Create(ctx context.Context, _ *request.CreateUser) (*response.CreateUser, error) {
	user, err := u.userSvc.Create(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewCreateUser(user.ID), nil
}
