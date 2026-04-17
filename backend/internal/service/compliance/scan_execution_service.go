package compliance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	baselineProvider "kbmanage/backend/internal/integration/compliance/baseline"
	scannerProvider "kbmanage/backend/internal/integration/compliance/scanner"
	"kbmanage/backend/internal/repository"
)

type ScanExecutionListFilter struct {
	WorkspaceID   uint64
	ProjectID     uint64
	ProfileID     uint64
	Status        string
	TriggerSource string
}

type ExecuteScanInput struct {
	TriggerSource string `json:"triggerSource"`
}

type ScanExecutionService struct {
	repo             *repository.ComplianceScanExecutionRepository
	baselineRepo     *repository.ComplianceBaselineRepository
	profileRepo      *repository.ComplianceScanProfileRepository
	findingRepo      *repository.ComplianceFindingRepository
	evidenceRepo     *repository.ComplianceEvidenceRepository
	scope            *ScopeService
	scanner          scannerProvider.Provider
	snapshotProvider baselineProvider.Provider
	progressCache    *ProgressCache
}

func NewScanExecutionService(repo *repository.ComplianceScanExecutionRepository, baselineRepo *repository.ComplianceBaselineRepository, profileRepo *repository.ComplianceScanProfileRepository, findingRepo *repository.ComplianceFindingRepository, evidenceRepo *repository.ComplianceEvidenceRepository, scope *ScopeService, scanner scannerProvider.Provider, snapshotProvider baselineProvider.Provider, progressCache *ProgressCache) *ScanExecutionService {
	if scanner == nil {
		scanner = scannerProvider.NewMockProvider()
	}
	if snapshotProvider == nil {
		snapshotProvider = baselineProvider.NewStaticProvider()
	}
	return &ScanExecutionService{repo: repo, baselineRepo: baselineRepo, profileRepo: profileRepo, findingRepo: findingRepo, evidenceRepo: evidenceRepo, scope: scope, scanner: scanner, snapshotProvider: snapshotProvider, progressCache: progressCache}
}

func (s *ScanExecutionService) List(ctx context.Context, userID uint64, filter ScanExecutionListFilter) ([]domain.ScanExecution, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	if err := s.scope.ValidateScope(ctx, userID, filter.WorkspaceID, filter.ProjectID, PermissionComplianceRead); err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, repository.ComplianceScanExecutionListFilter{
		WorkspaceID:   uint64Ptr(filter.WorkspaceID),
		ProjectID:     uint64Ptr(filter.ProjectID),
		ProfileID:     uint64Ptr(filter.ProfileID),
		Status:        domain.ComplianceScanStatus(strings.TrimSpace(filter.Status)),
		TriggerSource: domain.ComplianceTriggerSource(strings.TrimSpace(filter.TriggerSource)),
	})
	if err != nil {
		return nil, err
	}
	return s.scope.FilterScanExecutions(ctx, userID, items, PermissionComplianceRead), nil
}

func (s *ScanExecutionService) Get(ctx context.Context, userID, executionID uint64) (*domain.ScanExecution, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	item, err := s.repo.GetByID(ctx, executionID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateExecutionScope(ctx, userID, item, PermissionComplianceRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ScanExecutionService) ExecuteNow(ctx context.Context, userID, profileID uint64, input ExecuteScanInput) (*domain.ScanExecution, error) {
	if s == nil || s.repo == nil || s.profileRepo == nil || s.baselineRepo == nil || s.findingRepo == nil || s.evidenceRepo == nil {
		return nil, ErrComplianceNotConfigured
	}
	profile, err := s.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateProfileScope(ctx, userID, profile, PermissionComplianceExecuteScan); err != nil {
		return nil, err
	}
	baseline, err := s.baselineRepo.GetByID(ctx, profile.BaselineID)
	if err != nil {
		return nil, err
	}
	snapshot, err := s.snapshotProvider.BuildSnapshot(ctx, baseline)
	if err != nil {
		return nil, err
	}
	snapshotJSON, err := marshalJSON(snapshot)
	if err != nil {
		return nil, err
	}
	scopeSnapshotJSON, err := marshalJSON(map[string]any{
		"scopeType":     profile.ScopeType,
		"clusterIds":    unmarshalUint64Slice(profile.ClusterRefsJSON),
		"nodeSelectors": unmarshalNodeSelectors(profile.NodeSelectorsJSON),
		"namespaces":    unmarshalStringSlice(profile.NamespaceRefsJSON),
		"resourceKinds": unmarshalStringSlice(profile.ResourceKindsJSON),
	})
	if err != nil {
		return nil, err
	}
	triggerSource := strings.TrimSpace(input.TriggerSource)
	if triggerSource == "" {
		triggerSource = string(domain.ComplianceTriggerSourceManual)
	}
	item := &domain.ScanExecution{
		ProfileID:            profile.ID,
		BaselineID:           baseline.ID,
		WorkspaceID:          profile.WorkspaceID,
		ProjectID:            profile.ProjectID,
		BaselineSnapshotJSON: snapshotJSON,
		BaselineVersionLabel: fmt.Sprintf("%s@%s", baseline.StandardType, baseline.Version),
		TriggerSource:        domain.ComplianceTriggerSource(triggerSource),
		Status:               domain.ComplianceScanStatusRunning,
		CoverageStatus:       domain.ComplianceCoverageStatusUnavailable,
		ScopeSnapshotJSON:    scopeSnapshotJSON,
		CreatedBy:            uint64Ptr(userID),
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.progressCache.Set(ctx, item.ID, ScanProgressSnapshot{Status: string(domain.ComplianceScanStatusRunning), Progress: 10, Message: "scan started"})
	result, err := s.scanner.Execute(ctx, scannerProvider.Request{ExecutionID: item.ID, Profile: *profile, Baseline: snapshot})
	if err != nil {
		_ = s.repo.UpdateFields(ctx, item.ID, map[string]any{"status": domain.ComplianceScanStatusFailed, "error_summary": err.Error(), "completed_at": time.Now().UTC()})
		return s.repo.GetByID(ctx, item.ID)
	}
	now := time.Now().UTC()
	findings := make([]domain.ComplianceFinding, 0, len(result.Findings))
	for _, candidate := range result.Findings {
		findings = append(findings, domain.ComplianceFinding{
			ScanExecutionID:   item.ID,
			ControlID:         candidate.ControlID,
			ControlTitle:      candidate.ControlTitle,
			Result:            candidate.Result,
			RiskLevel:         candidate.RiskLevel,
			ClusterID:         candidate.ClusterID,
			NodeName:          candidate.NodeName,
			Namespace:         candidate.Namespace,
			ResourceKind:      candidate.ResourceKind,
			ResourceName:      candidate.ResourceName,
			ResourceUID:       candidate.ResourceUID,
			Summary:           candidate.Summary,
			RemediationStatus: domain.ComplianceRemediationStatusOpen,
			DetectedAt:        now,
		})
	}
	if err := s.findingRepo.ReplaceByScanExecution(ctx, item.ID, findings); err != nil {
		return nil, err
	}
	storedFindings, err := s.findingRepo.List(ctx, repository.ComplianceFindingListFilter{ScanExecutionID: &item.ID})
	if err != nil {
		return nil, err
	}
	byControl := make(map[string]uint64, len(storedFindings))
	for _, finding := range storedFindings {
		byControl[finding.ControlID] = finding.ID
	}
	evidences := make([]domain.EvidenceRecord, 0)
	for _, candidate := range result.Findings {
		findingID := byControl[candidate.ControlID]
		for _, evidence := range candidate.Evidences {
			payloadJSON, _ := marshalJSON(evidence.Payload)
			evidences = append(evidences, domain.EvidenceRecord{
				FindingID:       findingID,
				EvidenceType:    evidence.EvidenceType,
				SourceRef:       evidence.SourceRef,
				CollectedAt:     now,
				Confidence:      evidence.Confidence,
				Summary:         evidence.Summary,
				ArtifactRef:     evidence.ArtifactRef,
				RedactionStatus: evidence.RedactionStatus,
				PayloadJSON:     payloadJSON,
			})
		}
	}
	if err := s.evidenceRepo.CreateBatch(ctx, evidences); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateFields(ctx, item.ID, map[string]any{
		"status":          result.Status,
		"coverage_status": result.CoverageStatus,
		"started_at":      result.StartedAt,
		"completed_at":    result.CompletedAt,
		"score":           result.Score,
		"pass_count":      result.PassCount,
		"fail_count":      result.FailCount,
		"warning_count":   result.WarningCount,
		"error_summary":   result.ErrorSummary,
	}); err != nil {
		return nil, err
	}
	_ = s.progressCache.Set(ctx, item.ID, ScanProgressSnapshot{Status: string(result.Status), Progress: 100, Message: "scan finished"})
	return s.repo.GetByID(ctx, item.ID)
}

func (s *ScanExecutionService) RunPending(ctx context.Context, limit int) (int, error) {
	if s == nil || s.repo == nil {
		return 0, ErrComplianceNotConfigured
	}
	items, err := s.repo.ListPending(ctx, limit)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range items {
		if _, err := s.ExecuteNow(ctx, derefUint64(item.CreatedBy), item.ProfileID, ExecuteScanInput{TriggerSource: string(item.TriggerSource)}); err == nil {
			count++
		}
	}
	return count, nil
}
