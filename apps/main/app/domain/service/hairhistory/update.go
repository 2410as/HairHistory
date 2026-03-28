package hairhistory

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/domain"
	"github.com/annasakai/hairhistorymemo/apps/main/app/domain/entity"
)

func (s *service) Update(ctx context.Context, historyID string, req entity.UpdateHairHistoryParams) (*entity.HairHistory, error) {
	h, err := s.hairHistoryRepo.Update(ctx, historyID, req)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, domain.ErrNotFound
	}
	return h, nil
}
