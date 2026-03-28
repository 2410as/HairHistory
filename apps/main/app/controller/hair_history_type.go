package controller

import "github.com/annasakai/hairhistorymemo/apps/main/app/usecase"

type HairHistory struct {
	hairHistoryUsecase *usecase.HairHistory
}

func NewHairHistory(hairHistoryUsecase *usecase.HairHistory) *HairHistory {
	return &HairHistory{hairHistoryUsecase: hairHistoryUsecase}
}
