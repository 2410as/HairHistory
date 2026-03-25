package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

type ListHistoriesRequest struct {
	UserID string
}

type ListHistoriesResponse struct {
	Histories []domain.HairHistory
}

type ListHistoriesUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewListHistoriesUsecase(hairHistoryRepo domain.HairHistoryRepository) *ListHistoriesUsecase {
	return &ListHistoriesUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *ListHistoriesUsecase) Execute(ctx context.Context, req ListHistoriesRequest) (ListHistoriesResponse, error) {
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

