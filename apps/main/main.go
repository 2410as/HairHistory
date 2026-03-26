package main

import (
	"log"
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
	"github.com/annasakai/hairhistorymemo/apps/main/app/infra"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

func main() {
	// TODO: PostgreSQL接続・初期化をここに入れる
	userRepo := &infra.UserRepositoryPG{}
	hairHistoryRepo := &infra.HairHistoryRepositoryPG{}

	deps := controller.Deps{
		CreateUser:      usecase.NewCreateUserUsecase(userRepo),
		ListHistories:   usecase.NewListHistoriesUsecase(hairHistoryRepo),
		CreateHistory:   usecase.NewCreateHistoryUsecase(hairHistoryRepo),
		UpdateHistory:   usecase.NewUpdateHistoryUsecase(hairHistoryRepo),
		DeleteHistory:   usecase.NewDeleteHistoryUsecase(hairHistoryRepo),
	}

	handler := controller.NewRouter(deps)

	addr := ":8080"
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

