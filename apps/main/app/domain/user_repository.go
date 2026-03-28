package domain

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

// UserRepository is the persistence contract for User (implemented in infra).
type UserRepository interface {
	Create(ctx context.Context) (*entity.User, error)
	GetByID(ctx context.Context, userID string) (*entity.User, error)
}
