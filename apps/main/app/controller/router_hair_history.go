package controller

import "net/http"

func registerHairHistoryRoutes(mux *http.ServeMux, deps Deps) {
	hairHistoryController := NewHairHistory(deps.HairHistory)

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
}
