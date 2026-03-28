package request

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DeleteHistory struct {
	HistoryID string
}

func NewDeleteHistory(httpReq *http.Request) (*DeleteHistory, error) {
	id := chi.URLParam(httpReq, "historyId")
	if id == "" {
		return nil, errors.New("invalid historyId")
	}
	return &DeleteHistory{HistoryID: id}, nil
}
