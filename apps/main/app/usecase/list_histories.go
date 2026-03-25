package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type ListHistoriesResponse struct {
	Histories []domain.HairHistory `json:"histories"`
}

type ListHistoriesUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewListHistoriesUsecase(hairHistoryRepo domain.HairHistoryRepository) *ListHistoriesUsecase {
	return &ListHistoriesUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *ListHistoriesUsecase) Execute(ctx context.Context, req *request.ListHistories) (ListHistoriesResponse, error) {
	list, err := u.hairHistoryRepo.ListByUserID(ctx, req.UserID)
	if err != nil {
		return ListHistoriesResponse{}, err
	}

	histories := make([]domain.HairHistory, 0, len(list))
	for _, h := range list {
		if h == nil {
			continue
		}
		histories = append(histories, *h)
	}

	return ListHistoriesResponse{Histories: histories}, nil
}

