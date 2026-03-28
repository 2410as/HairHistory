package controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
)

// healthPingTimeout returns the DB ping deadline for GET /api/health.
// Override with HAIR_HEALTH_PING_TIMEOUT (e.g. "2s", "500ms"); invalid or empty defaults to 2s.
func healthPingTimeout() time.Duration {
	const defaultTimeout = 2 * time.Second
	s := os.Getenv("HAIR_HEALTH_PING_TIMEOUT")
	if s == "" {
		return defaultTimeout
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return defaultTimeout
	}
	return d
}

func (a Health) Get(w http.ResponseWriter, r *http.Request) {
	if a.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), healthPingTimeout())
		defer cancel()
		if err := a.db.Ping(ctx); err != nil {
			render.ErrorJSON(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
	}
	render.JSON(w, map[string]any{"ok": true})
}
