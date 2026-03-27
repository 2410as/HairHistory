package controller

import (
	"net/http"
)

func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()

	healthController := NewHealth()
	usersController := NewUsers(deps.User)
	hairHistoryController := NewHairHistory(deps.HairHistory)

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		healthController.Get(w, r)
	})

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		usersController.Create(w, r)
	})

	// /api/users/{userId}/histories
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			hairHistoryController.List(w, r)
		case http.MethodPost:
			hairHistoryController.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /api/histories/{historyId}
	mux.HandleFunc("/api/histories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			hairHistoryController.Update(w, r)
		case http.MethodDelete:
			hairHistoryController.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
