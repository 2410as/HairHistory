package user

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

func (s *service) Create(ctx context.Context) (*domain.User, error) {
	return s.userRepo.Create(ctx)
}
