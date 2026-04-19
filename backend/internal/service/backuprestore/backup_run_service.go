package backuprestore

import "context"

func (s *Service) RunBackupPolicy(ctx context.Context, userID, policyID uint64) (any, error) {
	return s.RunPolicy(ctx, userID, policyID)
}
