package controller

import "github.com/annasakai/hairhistorymemo/apps/main/app/usecase"

type Deps struct {
	User        *usecase.User
	HairHistory *usecase.HairHistory
}

