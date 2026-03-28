package response

import "github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"

type UpdateHistory struct {
	Ent *historyEnt `json:"ent"`
}

func NewUpdateHistory(h entity.HairHistory) *UpdateHistory {
	return &UpdateHistory{
		Ent: newHistoryEnt(h),
	}
}
