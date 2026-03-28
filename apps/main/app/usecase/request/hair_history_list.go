package request

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ListHistories struct {
	UserID string
}

func NewListHistories(httpReq *http.Request) (*ListHistories, error) {
	userID := chi.URLParam(httpReq, "userId")
	if userID == "" {
		return nil, errors.New("invalid userId")
	}
	return &ListHistories{UserID: userID}, nil
}
