package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/infra/logging"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

type Users struct {
	createUserUsecase *usecase.CreateUserUsecase
}

func NewUsers(createUserUsecase *usecase.CreateUserUsecase) *Users {
	return &Users{createUserUsecase: createUserUsecase}
}

func (a Users) Create(w http.ResponseWriter, r *http.Request) {
	req := usecase.CreateUserRequest{}

	res, err := a.createUserUsecase.Execute(r.Context(), req)
	if err != nil {
		logging.Errorf(r.Context(), "CreateUserUsecase.Execute err=%v", err)
		render.ErrorInternalServer(r.Context(), w, err)
		return
	}
	render.JSON(w, res)
}

