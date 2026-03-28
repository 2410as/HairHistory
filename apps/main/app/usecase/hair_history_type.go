package usecase

import hairhistorysvc "github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/hairhistory"

type HairHistory struct {
	hairHistorySvc hairhistorysvc.Service
}

func NewHairHistory(hairHistorySvc hairhistorysvc.Service) *HairHistory {
	return &HairHistory{hairHistorySvc: hairHistorySvc}
}
