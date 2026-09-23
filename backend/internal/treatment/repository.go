package treatment

import (
	"context"
	"errors"

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

const treatmentColumns = `id, user_id, treated_on, services, salon_name, memo, cost, created_at, updated_at`

const createQuery = `
INSERT INTO treatments (id, user_id, treated_on, services, salon_name, memo, cost)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING ` + treatmentColumns

func (r *PostgresRepository) Create(ctx context.Context, params CreateParams) (Treatment, error) {
	row := r.pool.QueryRow(ctx, createQuery,
		params.ID, params.UserID, params.TreatedOn, params.Services,
		params.SalonName, params.Memo, params.Cost)
	return scanTreatment(row)
}

const listByUserIDQuery = `
SELECT ` + treatmentColumns + `
FROM treatments
WHERE user_id = $1
ORDER BY treated_on DESC, created_at DESC`

func (r *PostgresRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]Treatment, error) {
	rows, err := r.pool.Query(ctx, listByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	treatments := make([]Treatment, 0)
	for rows.Next() {
		item, err := scanTreatment(rows)
		if err != nil {
			return nil, err
		}
		treatments = append(treatments, item)
	}
	return treatments, rows.Err()
}

const getByIDQuery = `
SELECT ` + treatmentColumns + `
FROM treatments
WHERE id = $1`

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (Treatment, error) {
	found, err := scanTreatment(r.pool.QueryRow(ctx, getByIDQuery, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return Treatment{}, ErrNotFound
	}
	return found, err
}

const updateQuery = `
UPDATE treatments
SET treated_on = $2,
    services = $3,
    salon_name = $4,
    memo = $5,
    cost = $6,
    updated_at = now()
WHERE id = $1
RETURNING ` + treatmentColumns

func (r *PostgresRepository) Update(ctx context.Context, params UpdateParams) (Treatment, error) {
	row := r.pool.QueryRow(ctx, updateQuery,
		params.ID, params.TreatedOn, params.Services, params.SalonName, params.Memo, params.Cost)
	updated, err := scanTreatment(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Treatment{}, ErrNotFound
	}
	return updated, err
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM treatments WHERE id = $1`, id)
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

func scanTreatment(row rowScanner) (Treatment, error) {
	var item Treatment
	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.TreatedOn,
		&item.Services,
		&item.SalonName,
		&item.Memo,
		&item.Cost,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return Treatment{}, err
	}
	return item, nil
}
