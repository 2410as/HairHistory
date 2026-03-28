package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// Update handles PUT /api/histories/{historyId}
func (a HairHistory) Update(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewUpdateHistory(r)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.Update(r.Context(), req)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, res)
}
