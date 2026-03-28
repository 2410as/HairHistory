package controller

import (
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller/render"
)

func (a Health) Get(w http.ResponseWriter, r *http.Request) {
	_ = r
	render.JSON(w, map[string]any{"ok": true})
}
