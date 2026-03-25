package render

import (
	"context"
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(v)
}

func ErrorJSON(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": message,
	})
}

func ErrorBadRequest(_ context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		ErrorJSON(w, "bad request", http.StatusBadRequest)
		return
	}
	ErrorJSON(w, err.Error(), http.StatusBadRequest)
}

func ErrorInternalServer(_ context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		ErrorJSON(w, "internal server error", http.StatusInternalServerError)
		return
	}
	ErrorJSON(w, err.Error(), http.StatusInternalServerError)
}

