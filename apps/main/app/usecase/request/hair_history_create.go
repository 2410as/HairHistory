package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

type CreateHistory struct {
	UserID      string               `json:"-"`
	Date        time.Time            `json:"date"`
	Services    []entity.ServiceType `json:"services"`
	SalonName   string               `json:"salonName"`
	StylistName string               `json:"stylistName"`
	Memo        string               `json:"memo"`
}

func (r *CreateHistory) Validate() error {
	if r.SalonName == "" {
		return errors.New("salonName is required")
	}
	if r.StylistName == "" {
		return errors.New("stylistName is required")
	}
	if r.Date.IsZero() {
		return errors.New("date is required")
	}
	return nil
}

func NewCreateHistory(httpReq *http.Request) (*CreateHistory, error) {
	req := &CreateHistory{}
	userID := chi.URLParam(httpReq, "userId")
	if userID == "" {
		return nil, errors.New("invalid userId")
	}
	req.UserID = userID

	if err := json.NewDecoder(httpReq.Body).Decode(req); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return req, nil
}
