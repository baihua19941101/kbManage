package clusterlifecycle

import "context"

func (s *Service) IssueRegistrationBundle(ctx context.Context, userID uint64, input RegisterClusterInput) (*RegistrationBundle, error) {
	return s.RegisterCluster(ctx, userID, input)
}
