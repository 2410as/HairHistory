package request

import "errors"

type ListHistories struct {
	UserID string
}

func NewListHistories(userID string) (*ListHistories, error) {
	if userID == "" {
		return nil, errors.New("invalid userId")
	}
	return &ListHistories{UserID: userID}, nil
}
