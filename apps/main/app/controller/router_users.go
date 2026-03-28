package controller

import "net/http"

func registerUserRoutes(mux *http.ServeMux, deps Deps) {
	usersController := NewUsers(deps.User)

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		usersController.Create(w, r)
	})
}
