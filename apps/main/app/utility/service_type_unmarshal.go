package utility

import (
	"encoding/json"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func UnmarshalServices(raw []byte) ([]entity.ServiceType, error) {
	var services []entity.ServiceType
	if err := json.Unmarshal(raw, &services); err != nil {
		return nil, err
	}
	return services, nil
}
