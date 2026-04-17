package compliance

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"
	authSvc "kbmanage/backend/internal/service/auth"
)

type ScopeFilter struct {
	WorkspaceID uint64
	ProjectID   uint64
	ClusterIDs  []uint64
}

type ScopeService struct {
	scopeAccess *authSvc.ScopeAccessService
}

func NewScopeService(scopeAccess *authSvc.ScopeAccessService) *ScopeService {
	return &ScopeService{scopeAccess: scopeAccess}
}

func (s *ScopeService) ValidateScope(ctx context.Context, userID, workspaceID, projectID uint64, permission string) error {
	if s == nil || s.scopeAccess == nil {
		return nil
	}
	if userID == 0 {
		return ErrComplianceScopeDenied
	}
	mapped := s.mapPermission(permission)
	if projectID != 0 {
		allowed, err := s.scopeAccess.HasScopePermission(ctx, userID, domain.ScopeTypeProject, workspaceID, projectID, mapped)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrComplianceScopeDenied
		}
		return nil
	}
	if workspaceID != 0 {
		allowed, err := s.scopeAccess.HasScopePermission(ctx, userID, domain.ScopeTypeWorkspace, workspaceID, 0, mapped)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrComplianceScopeDenied
		}
	}
	return nil
}

func (s *ScopeService) ValidateBaselineScope(ctx context.Context, userID uint64, baseline *domain.ComplianceBaseline, permission string) error {
	if baseline == nil {
		return errors.New("baseline is required")
	}
	return s.ValidateScope(ctx, userID, 0, 0, permission)
}

func (s *ScopeService) ValidateProfileScope(ctx context.Context, userID uint64, profile *domain.ScanProfile, permission string) error {
	if profile == nil {
		return errors.New("scan profile is required")
	}
	return s.ValidateScope(ctx, userID, derefUint64(profile.WorkspaceID), derefUint64(profile.ProjectID), permission)
}

func (s *ScopeService) ValidateExecutionScope(ctx context.Context, userID uint64, execution *domain.ScanExecution, permission string) error {
	if execution == nil {
		return errors.New("scan execution is required")
	}
	return s.ValidateScope(ctx, userID, derefUint64(execution.WorkspaceID), derefUint64(execution.ProjectID), permission)
}

func (s *ScopeService) FilterClustersByReadScope(ctx context.Context, userID uint64, clusterIDs []uint64) ([]uint64, error) {
	if s == nil || s.scopeAccess == nil {
		return clusterIDs, nil
	}
	allowedIDs, constrained, err := s.scopeAccess.ListClusterIDsByPermission(ctx, userID, s.mapPermission(PermissionComplianceRead))
	if err != nil {
		return nil, err
	}
	if !constrained || len(clusterIDs) == 0 {
		return clusterIDs, nil
	}
	allowed := make(map[uint64]struct{}, len(allowedIDs))
	for _, id := range allowedIDs {
		allowed[id] = struct{}{}
	}
	filtered := make([]uint64, 0, len(clusterIDs))
	for _, id := range clusterIDs {
		if _, ok := allowed[id]; ok {
			filtered = append(filtered, id)
		}
	}
	return filtered, nil
}

func (s *ScopeService) FilterScanProfiles(ctx context.Context, userID uint64, items []domain.ScanProfile, permission string) []domain.ScanProfile {
	filtered := make([]domain.ScanProfile, 0, len(items))
	for i := range items {
		if err := s.ValidateProfileScope(ctx, userID, &items[i], permission); err == nil {
			filtered = append(filtered, items[i])
		}
	}
	return filtered
}

func (s *ScopeService) FilterScanExecutions(ctx context.Context, userID uint64, items []domain.ScanExecution, permission string) []domain.ScanExecution {
	filtered := make([]domain.ScanExecution, 0, len(items))
	for i := range items {
		if err := s.ValidateExecutionScope(ctx, userID, &items[i], permission); err == nil {
			filtered = append(filtered, items[i])
		}
	}
	return filtered
}

func (s *ScopeService) CanViewRawEvidence(ctx context.Context, userID uint64, workspaceID, projectID uint64) bool {
	if err := s.ValidateScope(ctx, userID, workspaceID, projectID, PermissionComplianceExecuteScan); err == nil {
		return true
	}
	if err := s.ValidateScope(ctx, userID, workspaceID, projectID, PermissionComplianceManageBaseline); err == nil {
		return true
	}
	return false
}

func (s *ScopeService) FilterEvidencePayload(raw string, allowRaw bool) string {
	if allowRaw || strings.TrimSpace(raw) == "" {
		return raw
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return ""
	}
	masked, _ := marshalJSON(maskEvidencePayload(payload))
	return masked
}

func (s *ScopeService) mapPermission(permission string) string {
	switch strings.TrimSpace(permission) {
	case PermissionComplianceManageBaseline, PermissionComplianceManageRemediation, PermissionComplianceReviewException:
		return "securitypolicy:manage"
	case PermissionComplianceExecuteScan:
		return "securitypolicy:enforce"
	case PermissionComplianceExportArchive:
		return "audit:read"
	case PermissionComplianceRead, "":
		return "securitypolicy:read"
	default:
		return permission
	}
}
