package usecase

import usersvc "github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/user"

type User struct {
	userSvc usersvc.Service
}

func NewUser(userSvc usersvc.Service) *User {
	return &User{userSvc: userSvc}
}
