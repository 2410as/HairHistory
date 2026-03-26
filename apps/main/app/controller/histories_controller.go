package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
	"github.com/annasakai/hairhistorymemo/apps/main/app/infra/logging"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

type Histories struct {
	listHistoriesUsecase   *usecase.ListHistoriesUsecase
	createHistoryUsecase   *usecase.CreateHistoryUsecase
	updateHistoryUsecase   *usecase.UpdateHistoryUsecase
	deleteHistoryUsecase   *usecase.DeleteHistoryUsecase
}

func NewHistories(
	listHistoriesUsecase *usecase.ListHistoriesUsecase,
	createHistoryUsecase *usecase.CreateHistoryUsecase,
	updateHistoryUsecase *usecase.UpdateHistoryUsecase,
	deleteHistoryUsecase *usecase.DeleteHistoryUsecase,
) *Histories {
	return &Histories{
		listHistoriesUsecase: listHistoriesUsecase,
		createHistoryUsecase: createHistoryUsecase,
		updateHistoryUsecase: updateHistoryUsecase,
		deleteHistoryUsecase: deleteHistoryUsecase,
	}
}

// List handles GET /api/users/{userId}/histories
func (a Histories) List(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		render.ErrorBadRequest(r.Context(), w, errors.New("userId is required"))
		return
	}

	resp, err := a.listHistoriesUsecase.Execute(r.Context(), usecase.ListHistoriesRequest{UserID: userID})
	if err != nil {
		logging.Errorf(r.Context(), "ListHistoriesUsecase.Execute err=%v", err)
		render.ErrorInternalServer(r.Context(), w, err)
		return
	}
	render.JSON(w, resp)
}

// Create handles POST /api/users/{userId}/histories
func (a Histories) Create(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		render.ErrorBadRequest(r.Context(), w, errors.New("userId is required"))
		return
	}

	req, err := usecase.NewCreateHistoryRequest(r)
	if err != nil {
		logging.Warningf(r.Context(), "NewCreateHistoryRequest err=%v", err)
		render.ErrorBadRequest(r.Context(), w, err)
		return
	}

	resp, err := a.createHistoryUsecase.Execute(r.Context(), userID, req)
	if err != nil {
		logging.Errorf(r.Context(), "CreateHistoryUsecase.Execute err=%v", err)
		render.ErrorInternalServer(r.Context(), w, err)
		return
	}
	render.JSON(w, resp)
}

// Update handles PUT /api/histories/{historyId}
func (a Histories) Update(w http.ResponseWriter, r *http.Request) {
	historyID := chi.URLParam(r, "historyId")
	if historyID == "" {
		render.ErrorBadRequest(r.Context(), w, errors.New("historyId is required"))
		return
	}

	req, err := usecase.NewUpdateHistoryRequest(r)
	if err != nil {
		logging.Warningf(r.Context(), "NewUpdateHistoryRequest err=%v", err)
		render.ErrorBadRequest(r.Context(), w, err)
		return
	}

	resp, err := a.updateHistoryUsecase.Execute(r.Context(), historyID, req)
	if err != nil {
		logging.Errorf(r.Context(), "UpdateHistoryUsecase.Execute err=%v", err)
		render.ErrorInternalServer(r.Context(), w, err)
		return
	}
	render.JSON(w, resp)
}

// Delete handles DELETE /api/histories/{historyId}
func (a Histories) Delete(w http.ResponseWriter, r *http.Request) {
	historyID := chi.URLParam(r, "historyId")
	if historyID == "" {
		render.ErrorBadRequest(r.Context(), w, errors.New("historyId is required"))
		return
	}

	resp, err := a.deleteHistoryUsecase.Execute(r.Context(), historyID, usecase.DeleteHistoryRequest{})
	if err != nil {
		logging.Errorf(r.Context(), "DeleteHistoryUsecase.Execute err=%v", err)
		render.ErrorInternalServer(r.Context(), w, err)
		return
	}
	render.JSON(w, resp)
}

