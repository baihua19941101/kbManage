package clusterlifecycle

import "context"

func (s *Service) ValidateCluster(ctx context.Context, userID, clusterID uint64, input ValidationInput) (*ValidationResult, error) {
	return s.ValidateChange(ctx, userID, clusterID, input)
}
