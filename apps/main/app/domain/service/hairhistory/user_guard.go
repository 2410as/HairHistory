package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func (s *service) ensureUserExists(ctx context.Context, userID string) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return domain.ErrNotFound
	}
	if u.IsDeactivated {
		return domain.ErrNotFound
	}
	return nil
}
