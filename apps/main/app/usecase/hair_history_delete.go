package usecase

import (
	"context"

	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/response"
)

func (u HairHistory) Delete(ctx context.Context, r *request.DeleteHistory) (*response.DeleteHistory, error) {
	if err := u.hairHistorySvc.Delete(ctx, r.HistoryID); err != nil {
		return nil, err
	}
	return response.NewDeleteHistory(true), nil
}
