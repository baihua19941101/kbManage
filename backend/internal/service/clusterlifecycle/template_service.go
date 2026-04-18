package clusterlifecycle

import "context"

func (s *Service) ListTemplateViews(ctx context.Context, userID uint64, driverKey, infrastructureType string) ([]TemplateView, error) {
	items, err := s.ListTemplates(ctx, userID, driverKey, infrastructureType)
	if err != nil {
		return nil, err
	}
	out := make([]TemplateView, 0, len(items))
	for _, item := range items {
		out = append(out, TemplateView{
			ID:                 item.ID,
			Name:               item.Name,
			DriverKey:          item.DriverKey,
			InfrastructureType: item.InfrastructureType,
			Status:             string(item.Status),
		})
	}
	return out, nil
}

type TemplateView struct {
	ID                 uint64 `json:"id"`
	Name               string `json:"name"`
	DriverKey          string `json:"driverKey"`
	InfrastructureType string `json:"infrastructureType"`
	Status             string `json:"status"`
}
