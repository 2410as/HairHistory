package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, v any) {
	JSONWithStatus(w, v, http.StatusOK)
}

func JSONWithStatus(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
