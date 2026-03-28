package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/hairhistory"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/user"
	appinfra "github.com/annasakai/hairhistorymemo/apps/main/app/infra"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

// defaultDatabaseURL is used when DATABASE_URL is unset (local docker-compose Postgres).
const defaultDatabaseURL = "postgres://postgres:postgres@127.0.0.1:5432/hairhistory?sslmode=disable"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- PostgreSQL ---
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDatabaseURL
		log.Printf("DATABASE_URL not set; using default %s", defaultDatabaseURL)
	}
	pool, err := appinfra.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("db pool: %v", err)
	}
	defer pool.Close()

	// --- repositories (infra) ---
	userRepo, err := appinfra.NewUserRepositoryPG(pool)
	if err != nil {
		log.Fatalf("user repository: %v", err)
	}
	hairHistoryRepo, err := appinfra.NewHairHistoryRepositoryPG(pool)
	if err != nil {
		log.Fatalf("hair history repository: %v", err)
	}

	// --- domain services ---
	userSvc := user.NewService(userRepo)
	hairHistorySvc := hairhistory.NewService(hairHistoryRepo, userRepo)

	// --- usecases ---
	userUsecase := usecase.NewUser(userSvc)
	hairHistoryUsecase := usecase.NewHairHistory(hairHistorySvc)

	// --- HTTP (chi; ルート定義は controller 配下) ---
	deps := controller.Deps{
		User:        userUsecase,
		HairHistory: hairHistoryUsecase,
		DB:          pool,
	}
	handler := controller.NewRouter(deps)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Printf("api listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
}
