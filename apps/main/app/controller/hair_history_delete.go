package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// Delete handles DELETE /api/histories/{historyId}
func (a HairHistory) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewDeleteHistory(r)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.Delete(r.Context(), req)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, res)
}
