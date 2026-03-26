package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
)

type DeleteHistoryRequest struct{}

type DeleteHistoryResponse struct {
	OK bool
}

type DeleteHistoryUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewDeleteHistoryUsecase(hairHistoryRepo domain.HairHistoryRepository) *DeleteHistoryUsecase {
	return &DeleteHistoryUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *DeleteHistoryUsecase) Execute(ctx context.Context, historyID string, req DeleteHistoryRequest) (DeleteHistoryResponse, error) {
	if err := u.hairHistoryRepo.Delete(ctx, historyID); err != nil {
		return DeleteHistoryResponse{}, err
	}
	return DeleteHistoryResponse{OK: true}, nil
}

