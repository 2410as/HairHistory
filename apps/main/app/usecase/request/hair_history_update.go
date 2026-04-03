package request

import (
	"errors"
	"time"

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

func (r *UpdateHistory) Validate() error {
	if r.HistoryID == "" {
		return errors.New("invalid historyId")
	}
	if r.Services == nil {
		return nil
	}
	return entity.ValidateServices(*r.Services)
}

func NewUpdateHistory(req *UpdateHistory) (*UpdateHistory, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return req, nil
}
