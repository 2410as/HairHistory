package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		registerHealthRoutes(r)
		registerUserRoutes(r, deps)
		registerHairHistoryRoutes(r, deps)
	})
	return r
}
