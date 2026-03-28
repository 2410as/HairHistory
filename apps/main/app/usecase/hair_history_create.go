package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

func (u HairHistory) Create(ctx context.Context, r *request.CreateHistory) (*response.CreateHistory, error) {
	params := entity.CreateHairHistoryParams{
		Date:        r.Date,
		Services:    r.Services,
		SalonName:   r.SalonName,
		StylistName: r.StylistName,
		Memo:        r.Memo,
	}

	h, err := u.hairHistorySvc.Create(ctx, r.UserID, params)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return response.NewCreateHistory(entity.HairHistory{}), nil
	}
	return response.NewCreateHistory(*h), nil
}
