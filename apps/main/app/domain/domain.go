package domain

import (
	"context"
	"time"
)

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
	ID             string     `json:"id"`
	Name           *string    `json:"name,omitempty"`
	Email          *string    `json:"email,omitempty"`
	LastLoginAt   *time.Time `json:"lastLoginAt,omitempty"`
	IsDeactivated  bool       `json:"isDeactivated"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type HairHistory struct {
	ID           string        `json:"id"`
	UserID       string        `json:"userId"`
	Date         time.Time     `json:"date"`
	Services     []ServiceType `json:"services"`
	SalonName    string        `json:"salonName"`
	StylistName  string        `json:"stylistName"`
	Memo         string        `json:"memo"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

// UserRepository defines the persistence contract for User.
type UserRepository interface {
	Create(ctx context.Context) (*User, error)
	GetByID(ctx context.Context, userID string) (*User, error)
}

// HairHistoryRepository defines the persistence contract for HairHistory.
type HairHistoryRepository interface {
	ListByUserID(ctx context.Context, userID string) ([]*HairHistory, error)
	Create(ctx context.Context, userID string, req CreateHairHistoryParams) (*HairHistory, error)
	Update(ctx context.Context, historyID string, req UpdateHairHistoryParams) (*HairHistory, error)
	Delete(ctx context.Context, historyID string) error
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

