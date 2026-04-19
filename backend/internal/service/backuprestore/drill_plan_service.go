package backuprestore

import "context"

func (s *Service) CreatePlan(ctx context.Context, userID uint64, input CreateDRDrillPlanInput) (any, error) {
	return s.CreateDrillPlan(ctx, userID, input)
}
