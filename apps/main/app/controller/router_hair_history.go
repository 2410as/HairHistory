package controller

import "github.com/go-chi/chi/v5"

func registerHairHistoryRoutes(r chi.Router, deps Deps) {
	h := NewHairHistory(deps.HairHistory)
	r.Get("/users/{userId}/histories", h.List)
	r.Post("/users/{userId}/histories", h.Create)
	r.Put("/histories/{historyId}", h.Update)
	r.Delete("/histories/{historyId}", h.Delete)
}
