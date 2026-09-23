package share

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const linkColumns = `id, user_id, token, expires_at, revoked_at, created_at`

const createQuery = `
INSERT INTO share_links (id, user_id, token, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING ` + linkColumns

func (r *PostgresRepository) Create(ctx context.Context, link Link) (Link, error) {
	row := r.pool.QueryRow(ctx, createQuery, link.ID, link.UserID, link.Token, link.ExpiresAt, link.CreatedAt)
	return scanLink(row)
}

const listActiveByUserIDQuery = `
SELECT ` + linkColumns + `
FROM share_links
WHERE user_id = $1
  AND revoked_at IS NULL
  AND expires_at > $2
ORDER BY created_at DESC`

func (r *PostgresRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID, now time.Time) ([]Link, error) {
	rows, err := r.pool.Query(ctx, listActiveByUserIDQuery, userID, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]Link, 0)
	for rows.Next() {
		link, err := scanLink(rows)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

const findByTokenQuery = `
SELECT ` + linkColumns + `
FROM share_links
WHERE token = $1`

func (r *PostgresRepository) FindByToken(ctx context.Context, token string) (Link, error) {
	link, err := scanLink(r.pool.QueryRow(ctx, findByTokenQuery, token))
	if errors.Is(err, pgx.ErrNoRows) {
		return Link{}, ErrNotFound
	}
	return link, err
}

const revokeQuery = `
UPDATE share_links
SET revoked_at = $3
WHERE token = $2
  AND user_id = $1
  AND revoked_at IS NULL`

func (r *PostgresRepository) Revoke(ctx context.Context, userID uuid.UUID, token string, revokedAt time.Time) error {
	tag, err := r.pool.Exec(ctx, revokeQuery, userID, token, revokedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanLink(row rowScanner) (Link, error) {
	var link Link
	err := row.Scan(&link.ID, &link.UserID, &link.Token, &link.ExpiresAt, &link.RevokedAt, &link.CreatedAt)
	if err != nil {
		return Link{}, err
	}
	return link, nil
}
