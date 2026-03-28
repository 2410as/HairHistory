package response

import "github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"

type ListHistories struct {
	List []*historyEnt `json:"list"`
}

func NewListHistories(histories []entity.HairHistory) *ListHistories {
	return &ListHistories{
		List: newHistoryEntList(histories),
	}
}
