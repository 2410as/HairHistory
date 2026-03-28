package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
)

func (a Health) Get(w http.ResponseWriter, r *http.Request) {
	if a.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := a.db.Ping(ctx); err != nil {
			render.ErrorJSON(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
	}
	render.JSON(w, map[string]any{"ok": true})
}
