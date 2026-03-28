package main

import (
	"log"
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/hairhistory"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/user"
	"github.com/annasakai/hairhistorymemo/apps/main/app/infra"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

func main() {
	// TODO: PostgreSQL接続・初期化をここに入れる
	userRepo := &infra.UserRepositoryPG{}
	hairHistoryRepo := &infra.HairHistoryRepositoryPG{}

	userSvc := user.NewService(userRepo)
	hairHistorySvc := hairhistory.NewService(hairHistoryRepo)

	deps := controller.Deps{
		User:        usecase.NewUser(userSvc),
		HairHistory: usecase.NewHairHistory(hairHistorySvc),
	}

	handler := controller.NewRouter(deps)

	addr := ":8080"
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
