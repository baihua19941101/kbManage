package clusterlifecycle

import "context"

func (s *Service) ListClusterNodePools(ctx context.Context, userID, clusterID uint64) ([]NodePoolView, error) {
	items, err := s.ListNodePools(ctx, userID, clusterID)
	if err != nil {
		return nil, err
	}
	out := make([]NodePoolView, 0, len(items))
	for _, item := range items {
		out = append(out, NodePoolView{
			ID:           item.ID,
			Name:         item.Name,
			DesiredCount: item.DesiredCount,
			CurrentCount: item.CurrentCount,
			Status:       string(item.Status),
		})
	}
	return out, nil
}

type NodePoolView struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	DesiredCount int    `json:"desiredCount"`
	CurrentCount int    `json:"currentCount"`
	Status       string `json:"status"`
}
