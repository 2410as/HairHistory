package render

import (
	"errors"
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

// ErrorFromUsecase maps domain / infrastructure errors to HTTP status codes.
func ErrorFromUsecase(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		ErrorJSON(w, "not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrInvalidInput):
		ErrorJSON(w, err.Error(), http.StatusBadRequest)
	default:
		ErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}
