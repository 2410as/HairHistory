package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type HairHistory struct {
	hairHistoryUsecase *usecase.HairHistory
}

func NewHairHistory(hairHistoryUsecase *usecase.HairHistory) *HairHistory {
	return &HairHistory{hairHistoryUsecase: hairHistoryUsecase}
}

// List handles GET /api/users/{userId}/histories
func (a HairHistory) List(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewListHistories(r)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.List(r.Context(), req)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, res)
}

// Create handles POST /api/users/{userId}/histories
func (a HairHistory) Create(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewCreateHistory(r)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.hairHistoryUsecase.Create(r.Context(), req)
	if err != nil {
		render.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, res)
}

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

