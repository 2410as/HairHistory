package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func (s *service) ListByUserID(ctx context.Context, userID string) ([]*domain.HairHistory, error) {
	return s.hairHistoryRepo.ListByUserID(ctx, userID)
}
