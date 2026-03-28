package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// Create handles POST /api/users/{userId}/histories
func (a HairHistory) Create(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewCreateHistory(r)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.Create(r.Context(), req)
	if err != nil {
		render.ErrorFromUsecase(w, err)
		return
	}
	render.JSON(w, res)
}
