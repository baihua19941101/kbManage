package clusterlifecycle

import "context"

func (s *Service) RequestClusterDisable(ctx context.Context, userID, clusterID uint64, input DisableClusterInput) error {
	_, err := s.DisableCluster(ctx, userID, clusterID, input)
	return err
}

func (s *Service) RequestClusterRetirement(ctx context.Context, userID, clusterID uint64, input RetireClusterInput) error {
	_, err := s.RetireCluster(ctx, userID, clusterID, input)
	return err
}
