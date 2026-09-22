package treatment

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("treatment: record not found")

const DateLayout = "2006-01-02"

type Treatment struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TreatedOn time.Time
	Services  []string
	SalonName *string
	Memo      *string
	Cost      *int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Input struct {
	TreatedOn string
	Services  []string
	SalonName *string
	Memo      *string
	Cost      *int
}

type validInput struct {
	treatedOn time.Time
	services  []string
	salonName *string
	memo      *string
	cost      *int
}

type CreateParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TreatedOn time.Time
	Services  []string
	SalonName *string
	Memo      *string
	Cost      *int
}

type UpdateParams struct {
	ID        uuid.UUID
	TreatedOn time.Time
	Services  []string
	SalonName *string
	Memo      *string
	Cost      *int
}

type Repository interface {
	Create(ctx context.Context, params CreateParams) (Treatment, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Treatment, error)
	GetByID(ctx context.Context, id uuid.UUID) (Treatment, error)
	Update(ctx context.Context, params UpdateParams) (Treatment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
