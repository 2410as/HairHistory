package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type Users struct {
	userUsecase *usecase.User
}

func NewUsers(userUsecase *usecase.User) *Users {
	return &Users{userUsecase: userUsecase}
}

// Create handles POST /api/users
func (a Users) Create(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewCreateUser(r)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.userUsecase.Create(r.Context(), req)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, res)
}

