package backuprestore

import "context"

func (s *Service) CreateMigration(ctx context.Context, userID uint64, input CreateMigrationPlanInput) (any, error) {
	return s.CreateMigrationPlan(ctx, userID, input)
}
