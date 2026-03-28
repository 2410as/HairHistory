package utility

import (
	"encoding/json"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func MarshalServices(services []entity.ServiceType) ([]byte, error) {
	return json.Marshal(services)
}
