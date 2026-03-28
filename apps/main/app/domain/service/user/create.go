package user

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func (s *service) Create(ctx context.Context) (*entity.User, error) {
	return s.userRepo.Create(ctx)
}
