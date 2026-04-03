package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// Update handles PUT /api/histories/{historyId}
func (a HairHistory) Update(w http.ResponseWriter, r *http.Request) {
	in := &request.UpdateHistory{
		HistoryID: chi.URLParam(r, "historyId"),
	}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := request.NewUpdateHistory(in)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.Update(r.Context(), req)
	if err != nil {
		render.ErrorFromUsecase(w, err)
		return
	}
	render.JSON(w, res)
}
