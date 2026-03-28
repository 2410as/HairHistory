package response

import (
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

type historyEnt struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Date        time.Time `json:"date"`
	Services    []string  `json:"services"`
	SalonName   string    `json:"salonName"`
	StylistName string    `json:"stylistName"`
	Memo        string    `json:"memo"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newHistoryEnt(h entity.HairHistory) *historyEnt {
	services := make([]string, len(h.Services))
	for i, s := range h.Services {
		services[i] = string(s)
	}
	return &historyEnt{
		ID:          h.ID,
		UserID:      h.UserID,
		Date:        h.Date,
		Services:    services,
		SalonName:   h.SalonName,
		StylistName: h.StylistName,
		Memo:        h.Memo,
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
	}
}

func newHistoryEntList(histories []entity.HairHistory) []*historyEnt {
	out := make([]*historyEnt, 0, len(histories))
	for i := range histories {
		out = append(out, newHistoryEnt(histories[i]))
	}
	return out
}
