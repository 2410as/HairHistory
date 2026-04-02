package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/hairhistory"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/user"
	appinfra "github.com/annasakai/hairhistorymemo/apps/main/app/infra"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
	"github.com/annasakai/hairhistorymemo/apps/main/app/utility/uchi"
)

// defaultDatabaseURL is used when DATABASE_URL is unset (local docker-compose Postgres).
const defaultDatabaseURL = "postgres://postgres:postgres@127.0.0.1:5432/hairhistory?sslmode=disable"
const uuidRoutePattern = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"

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

	// --- controllers ---
	healthController := controller.NewHealth(pool)
	usersController := controller.NewUsers(userUsecase)
	hairHistoryController := controller.NewHairHistory(hairHistoryUsecase)

	// --- HTTP router ---
	r := uchi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Endpoint definitions live here (階層的に r.Route / r.Group を使用).
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthController.Get)

		r.Group(func(r chi.Router) {
			// users/{userId}/histories
			r.Route("/users/{userId:"+uuidRoutePattern+"}", func(r chi.Router) {
				r.Get("/histories", hairHistoryController.List)
				r.Post("/histories", hairHistoryController.Create)
			})

			// /users
			r.Post("/users", usersController.Create)

			// histories/{historyId}
			r.Route("/histories", func(r chi.Router) {
				r.Put("/{historyId:"+uuidRoutePattern+"}", hairHistoryController.Update)
				r.Delete("/{historyId:"+uuidRoutePattern+"}", hairHistoryController.Delete)
			})
		})
	})

	var handler http.Handler = r

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

func corsAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("HAIR_CORS_ORIGINS"))
	if raw == "" {
		return []string{"http://localhost:3000"}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:3000"}
	}
	return out
}
