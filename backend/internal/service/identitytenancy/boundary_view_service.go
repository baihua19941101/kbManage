package identitytenancy

import (
	"context"
)

func (s *Service) BuildBoundarySummary(ctx context.Context, unitID uint64) (map[string]any, error) {
	mappings, err := s.mappings.ListByUnitID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"unitId":       unitID,
		"mappingCount": len(mappings),
		"hasConflict":  false,
	}, nil
}
