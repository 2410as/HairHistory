package main

import (
	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/hairhistory"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/user"
	"github.com/annasakai/hairhistorymemo/apps/main/app/infra"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

// wireDeps builds repositories → domain services → usecases for the HTTP layer.
func wireDeps() controller.Deps {
	// TODO: PostgreSQL接続・初期化をここに入れる
	userRepo := &infra.UserRepositoryPG{}
	hairHistoryRepo := &infra.HairHistoryRepositoryPG{}

	userSvc := user.NewService(userRepo)
	hairHistorySvc := hairhistory.NewService(hairHistoryRepo)

	return controller.Deps{
		User:        usecase.NewUser(userSvc),
		HairHistory: usecase.NewHairHistory(hairHistorySvc),
	}
}
