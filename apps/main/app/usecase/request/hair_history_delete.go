package request

import "errors"

type DeleteHistory struct {
	HistoryID string
}

func NewDeleteHistory(id string) (*DeleteHistory, error) {
	if id == "" {
		return nil, errors.New("invalid historyId")
	}
	return &DeleteHistory{HistoryID: id}, nil
}
