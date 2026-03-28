package controller

import "net/http"

func registerHealthRoutes(mux *http.ServeMux) {
	healthController := NewHealth()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		healthController.Get(w, r)
	})
}
