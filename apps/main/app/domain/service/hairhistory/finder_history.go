package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func (s *service) ListByUserID(ctx context.Context, userID string) ([]*entity.HairHistory, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}
	return s.hairHistoryRepo.ListByUserID(ctx, userID)
}
