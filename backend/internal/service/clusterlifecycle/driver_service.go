package clusterlifecycle

import (
	"context"

	"kbmanage/backend/internal/domain"
)

func (s *Service) CreateDriver(ctx context.Context, userID uint64, input CreateDriverInput) (*domain.ClusterDriverVersion, error) {
	return s.UpsertDriver(ctx, userID, input)
}

func (s *Service) ListCapabilitiesByDriver(ctx context.Context, userID, driverID uint64) ([]domain.CapabilityMatrixEntry, error) {
	return s.ListCapabilities(ctx, userID, driverID)
}
