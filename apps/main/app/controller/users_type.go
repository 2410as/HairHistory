package controller

import "github.com/annasakai/hairhistorymemo/apps/main/app/usecase"

type Users struct {
	userUsecase *usecase.User
}

func NewUsers(userUsecase *usecase.User) *Users {
	return &Users{userUsecase: userUsecase}
}
