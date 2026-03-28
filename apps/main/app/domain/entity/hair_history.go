package entity

import "time"

type HairHistory struct {
	ID          string
	UserID      string
	Date        time.Time
	Services    []ServiceType
	SalonName   string
	StylistName string
	Memo        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e *HairHistory) Created() {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
}

func (e *HairHistory) Updated() {
	e.UpdatedAt = time.Now()
}
