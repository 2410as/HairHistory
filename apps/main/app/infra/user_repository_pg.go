package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

// UserRepositoryPG persists users in PostgreSQL.
type UserRepositoryPG struct {
	Pool *pgxpool.Pool
}

var _ domain.UserRepository = (*UserRepositoryPG)(nil)

// NewUserRepositoryPG returns a repository backed by pool. pool must be non-nil.
func NewUserRepositoryPG(pool *pgxpool.Pool) (*UserRepositoryPG, error) {
	if pool == nil {
		return nil, errors.New("infra: nil pool")
	}
	return &UserRepositoryPG{Pool: pool}, nil
}

func (r *UserRepositoryPG) Create(ctx context.Context) (*entity.User, error) {
	id := uuid.New().String()

	row := r.Pool.QueryRow(ctx, `
INSERT INTO users (id, is_deactivated, created_at, updated_at)
VALUES ($1, false, now(), now())
RETURNING id, name, email, last_login_at, is_deactivated, created_at, updated_at
`, id)

	return scanUser(row)
}

func (r *UserRepositoryPG) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	row := r.Pool.QueryRow(ctx, `
SELECT id, name, email, last_login_at, is_deactivated, created_at, updated_at
FROM users WHERE id = $1
`, userID)

	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User
	var name, email sql.NullString
	var lastLogin sql.NullTime

	if err := row.Scan(
		&u.ID,
		&name,
		&email,
		&lastLogin,
		&u.IsDeactivated,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	if name.Valid {
		s := name.String
		u.Name = &s
	}
	if email.Valid {
		s := email.String
		u.Email = &s
	}
	if lastLogin.Valid {
		t := lastLogin.Time
		u.LastLoginAt = &t
	}
	return &u, nil
}
