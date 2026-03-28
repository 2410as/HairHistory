package entity

import "time"

type CreateHairHistoryParams struct {
	Date        time.Time
	Services    []ServiceType
	SalonName   string
	StylistName string
	Memo        string
}
