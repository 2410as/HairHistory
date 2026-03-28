package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func (s *service) Update(ctx context.Context, historyID string, req domain.UpdateHairHistoryParams) (*domain.HairHistory, error) {
	return s.hairHistoryRepo.Update(ctx, historyID, req)
}
