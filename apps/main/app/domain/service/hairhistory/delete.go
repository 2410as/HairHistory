package hairhistory

import (
	"context"
)

func (s *service) Delete(ctx context.Context, historyID string) error {
	return s.hairHistoryRepo.Delete(ctx, historyID)
}
