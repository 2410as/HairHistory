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
	ID            string
	Name          *string
	Email         *string
	LastLoginAt  *time.Time
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

