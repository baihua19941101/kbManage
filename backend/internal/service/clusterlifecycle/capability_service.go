package clusterlifecycle

import "context"

func (s *Service) ListDriverCapabilityMatrix(ctx context.Context, userID, driverID uint64) ([]CapabilityMatrixView, error) {
	items, err := s.ListCapabilities(ctx, userID, driverID)
	if err != nil {
		return nil, err
	}
	out := make([]CapabilityMatrixView, 0, len(items))
	for _, item := range items {
		out = append(out, CapabilityMatrixView{
			ID:                  item.ID,
			CapabilityDomain:    item.CapabilityDomain,
			SupportLevel:        string(item.SupportLevel),
			CompatibilityStatus: string(item.CompatibilityStatus),
		})
	}
	return out, nil
}

type CapabilityMatrixView struct {
	ID                  uint64 `json:"id"`
	CapabilityDomain    string `json:"capabilityDomain"`
	SupportLevel        string `json:"supportLevel"`
	CompatibilityStatus string `json:"compatibilityStatus"`
}
