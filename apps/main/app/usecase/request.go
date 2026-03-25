package usecase

import (
	"encoding/json"
	"net/http"
)

func NewCreateHistoryRequest(r *http.Request) (CreateHistoryRequest, error) {
	var req CreateHistoryRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		return CreateHistoryRequest{}, err
	}
	return req, nil
}

func NewUpdateHistoryRequest(r *http.Request) (UpdateHistoryRequest, error) {
	var req UpdateHistoryRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		return UpdateHistoryRequest{}, err
	}
	return req, nil
}

