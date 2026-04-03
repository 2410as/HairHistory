package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// List handles GET /api/users/{userId}/histories
func (a HairHistory) List(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewListHistories(chi.URLParam(r, "userId"))
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.List(r.Context(), req)
	if err != nil {
		render.ErrorFromUsecase(w, err)
		return
	}
	render.JSON(w, res)
}
