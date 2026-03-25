package utility

import (
	"encoding/json"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func MarshalServices(services []domain.ServiceType) ([]byte, error) {
	return json.Marshal(services)
}

func UnmarshalServices(raw []byte) ([]domain.ServiceType, error) {
	var services []domain.ServiceType
	if err := json.Unmarshal(raw, &services); err != nil {
		return nil, err
	}
	return services, nil
}

