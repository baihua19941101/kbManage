package backuprestore

import "context"

func (s *Service) PublishDrillReport(ctx context.Context, userID, recordID uint64) (any, error) {
	return s.GenerateDrillReport(ctx, userID, recordID)
}
