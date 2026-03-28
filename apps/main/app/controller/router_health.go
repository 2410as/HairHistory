package controller

import "github.com/go-chi/chi/v5"

func registerHealthRoutes(r chi.Router) {
	healthController := NewHealth()
	r.Get("/health", healthController.Get)
}
