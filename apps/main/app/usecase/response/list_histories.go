package response

import "github.com/annasakai/hairhistorymemo/apps/main/app/domain"

type ListHistories struct {
	Histories []domain.HairHistory `json:"histories"`
}

func NewListHistories(histories []domain.HairHistory) *ListHistories {
	return &ListHistories{Histories: histories}
}

