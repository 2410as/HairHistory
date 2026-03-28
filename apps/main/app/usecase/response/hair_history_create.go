package response

import "github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"

type CreateHistory struct {
	Ent *historyEnt `json:"ent"`
}

func NewCreateHistory(h entity.HairHistory) *CreateHistory {
	return &CreateHistory{
		Ent: newHistoryEnt(h),
	}
}
