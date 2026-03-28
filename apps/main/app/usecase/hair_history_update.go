package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

func (u HairHistory) Update(ctx context.Context, r *request.UpdateHistory) (*response.UpdateHistory, error) {
	params := entity.UpdateHairHistoryParams{
		Date:        r.Date,
		Services:    r.Services,
		SalonName:   r.SalonName,
		StylistName: r.StylistName,
		Memo:        r.Memo,
	}

	h, err := u.hairHistorySvc.Update(ctx, r.HistoryID, params)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return response.NewUpdateHistory(entity.HairHistory{}), nil
	}
	return response.NewUpdateHistory(*h), nil
}
