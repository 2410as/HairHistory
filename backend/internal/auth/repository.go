package auth

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

const upsertUserQuery = `
INSERT INTO users (id, google_sub, email, name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (google_sub) DO UPDATE
SET email = EXCLUDED.email,
    name = EXCLUDED.name,
    updated_at = now()
RETURNING id, google_sub, email, name, created_at, updated_at`

func (r *PostgresRepository) UpsertUser(ctx context.Context, params UpsertUserParams) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, upsertUserQuery, params.ID, params.GoogleSub, params.Email, params.Name).
		Scan(&user.ID, &user.GoogleSub, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

const findUserByIDQuery = `
SELECT id, google_sub, email, name, created_at, updated_at
FROM users
WHERE id = $1`

func (r *PostgresRepository) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, findUserByIDQuery, id).
		Scan(&user.ID, &user.GoogleSub, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

const createSessionQuery = `
INSERT INTO sessions (id, user_id, token_hash, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5)`

func (r *PostgresRepository) CreateSession(ctx context.Context, session Session) error {
	_, err := r.pool.Exec(ctx, createSessionQuery, session.ID, session.UserID, session.TokenHash, session.ExpiresAt, session.CreatedAt)
	return err
}

const findSessionByTokenHashQuery = `
SELECT id, user_id, token_hash, expires_at, created_at
FROM sessions
WHERE token_hash = $1`

func (r *PostgresRepository) FindSessionByTokenHash(ctx context.Context, tokenHash string) (Session, error) {
	var session Session
	err := r.pool.QueryRow(ctx, findSessionByTokenHashQuery, tokenHash).
		Scan(&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt, &session.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Session{}, ErrNotFound
	}
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func (r *PostgresRepository) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}

func (r *PostgresRepository) DeleteExpiredSessions(ctx context.Context, now time.Time) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at <= $1`, now)
	return err
}
