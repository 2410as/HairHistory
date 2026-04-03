package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

// Create handles POST /api/users
func (a Users) Create(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewCreateUser()
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.userUsecase.Create(r.Context(), req)
	if err != nil {
		render.ErrorFromUsecase(w, err)
		return
	}
	render.JSONWithStatus(w, res, http.StatusCreated)
}
