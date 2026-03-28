package user

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func (s *service) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
