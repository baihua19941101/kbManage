package compliance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	baselineProvider "kbmanage/backend/internal/integration/compliance/baseline"
	scannerProvider "kbmanage/backend/internal/integration/compliance/scanner"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"
)

const (
	PermissionComplianceRead              = "compliance:read"
	PermissionComplianceManageBaseline    = "compliance:manage-baseline"
	PermissionComplianceExecuteScan       = "compliance:execute-scan"
	PermissionComplianceManageRemediation = "compliance:manage-remediation"
	PermissionComplianceReviewException   = "compliance:review-exception"
	PermissionComplianceExportArchive     = "compliance:export-archive"
)

var (
	ErrComplianceScopeDenied = errors.New("compliance scope access denied")
)

type Service struct {
	Baselines *BaselineService
	Profiles  *ScanProfileService
	Scans     *ScanExecutionService
	Findings  *FindingService
	Evidence  *EvidenceService
	Scope     *ScopeService
}

func NewService(
	baselineRepo *repository.ComplianceBaselineRepository,
	profileRepo *repository.ComplianceScanProfileRepository,
	scanRepo *repository.ComplianceScanExecutionRepository,
	findingRepo *repository.ComplianceFindingRepository,
	evidenceRepo *repository.ComplianceEvidenceRepository,
	remediationRepo *repository.ComplianceRemediationRepository,
	exceptionRepo *repository.ComplianceExceptionRepository,
	recheckRepo *repository.ComplianceRecheckRepository,
	trendRepo *repository.ComplianceTrendRepository,
	exportRepo *repository.ComplianceExportRepository,
	scopeAccess *authSvc.ScopeAccessService,
	baselineSnapshotProvider baselineProvider.Provider,
	scanner scannerProvider.Provider,
	progressCache *ProgressCache,
	exportCache *ExportCache,
	scheduleCache *ScheduleCache,
) *Service {
	scopeSvc := NewScopeService(scopeAccess)
	baselines := NewBaselineService(baselineRepo, scopeSvc, baselineSnapshotProvider)
	profiles := NewScanProfileService(profileRepo, baselineRepo, scopeSvc, scheduleCache)
	scans := NewScanExecutionService(scanRepo, baselineRepo, profileRepo, findingRepo, evidenceRepo, scopeSvc, scanner, baselineSnapshotProvider, progressCache)
	findings := NewFindingService(findingRepo, scanRepo, scopeSvc)
	evidence := NewEvidenceService(findingRepo, evidenceRepo, scanRepo, scopeSvc)
	_ = remediationRepo
	_ = exceptionRepo
	_ = recheckRepo
	_ = trendRepo
	_ = exportRepo
	_ = exportCache
	return &Service{Baselines: baselines, Profiles: profiles, Scans: scans, Findings: findings, Evidence: evidence, Scope: scopeSvc}
}

func uint64Ptr(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}

func derefUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func marshalJSON(v any) (string, error) {
	if v == nil {
		return "", nil
	}
	payload, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func unmarshalStringSlice(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return []string{}
	}
	return out
}

func unmarshalUint64Slice(raw string) []uint64 {
	if strings.TrimSpace(raw) == "" {
		return []uint64{}
	}
	var out []uint64
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return []uint64{}
	}
	return out
}

func unmarshalNodeSelectors(raw string) []map[string]string {
	if strings.TrimSpace(raw) == "" {
		return []map[string]string{}
	}
	var out []map[string]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return []map[string]string{}
	}
	return out
}

func parseRFC3339Optional(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("invalid RFC3339 time: %w", err)
	}
	return &parsed, nil
}

func containsString(values []string, target string) bool {
	for _, item := range values {
		if strings.EqualFold(strings.TrimSpace(item), strings.TrimSpace(target)) {
			return true
		}
	}
	return false
}

func sortAndCompactStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func cloneMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func maskEvidencePayload(payload map[string]any) map[string]any {
	if len(payload) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(payload))
	for k := range payload {
		out[k] = "***"
	}
	return out
}

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
