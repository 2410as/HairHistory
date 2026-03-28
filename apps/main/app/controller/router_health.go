package controller

import "github.com/go-chi/chi/v5"

func registerHealthRoutes(r chi.Router, deps Deps) {
	healthController := NewHealth(deps.DB)
	r.Get("/health", healthController.Get)
}
