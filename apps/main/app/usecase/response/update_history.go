package response

import "github.com/annasakai/hairhistorymemo/apps/main/app/domain"

type UpdateHistory struct {
	History domain.HairHistory `json:"history"`
}

func NewUpdateHistory(h domain.HairHistory) *UpdateHistory {
	return &UpdateHistory{History: h}
}

