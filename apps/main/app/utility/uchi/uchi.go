package uchi

import "github.com/go-chi/chi/v5"

// NewRouter is a thin wrapper around chi.NewRouter.
// This exists to match the "uchi.NewRouter()" initialization style you requested.
func NewRouter() chi.Router {
	return chi.NewRouter()
}
