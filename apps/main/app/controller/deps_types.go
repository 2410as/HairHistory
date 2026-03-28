package controller

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase"
)

// DBPing is optional; when set, GET /api/health verifies the database is reachable.
type DBPing interface {
	Ping(ctx context.Context) error
}

// Deps bundles usecase instances for HTTP handlers.
type Deps struct {
	User        *usecase.User
	HairHistory *usecase.HairHistory
	DB          DBPing
}
