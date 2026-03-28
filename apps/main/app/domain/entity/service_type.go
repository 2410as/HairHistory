package entity

// ServiceType is a type-safe representation of salon menu items (color, perm, etc.).
// NOTE: Not related to domain/service (application services).
type ServiceType string

const (
	ServiceTypeColor        ServiceType = "color"
	ServiceTypeBleach       ServiceType = "bleach"
	ServiceTypeStraightPerm ServiceType = "straight_perm" // 縮毛矯正
	ServiceTypeTreatment    ServiceType = "treatment"
	ServiceTypeOther        ServiceType = "other"
)
