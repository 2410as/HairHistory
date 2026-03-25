package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type DeleteHistoryResponse struct {
	OK bool `json:"ok"`
}

type DeleteHistoryUsecase struct {
	hairHistoryRepo domain.HairHistoryRepository
}

func NewDeleteHistoryUsecase(hairHistoryRepo domain.HairHistoryRepository) *DeleteHistoryUsecase {
	return &DeleteHistoryUsecase{hairHistoryRepo: hairHistoryRepo}
}

func (u *DeleteHistoryUsecase) Execute(ctx context.Context, req *request.DeleteHistory) (DeleteHistoryResponse, error) {
	if err := u.hairHistoryRepo.Delete(ctx, req.HistoryID); err != nil {
		return DeleteHistoryResponse{}, err
	}
	return DeleteHistoryResponse{OK: true}, nil
}

