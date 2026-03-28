package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

type UpdateHistory struct {
	HistoryID   string                `json:"-"`
	Date        *time.Time            `json:"date,omitempty"`
	Services    *[]entity.ServiceType `json:"services,omitempty"`
	SalonName   *string               `json:"salonName,omitempty"`
	StylistName *string               `json:"stylistName,omitempty"`
	Memo        *string               `json:"memo,omitempty"`
}

func NewUpdateHistory(httpReq *http.Request) (*UpdateHistory, error) {
	req := &UpdateHistory{}
	id := chi.URLParam(httpReq, "historyId")
	if id == "" {
		return nil, errors.New("invalid historyId")
	}
	req.HistoryID = id

	if err := json.NewDecoder(httpReq.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}
