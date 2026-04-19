package backuprestore

import "context"

func (s *Service) SubmitRestore(ctx context.Context, userID uint64, input CreateRestoreJobInput) (any, error) {
	return s.CreateRestoreJob(ctx, userID, input)
}
