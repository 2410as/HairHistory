package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
	// Minimal validation for now.
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
	p := strings.TrimPrefix(httpReq.URL.Path, "/api/users/")
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	if len(parts) != 2 || parts[1] != "histories" || parts[0] == "" {
		return nil, errors.New("invalid userId path")
	}
	req.UserID = parts[0]

	if err := json.NewDecoder(httpReq.Body).Decode(req); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return req, nil
}
