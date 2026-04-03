package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// Create handles POST /api/users/{userId}/histories
func (a HairHistory) Create(w http.ResponseWriter, r *http.Request) {
	in := &request.CreateHistory{
		UserID: chi.URLParam(r, "userId"),
	}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := request.NewCreateHistory(in)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.Create(r.Context(), req)
	if err != nil {
		render.ErrorFromUsecase(w, err)
		return
	}
	render.JSONWithStatus(w, res, http.StatusCreated)
}
