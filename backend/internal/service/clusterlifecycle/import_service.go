package clusterlifecycle

import "context"

func (s *Service) SubmitImportCluster(ctx context.Context, userID uint64, input ImportClusterInput) (*ClusterImportResult, error) {
	op, cluster, err := s.ImportCluster(ctx, userID, input)
	if err != nil {
		return nil, err
	}
	return &ClusterImportResult{Operation: op, Cluster: cluster}, nil
}

type ClusterImportResult struct {
	Operation any `json:"operation"`
	Cluster   any `json:"cluster"`
}
