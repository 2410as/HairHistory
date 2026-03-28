package controller

import "github.com/go-chi/chi/v5"

func registerUserRoutes(r chi.Router, deps Deps) {
	usersController := NewUsers(deps.User)
	r.Post("/users", usersController.Create)
}
