package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

// HairHistoryRepositoryPG persists hair histories in PostgreSQL.
type HairHistoryRepositoryPG struct {
	Pool *pgxpool.Pool
}

var _ domain.HairHistoryRepository = (*HairHistoryRepositoryPG)(nil)

func servicesToJSON(s []entity.ServiceType) ([]byte, error) {
	strs := make([]string, len(s))
	for i, v := range s {
		strs[i] = string(v)
	}
	return json.Marshal(strs)
}

func servicesFromJSON(b []byte) ([]entity.ServiceType, error) {
	var strs []string
	if err := json.Unmarshal(b, &strs); err != nil {
		return nil, err
	}
	out := make([]entity.ServiceType, len(strs))
	for i, x := range strs {
		out[i] = entity.ServiceType(x)
	}
	return out, nil
}

func (r *HairHistoryRepositoryPG) ListByUserID(ctx context.Context, userID string) ([]*entity.HairHistory, error) {
	if r.Pool == nil {
		return nil, errors.New("infra: nil pool")
	}
	rows, err := r.Pool.Query(ctx, `
SELECT id, user_id, date, services, salon_name, stylist_name, memo, created_at, updated_at
FROM hair_histories
WHERE user_id = $1
ORDER BY date DESC, created_at DESC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list histories: %w", err)
	}
	defer rows.Close()

	var list []*entity.HairHistory
	for rows.Next() {
		h, err := scanHairHistory(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *HairHistoryRepositoryPG) Create(ctx context.Context, userID string, req entity.CreateHairHistoryParams) (*entity.HairHistory, error) {
	if r.Pool == nil {
		return nil, errors.New("infra: nil pool")
	}
	servicesJSON, err := servicesToJSON(req.Services)
	if err != nil {
		return nil, fmt.Errorf("services json: %w", err)
	}
	date := toDateUTC(req.Date)

	row := r.Pool.QueryRow(ctx, `
INSERT INTO hair_histories (user_id, date, services, salon_name, stylist_name, memo, created_at, updated_at)
VALUES ($1, $2, $3::jsonb, $4, $5, $6, now(), now())
RETURNING id, user_id, date, services, salon_name, stylist_name, memo, created_at, updated_at
`, userID, date, servicesJSON, req.SalonName, req.StylistName, req.Memo)

	h, err := scanHairHistory(row)
	if err != nil {
		return nil, fmt.Errorf("create history: %w", err)
	}
	return h, nil
}

func (r *HairHistoryRepositoryPG) Update(ctx context.Context, historyID string, req entity.UpdateHairHistoryParams) (*entity.HairHistory, error) {
	if r.Pool == nil {
		return nil, errors.New("infra: nil pool")
	}
	if _, err := uuid.Parse(historyID); err != nil {
		return nil, fmt.Errorf("invalid history id: %w", err)
	}

	existing, err := r.getByID(ctx, historyID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	date := existing.Date
	if req.Date != nil {
		date = toDateUTC(*req.Date)
	}
	services := existing.Services
	if req.Services != nil {
		services = *req.Services
	}
	salonName := existing.SalonName
	if req.SalonName != nil {
		salonName = *req.SalonName
	}
	stylistName := existing.StylistName
	if req.StylistName != nil {
		stylistName = *req.StylistName
	}
	memo := existing.Memo
	if req.Memo != nil {
		memo = *req.Memo
	}

	servicesJSON, err := servicesToJSON(services)
	if err != nil {
		return nil, fmt.Errorf("services json: %w", err)
	}

	row := r.Pool.QueryRow(ctx, `
UPDATE hair_histories
SET date = $2, services = $3::jsonb, salon_name = $4, stylist_name = $5, memo = $6, updated_at = now()
WHERE id = $1
RETURNING id, user_id, date, services, salon_name, stylist_name, memo, created_at, updated_at
`, historyID, date, servicesJSON, salonName, stylistName, memo)

	h, err := scanHairHistory(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("update history: %w", err)
	}
	return h, nil
}

func (r *HairHistoryRepositoryPG) Delete(ctx context.Context, historyID string) error {
	if r.Pool == nil {
		return errors.New("infra: nil pool")
	}
	if _, err := uuid.Parse(historyID); err != nil {
		return fmt.Errorf("invalid history id: %w", err)
	}
	ct, err := r.Pool.Exec(ctx, `DELETE FROM hair_histories WHERE id = $1`, historyID)
	if err != nil {
		return fmt.Errorf("delete history: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("history not found: %s", historyID)
	}
	return nil
}

func (r *HairHistoryRepositoryPG) getByID(ctx context.Context, id string) (*entity.HairHistory, error) {
	row := r.Pool.QueryRow(ctx, `
SELECT id, user_id, date, services, salon_name, stylist_name, memo, created_at, updated_at
FROM hair_histories WHERE id = $1
`, id)
	h, err := scanHairHistory(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return h, nil
}

func scanHairHistory(row pgx.Row) (*entity.HairHistory, error) {
	var h entity.HairHistory
	var servicesBytes []byte

	if err := row.Scan(
		&h.ID,
		&h.UserID,
		&h.Date,
		&servicesBytes,
		&h.SalonName,
		&h.StylistName,
		&h.Memo,
		&h.CreatedAt,
		&h.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan hair_history: %w", err)
	}
	svcs, err := servicesFromJSON(servicesBytes)
	if err != nil {
		return nil, fmt.Errorf("services from json: %w", err)
	}
	h.Services = svcs
	return &h, nil
}

func toDateUTC(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}
