package infra

import (
	"context"
	"errors"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

// UserRepositoryPG is the PostgreSQL implementation.
// NOTE: Stub only for scaffolding.
type UserRepositoryPG struct{}

var _ domain.UserRepository = (*UserRepositoryPG)(nil)

func (r *UserRepositoryPG) Create(ctx context.Context) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepositoryPG) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}
