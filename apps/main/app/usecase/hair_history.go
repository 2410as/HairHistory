package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	hairhistorysvc "github.com/annasakai/hairhistorymemo/apps/main/app/domain/service/hairhistory"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

type HairHistory struct {
	hairHistorySvc hairhistorysvc.Service
}

func NewHairHistory(hairHistorySvc hairhistorysvc.Service) *HairHistory {
	return &HairHistory{hairHistorySvc: hairHistorySvc}
}

func (u HairHistory) List(ctx context.Context, r *request.ListHistories) (*response.ListHistories, error) {
	list, err := u.hairHistorySvc.ListByUserID(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	histories := make([]domain.HairHistory, 0, len(list))
	for _, h := range list {
		if h == nil {
			continue
		}
		histories = append(histories, *h)
	}
	return response.NewListHistories(histories), nil
}

func (u HairHistory) Create(ctx context.Context, r *request.CreateHistory) (*response.CreateHistory, error) {
	params := domain.CreateHairHistoryParams{
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
		return response.NewCreateHistory(domain.HairHistory{}), nil
	}
	return response.NewCreateHistory(*h), nil
}

func (u HairHistory) Update(ctx context.Context, r *request.UpdateHistory) (*response.UpdateHistory, error) {
	params := domain.UpdateHairHistoryParams{
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
		return response.NewUpdateHistory(domain.HairHistory{}), nil
	}
	return response.NewUpdateHistory(*h), nil
}

func (u HairHistory) Delete(ctx context.Context, r *request.DeleteHistory) (*response.DeleteHistory, error) {
	if err := u.hairHistorySvc.Delete(ctx, r.HistoryID); err != nil {
		return nil, err
	}
	return response.NewDeleteHistory(true), nil
}
