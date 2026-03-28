package request

import (
	"encoding/json"
	"net/http"
	"strings"
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

func NewUpdateHistory(httpReq *http.Request) (*UpdateHistory, error) {
	req := &UpdateHistory{}
	p := strings.TrimPrefix(httpReq.URL.Path, "/api/histories/")
	p = strings.Trim(p, "/")
	req.HistoryID = p

	if err := json.NewDecoder(httpReq.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}
