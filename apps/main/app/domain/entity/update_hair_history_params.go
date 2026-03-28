package entity

import "time"

type UpdateHairHistoryParams struct {
	Date        *time.Time
	Services    *[]ServiceType
	SalonName   *string
	StylistName *string
	Memo        *string
}
