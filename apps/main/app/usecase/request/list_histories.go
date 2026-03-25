package request

import (
	"errors"
	"net/http"
	"strings"
)

type ListHistories struct {
	UserID string
}

func NewListHistories(httpReq *http.Request) (*ListHistories, error) {
	// Expected path: /api/users/{userId}/histories
	p := strings.TrimPrefix(httpReq.URL.Path, "/api/users/")
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	if len(parts) != 2 || parts[1] != "histories" || parts[0] == "" {
		return nil, errors.New("invalid userId path")
	}
	return &ListHistories{UserID: parts[0]}, nil
}

