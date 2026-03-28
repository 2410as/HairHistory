package repository

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

// User defines the persistence contract for User.
type User interface {
	Create(ctx context.Context) (*domain.User, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
}
