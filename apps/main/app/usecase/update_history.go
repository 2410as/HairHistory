package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type UpdateHistoryResponse struct {
	History domain.HairHistory `json:"history"`
}

type UpdateHistoryUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewUpdateHistoryUsecase(hairHistoryRepo domain.HairHistoryRepository) *UpdateHistoryUsecase {
	return &UpdateHistoryUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *UpdateHistoryUsecase) Execute(ctx context.Context, req *request.UpdateHistory) (UpdateHistoryResponse, error) {
	params := domain.UpdateHairHistoryParams{
		Date:        req.Date,
		Services:    req.Services,
		SalonName:   req.SalonName,
		StylistName: req.StylistName,
		Memo:        req.Memo,
	}

	h, err := u.hairHistoryRepo.Update(ctx, req.HistoryID, params)
	if err != nil {
		return UpdateHistoryResponse{}, err
	}
	if h == nil {
		return UpdateHistoryResponse{}, nil
	}

	return UpdateHistoryResponse{History: *h}, nil
}

