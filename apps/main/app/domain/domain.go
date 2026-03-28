package domain

import "time"

// ServiceType is a type-safe representation of hair services.
type ServiceType string

const (
	ServiceTypeColor        ServiceType = "color"
	ServiceTypeBleach       ServiceType = "bleach"
	ServiceTypeStraightPerm ServiceType = "straight_perm" // 縮毛矯正
	ServiceTypeTreatment    ServiceType = "treatment"
	ServiceTypeOther        ServiceType = "other"
)

type User struct {
	ID            string
	Name          *string
	Email         *string
	LastLoginAt   *time.Time
	IsDeactivated bool
	CreatedAt     time.Time
}

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

type CreateHairHistoryParams struct {
	Date        time.Time
	Services    []ServiceType
	SalonName   string
	StylistName string
	Memo        string
}

type UpdateHairHistoryParams struct {
	Date        *time.Time
	Services    *[]ServiceType
	SalonName   *string
	StylistName *string
	Memo        *string
}
