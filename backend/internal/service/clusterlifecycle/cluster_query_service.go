package clusterlifecycle

import "context"

func (s *Service) ListClusterRecords(ctx context.Context, userID uint64, filter ClusterListFilter) ([]ClusterListView, error) {
	items, err := s.ListClusters(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	out := make([]ClusterListView, 0, len(items))
	for _, item := range items {
		out = append(out, ClusterListView{
			ID:                 item.ID,
			Name:               item.Name,
			DisplayName:        item.DisplayName,
			Status:             string(item.Status),
			RegistrationStatus: string(item.RegistrationStatus),
			HealthStatus:       string(item.HealthStatus),
			KubernetesVersion:  item.KubernetesVersion,
			FailureReason:      item.RetirementReason,
		})
	}
	return out, nil
}

type ClusterListView struct {
	ID                 uint64 `json:"id"`
	Name               string `json:"name"`
	DisplayName        string `json:"displayName"`
	Status             string `json:"status"`
	RegistrationStatus string `json:"registrationStatus"`
	HealthStatus       string `json:"healthStatus"`
	KubernetesVersion  string `json:"kubernetesVersion"`
	FailureReason      string `json:"failureReason,omitempty"`
}
