package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func (s *service) Update(ctx context.Context, historyID string, req entity.UpdateHairHistoryParams) (*entity.HairHistory, error) {
	return s.hairHistoryRepo.Update(ctx, historyID, req)
}
