package backuprestore

import "context"

func (s *Service) ValidateRestore(ctx context.Context, userID uint64, input CreateRestoreJobInput) (*PrecheckResult, error) {
	return s.ValidateRestoreJob(ctx, userID, input)
}
