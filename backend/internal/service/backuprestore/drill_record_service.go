package backuprestore

import "context"

func (s *Service) StartPlan(ctx context.Context, userID, planID uint64) (any, error) {
	return s.RunDrillPlan(ctx, userID, planID)
}
