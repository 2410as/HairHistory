package controller

import "net/http"

func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()
	registerHealthRoutes(mux)
	registerUserRoutes(mux, deps)
	registerHairHistoryRoutes(mux, deps)
	return mux
}
