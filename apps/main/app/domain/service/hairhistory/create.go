package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func (s *service) Create(ctx context.Context, userID string, req entity.CreateHairHistoryParams) (*entity.HairHistory, error) {
	return s.hairHistoryRepo.Create(ctx, userID, req)
}
