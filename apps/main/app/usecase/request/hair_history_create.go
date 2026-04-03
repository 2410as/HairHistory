package request

import (
	"errors"
	"time"

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
	if r.UserID == "" {
		return errors.New("invalid userId")
	}
	if r.Date.IsZero() {
		return errors.New("date is required")
	}
	return entity.ValidateServices(r.Services)
}

func NewCreateHistory(req *CreateHistory) (*CreateHistory, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return req, nil
}
