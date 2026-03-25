package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
	"github.com/annasakai/hairhistorymemo/apps/main/app/infra"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

func main() {
	// TODO: PostgreSQL接続・初期化をここに入れる
	userRepo := &infra.UserRepositoryPG{}
	hairHistoryRepo := &infra.HairHistoryRepositoryPG{}

	createUserUsecase := usecase.NewCreateUserUsecase(userRepo)
	listHistoriesUsecase := usecase.NewListHistoriesUsecase(hairHistoryRepo)
	createHistoryUsecase := usecase.NewCreateHistoryUsecase(hairHistoryRepo)
	updateHistoryUsecase := usecase.NewUpdateHistoryUsecase(hairHistoryRepo)
	deleteHistoryUsecase := usecase.NewDeleteHistoryUsecase(hairHistoryRepo)

	deps := controller.Deps{
		CreateUser:      createUserUsecase,
		ListHistories:   listHistoriesUsecase,
		CreateHistory:   createHistoryUsecase,
		UpdateHistory:   updateHistoryUsecase,
		DeleteHistory:   deleteHistoryUsecase,
	}

	healthController := controller.NewHealth()
	usersController := controller.NewUsers(deps.CreateUser)
	historiesController := controller.NewHistories(
		deps.ListHistories,
		deps.CreateHistory,
		deps.UpdateHistory,
		deps.DeleteHistory,
	)

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthController.Get)

		r.Post("/users", usersController.Create)

		r.Route("/users/{userId}/histories", func(r chi.Router) {
			r.Get("/", historiesController.List)
			r.Post("/", historiesController.Create)
		})

		r.Route("/histories/{historyId}", func(r chi.Router) {
			r.Put("/", historiesController.Update)
			r.Delete("/", historiesController.Delete)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

