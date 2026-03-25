package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type CreateHistoryResponse struct {
	History domain.HairHistory `json:"history"`
}

type CreateHistoryUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewCreateHistoryUsecase(hairHistoryRepo domain.HairHistoryRepository) *CreateHistoryUsecase {
	return &CreateHistoryUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *CreateHistoryUsecase) Execute(ctx context.Context, req *request.CreateHistory) (CreateHistoryResponse, error) {
	params := domain.CreateHairHistoryParams{
		Date:        req.Date,
		Services:    req.Services,
		SalonName:   req.SalonName,
		StylistName: req.StylistName,
		Memo:        req.Memo,
	}

	h, err := u.hairHistoryRepo.Create(ctx, req.UserID, params)
	if err != nil {
		return CreateHistoryResponse{}, err
	}
	if h == nil {
		return CreateHistoryResponse{}, nil
	}

	return CreateHistoryResponse{History: *h}, nil
}

