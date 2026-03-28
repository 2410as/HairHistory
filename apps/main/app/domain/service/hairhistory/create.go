package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func (s *service) Create(ctx context.Context, userID string, req domain.CreateHairHistoryParams) (*domain.HairHistory, error) {
	return s.hairHistoryRepo.Create(ctx, userID, req)
}
