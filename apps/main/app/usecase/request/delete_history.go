package request

import (
	"errors"
	"net/http"
	"strings"
)

type DeleteHistory struct {
	HistoryID string
}

func NewDeleteHistory(httpReq *http.Request) (*DeleteHistory, error) {
	// Expected path: /api/histories/{historyId}
	p := strings.TrimPrefix(httpReq.URL.Path, "/api/histories/")
	p = strings.Trim(p, "/")
	if p == "" {
		return nil, errors.New("invalid historyId path")
	}
	return &DeleteHistory{HistoryID: p}, nil
}

