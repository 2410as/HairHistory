package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

func (u HairHistory) List(ctx context.Context, r *request.ListHistories) (*response.ListHistories, error) {
	list, err := u.hairHistorySvc.ListByUserID(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	histories := make([]entity.HairHistory, 0, len(list))
	for _, h := range list {
		if h == nil {
			continue
		}
		histories = append(histories, *h)
	}
	return response.NewListHistories(histories), nil
}
