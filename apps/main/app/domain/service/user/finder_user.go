package user

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func (s *service) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
