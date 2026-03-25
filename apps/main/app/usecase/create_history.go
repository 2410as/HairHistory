package usecase

import (
	"context"
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

type CreateHistoryRequest struct {
	Date        time.Time           `json:"date"`
	Services    []domain.ServiceType `json:"services"`
	SalonName   string              `json:"salonName"`
	StylistName string              `json:"stylistName"`
	Memo        string              `json:"memo"`
}

type CreateHistoryResponse struct {
	History domain.HairHistory `json:"history"`
}

type CreateHistoryUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewCreateHistoryUsecase(hairHistoryRepo domain.HairHistoryRepository) *CreateHistoryUsecase {
	return &CreateHistoryUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *CreateHistoryUsecase) Execute(ctx context.Context, userID string, req CreateHistoryRequest) (CreateHistoryResponse, error) {
	params := domain.CreateHairHistoryParams{
		Date:        req.Date,
		Services:    req.Services,
		SalonName:   req.SalonName,
		StylistName: req.StylistName,
		Memo:        req.Memo,
	}

	h, err := u.hairHistoryRepo.Create(ctx, userID, params)
	if err != nil {
		return CreateHistoryResponse{}, err
	}
	if h == nil {
		return CreateHistoryResponse{}, nil
	}

	return CreateHistoryResponse{History: *h}, nil
}

