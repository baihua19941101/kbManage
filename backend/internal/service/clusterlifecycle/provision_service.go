package clusterlifecycle

import "context"

func (s *Service) ProvisionCluster(ctx context.Context, userID uint64, input CreateClusterInput) (*ClusterProvisionResult, error) {
	op, cluster, err := s.CreateCluster(ctx, userID, input)
	if err != nil {
		return nil, err
	}
	return &ClusterProvisionResult{Operation: op, Cluster: cluster}, nil
}

type ClusterProvisionResult struct {
	Operation any `json:"operation"`
	Cluster   any `json:"cluster"`
}
