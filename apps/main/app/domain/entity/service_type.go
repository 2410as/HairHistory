package entity

import (
	"errors"
	"fmt"
)

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

// IsKnownServiceType reports whether s is a supported service code for API / DB.
func IsKnownServiceType(s ServiceType) bool {
	switch s {
	case ServiceTypeColor, ServiceTypeBleach, ServiceTypeStraightPerm, ServiceTypeTreatment, ServiceTypeOther:
		return true
	default:
		return false
	}
}

// ValidateServices ensures at least one service and only known types.
func ValidateServices(services []ServiceType) error {
	if len(services) == 0 {
		return errServicesRequired
	}
	for _, s := range services {
		if !IsKnownServiceType(s) {
			return errUnknownServiceType(s)
		}
	}
	return nil
}

var errServicesRequired = errors.New("at least one service is required")

func errUnknownServiceType(s ServiceType) error {
	return fmt.Errorf("unknown service type: %q", s)
}
