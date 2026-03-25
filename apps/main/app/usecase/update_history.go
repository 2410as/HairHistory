package usecase

import (
	"context"
	"time"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

type UpdateHistoryRequest struct {
	Date        *time.Time           `json:"date,omitempty"`
	Services    *[]domain.ServiceType `json:"services,omitempty"`
	SalonName   *string             `json:"salonName,omitempty"`
	StylistName *string             `json:"stylistName,omitempty"`
	Memo        *string             `json:"memo,omitempty"`
}

type UpdateHistoryResponse struct {
	History domain.HairHistory `json:"history"`
}

type UpdateHistoryUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewUpdateHistoryUsecase(hairHistoryRepo domain.HairHistoryRepository) *UpdateHistoryUsecase {
	return &UpdateHistoryUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *UpdateHistoryUsecase) Execute(ctx context.Context, historyID string, req UpdateHistoryRequest) (UpdateHistoryResponse, error) {
	params := domain.UpdateHairHistoryParams{
		Date:        req.Date,
		Services:    req.Services,
		SalonName:   req.SalonName,
		StylistName: req.StylistName,
		Memo:        req.Memo,
	}

	h, err := u.hairHistoryRepo.Update(ctx, historyID, params)
	if err != nil {
		return UpdateHistoryResponse{}, err
	}
	if h == nil {
		return UpdateHistoryResponse{}, nil
	}

	return UpdateHistoryResponse{History: *h}, nil
}

