package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/2410as/HairHistory/backend/internal/auth"
	"github.com/2410as/HairHistory/backend/internal/config"
	"github.com/2410as/HairHistory/backend/internal/db"
	"github.com/2410as/HairHistory/backend/internal/httpx"
	"github.com/2410as/HairHistory/backend/internal/share"
	"github.com/2410as/HairHistory/backend/internal/treatment"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	if err := run(); err != nil {
		slog.Error("server terminated", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	authRepo := auth.NewPostgresRepository(pool)
	treatmentRepo := treatment.NewPostgresRepository(pool)
	shareRepo := share.NewPostgresRepository(pool)

	authUsecase := auth.NewUsecase(authRepo, auth.NewGoogleVerifier(cfg.GoogleClientID), cfg.SessionTTL, time.Now)
	treatmentUsecase := treatment.NewUsecase(treatmentRepo)
	shareUsecase := share.NewUsecase(shareRepo, treatmentRepo, time.Now)

	authHandler := auth.NewHandler(authUsecase, cfg.CookieSecure(), !cfg.IsProduction())
	treatmentHandler := treatment.NewHandler(treatmentUsecase)
	shareHandler := share.NewHandler(shareUsecase)

	if err := authUsecase.PurgeExpiredSessions(ctx); err != nil {
		slog.Warn("purge expired sessions failed", "error", err)
	}

	router := chi.NewRouter()
	router.Use(httpx.RequestID, httpx.Recoverer, httpx.Logging, httpx.CORS(cfg.CORSAllowedOrigins))

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Route("/api", func(api chi.Router) {
		api.Mount("/auth", authHandler.Routes())
		api.Mount("/public/shares", shareHandler.PublicRoutes())

		api.Group(func(protected chi.Router) {
			protected.Use(auth.RequireUser(authUsecase))
			protected.Mount("/treatments", treatmentHandler.Routes())
			protected.Mount("/shares", shareHandler.Routes())
		})
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, r, httpx.NotFound("endpoint not found"))
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server listening", "port", cfg.Port, "env", cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
