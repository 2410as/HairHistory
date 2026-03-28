package response

import (
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

// ---

type ListHistories struct {
	List []*historyEnt `json:"list"`
}

func NewListHistories(histories []domain.HairHistory) *ListHistories {
	return &ListHistories{
		List: newHistoryEntList(histories),
	}
}

// ---

type CreateHistory struct {
	Ent *historyEnt `json:"ent"`
}

func NewCreateHistory(h domain.HairHistory) *CreateHistory {
	return &CreateHistory{
		Ent: newHistoryEnt(h),
	}
}

// ---

type UpdateHistory struct {
	Ent *historyEnt `json:"ent"`
}

func NewUpdateHistory(h domain.HairHistory) *UpdateHistory {
	return &UpdateHistory{
		Ent: newHistoryEnt(h),
	}
}

// ---

type DeleteHistory struct {
	OK bool `json:"ok"`
}

func NewDeleteHistory(ok bool) *DeleteHistory {
	return &DeleteHistory{OK: ok}
}

// ---

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

func newHistoryEnt(h domain.HairHistory) *historyEnt {
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

func newHistoryEntList(histories []domain.HairHistory) []*historyEnt {
	out := make([]*historyEnt, 0, len(histories))
	for i := range histories {
		out = append(out, newHistoryEnt(histories[i]))
	}
	return out
}
