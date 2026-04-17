package compliance

import (
	"context"
	"sort"
	"strings"
	"time"
)

type ComplianceTrendPoint struct {
	WindowStart               time.Time `json:"windowStart"`
	WindowEnd                 time.Time `json:"windowEnd"`
	ScoreAvg                  float64   `json:"scoreAvg"`
	CoverageRate              float64   `json:"coverageRate"`
	RemediationCompletionRate float64   `json:"remediationCompletionRate"`
	HighRiskOpenCount         int       `json:"highRiskOpenCount"`
	BaselineVersion           string    `json:"baselineVersion,omitempty"`
	WorkspaceID               uint64    `json:"workspaceId,omitempty"`
	ProjectID                 uint64    `json:"projectId,omitempty"`
	ScopeType                 string    `json:"scopeType,omitempty"`
	ScopeRef                  string    `json:"scopeRef,omitempty"`
	BaselineID                string    `json:"baselineId,omitempty"`
	GeneratedAt               time.Time `json:"generatedAt"`
}

type TrendComparisonBasis struct {
	BaselineID       string   `json:"baselineId,omitempty"`
	BaselineVersions []string `json:"baselineVersions,omitempty"`
	MixedBaseline    bool     `json:"mixedBaselineFlag"`
}

type ComplianceTrendResponse struct {
	Points          []ComplianceTrendPoint `json:"points"`
	ComparisonBasis TrendComparisonBasis   `json:"comparisonBasis"`
}

type TrendFilter struct {
	WorkspaceID uint64
	ProjectID   uint64
	BaselineID  string
	ScopeType   string
	ScopeRef    string
	TimeFrom    *time.Time
	TimeTo      *time.Time
}

type TrendService struct {
	store    *complianceStore
	overview *OverviewService
	now      func() time.Time
}

func NewTrendService() *TrendService {
	return &TrendService{store: defaultComplianceStore, overview: NewOverviewService(), now: time.Now}
}

func (s *TrendService) RecordSnapshot(ctx context.Context, filter TrendFilter) (*ComplianceTrendPoint, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	now := s.now()
	windowEnd := now.UTC().Truncate(time.Hour)
	windowStart := windowEnd.Add(-24 * time.Hour)
	if filter.TimeTo != nil {
		windowEnd = filter.TimeTo.UTC()
	}
	if filter.TimeFrom != nil {
		windowStart = filter.TimeFrom.UTC()
	}
	overview, err := s.overview.GetOverview(ctx, OverviewFilter{WorkspaceID: filter.WorkspaceID, ProjectID: filter.ProjectID, GroupBy: "cluster", TimeFrom: filter.TimeFrom, TimeTo: filter.TimeTo})
	if err != nil {
		return nil, err
	}
	baselineVersion := latestBaselineVersionLocked(s.store, strings.TrimSpace(filter.BaselineID))
	point := &ComplianceTrendPoint{
		WindowStart:               windowStart,
		WindowEnd:                 windowEnd,
		ScoreAvg:                  maxFloat(0, 100-float64(overview.OpenFindingsCount*10)),
		CoverageRate:              overview.CoverageRate,
		RemediationCompletionRate: overview.RemediationCompletionRate,
		HighRiskOpenCount:         overview.HighRiskOpenCount,
		BaselineVersion:           baselineVersion,
		WorkspaceID:               filter.WorkspaceID,
		ProjectID:                 filter.ProjectID,
		ScopeType:                 normalizeTrendScopeType(filter.ScopeType),
		ScopeRef:                  strings.TrimSpace(filter.ScopeRef),
		BaselineID:                strings.TrimSpace(filter.BaselineID),
		GeneratedAt:               now,
	}
	s.store.mu.Lock()
	s.store.snapshots = append(s.store.snapshots, point)
	s.store.mu.Unlock()
	return cloneTrendPoint(point), nil
}

func (s *TrendService) GetTrends(_ context.Context, filter TrendFilter) (*ComplianceTrendResponse, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	points := make([]ComplianceTrendPoint, 0, len(s.store.snapshots))
	versions := make(map[string]struct{})
	for _, point := range s.store.snapshots {
		if !matchTrendFilter(point, filter) {
			continue
		}
		copyItem := cloneTrendPoint(point)
		points = append(points, *copyItem)
		if version := strings.TrimSpace(copyItem.BaselineVersion); version != "" {
			versions[version] = struct{}{}
		}
	}
	sort.Slice(points, func(i, j int) bool { return points[i].WindowStart.Before(points[j].WindowStart) })
	basis := TrendComparisonBasis{BaselineID: strings.TrimSpace(filter.BaselineID), BaselineVersions: make([]string, 0, len(versions))}
	for version := range versions {
		basis.BaselineVersions = append(basis.BaselineVersions, version)
	}
	sort.Strings(basis.BaselineVersions)
	basis.MixedBaseline = len(basis.BaselineVersions) > 1
	return &ComplianceTrendResponse{Points: points, ComparisonBasis: basis}, nil
}

func cloneTrendPoint(item *ComplianceTrendPoint) *ComplianceTrendPoint {
	if item == nil {
		return nil
	}
	copyItem := *item
	return &copyItem
}

func matchTrendFilter(point *ComplianceTrendPoint, filter TrendFilter) bool {
	if point == nil {
		return false
	}
	if filter.WorkspaceID != 0 && point.WorkspaceID != filter.WorkspaceID {
		return false
	}
	if filter.ProjectID != 0 && point.ProjectID != filter.ProjectID {
		return false
	}
	if filter.BaselineID != "" && point.BaselineID != strings.TrimSpace(filter.BaselineID) {
		return false
	}
	if filter.ScopeType != "" && !strings.EqualFold(point.ScopeType, filter.ScopeType) {
		return false
	}
	if filter.ScopeRef != "" && point.ScopeRef != strings.TrimSpace(filter.ScopeRef) {
		return false
	}
	if filter.TimeFrom != nil && point.WindowStart.Before(*filter.TimeFrom) {
		return false
	}
	if filter.TimeTo != nil && point.WindowEnd.After(*filter.TimeTo) {
		return false
	}
	return true
}

func normalizeTrendScopeType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "cluster", "workspace", "project":
		return strings.ToLower(strings.TrimSpace(raw))
	default:
		return "cluster"
	}
}

func latestBaselineVersionLocked(store *complianceStore, baselineID string) string {
	if strings.TrimSpace(baselineID) == "" {
		return "mixed"
	}
	for _, point := range store.snapshots {
		if point.BaselineID == baselineID && strings.TrimSpace(point.BaselineVersion) != "" {
			return point.BaselineVersion
		}
	}
	return "v1"
}
