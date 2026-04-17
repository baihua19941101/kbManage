package compliance

import (
	"context"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type EvidenceService struct {
	findingRepo  *repository.ComplianceFindingRepository
	evidenceRepo *repository.ComplianceEvidenceRepository
	scanRepo     *repository.ComplianceScanExecutionRepository
	scope        *ScopeService
}

func NewEvidenceService(findingRepo *repository.ComplianceFindingRepository, evidenceRepo *repository.ComplianceEvidenceRepository, scanRepo *repository.ComplianceScanExecutionRepository, scope *ScopeService) *EvidenceService {
	return &EvidenceService{findingRepo: findingRepo, evidenceRepo: evidenceRepo, scanRepo: scanRepo, scope: scope}
}

func (s *EvidenceService) ListByFinding(ctx context.Context, userID, findingID uint64) ([]domain.EvidenceRecord, error) {
	if s == nil || s.findingRepo == nil || s.evidenceRepo == nil || s.scanRepo == nil {
		return nil, ErrComplianceNotConfigured
	}
	finding, err := s.findingRepo.GetByID(ctx, findingID)
	if err != nil {
		return nil, err
	}
	execution, err := s.scanRepo.GetByID(ctx, finding.ScanExecutionID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateExecutionScope(ctx, userID, execution, PermissionComplianceRead); err != nil {
		return nil, err
	}
	allowRaw := s.scope.CanViewRawEvidence(ctx, userID, derefUint64(execution.WorkspaceID), derefUint64(execution.ProjectID))
	items, err := s.evidenceRepo.ListByFindingID(ctx, findingID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].PayloadJSON = s.scope.FilterEvidencePayload(items[i].PayloadJSON, allowRaw)
	}
	return items, nil
}
