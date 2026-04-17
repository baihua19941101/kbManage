package compliance

import (
	"context"
	"errors"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type FindingListFilter struct {
	WorkspaceID     uint64
	ProjectID       uint64
	ScanExecutionID uint64
	Result          string
	RiskLevel       string
}

type FindingService struct {
	repo     *repository.ComplianceFindingRepository
	scanRepo *repository.ComplianceScanExecutionRepository
	scope    *ScopeService
}

func NewFindingService(repo *repository.ComplianceFindingRepository, scanRepo *repository.ComplianceScanExecutionRepository, scope *ScopeService) *FindingService {
	return &FindingService{repo: repo, scanRepo: scanRepo, scope: scope}
}

func (s *FindingService) List(ctx context.Context, userID uint64, filter FindingListFilter) ([]domain.ComplianceFinding, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	if err := s.scope.ValidateScope(ctx, userID, filter.WorkspaceID, filter.ProjectID, PermissionComplianceRead); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, repository.ComplianceFindingListFilter{
		WorkspaceID:     uint64Ptr(filter.WorkspaceID),
		ProjectID:       uint64Ptr(filter.ProjectID),
		ScanExecutionID: uint64Ptr(filter.ScanExecutionID),
		Result:          domain.ComplianceFindingResult(filter.Result),
		RiskLevel:       domain.ComplianceRiskLevel(filter.RiskLevel),
	})
}

func (s *FindingService) Get(ctx context.Context, userID, findingID uint64) (*domain.ComplianceFinding, error) {
	if s == nil || s.repo == nil || s.scanRepo == nil {
		return nil, ErrComplianceNotConfigured
	}
	item, err := s.repo.GetByID(ctx, findingID)
	if err != nil {
		return nil, err
	}
	execution, err := s.scanRepo.GetByID(ctx, item.ScanExecutionID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateExecutionScope(ctx, userID, execution, PermissionComplianceRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *FindingService) GetExecution(ctx context.Context, userID, findingID uint64) (*domain.ScanExecution, error) {
	item, err := s.Get(ctx, userID, findingID)
	if err != nil {
		return nil, err
	}
	if item.ScanExecutionID == 0 {
		return nil, errors.New("scan execution id is required")
	}
	return s.scanRepo.GetByID(ctx, item.ScanExecutionID)
}
