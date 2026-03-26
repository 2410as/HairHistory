package usecase

import (
	"context"
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

type CreateHistoryRequest struct {
	Date        time.Time
	Services    []domain.ServiceType
	SalonName   string
	StylistName string
	Memo        string
}

type CreateHistoryResponse struct {
	History domain.HairHistory
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

