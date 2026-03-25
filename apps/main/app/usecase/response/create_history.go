package response

import "github.com/annasakai/hairhistorymemo/apps/main/app/domain"

type CreateHistory struct {
	History domain.HairHistory `json:"history"`
}

func NewCreateHistory(h domain.HairHistory) *CreateHistory {
	return &CreateHistory{History: h}
}

