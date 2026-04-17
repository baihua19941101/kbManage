package compliance

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ComplianceOverview struct {
	CoverageRate              float64                   `json:"coverageRate"`
	OpenFindingsCount         int                       `json:"openFindingsCount"`
	HighRiskOpenCount         int                       `json:"highRiskOpenCount"`
	RemediationCompletionRate float64                   `json:"remediationCompletionRate"`
	Groups                    []ComplianceOverviewGroup `json:"groups"`
}

type ComplianceOverviewGroup struct {
	GroupKey                  string  `json:"groupKey"`
	ScoreAvg                  float64 `json:"scoreAvg"`
	CoverageRate              float64 `json:"coverageRate"`
	OpenFindingsCount         int     `json:"openFindingsCount"`
	HighRiskOpenCount         int     `json:"highRiskOpenCount"`
	RemediationCompletionRate float64 `json:"remediationCompletionRate"`
}

type OverviewFilter struct {
	WorkspaceID uint64
	ProjectID   uint64
	GroupBy     string
	TimeFrom    *time.Time
	TimeTo      *time.Time
}

type OverviewService struct {
	store *complianceStore
}

func NewOverviewService() *OverviewService {
	return &OverviewService{store: defaultComplianceStore}
}

func (s *OverviewService) GetOverview(_ context.Context, filter OverviewFilter) (*ComplianceOverview, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	groupBy := normalizeOverviewGroupBy(filter.GroupBy)
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()

	groups := make(map[string]*ComplianceOverviewGroup)
	total := 0
	openCount := 0
	highRiskOpen := 0
	totalTasks := 0
	doneTasks := 0
	coveredCount := 0

	for _, finding := range s.store.findings {
		if !matchFindingOverviewFilter(finding, filter) {
			continue
		}
		total++
		groupKey := overviewGroupKey(groupBy, finding)
		group := ensureOverviewGroup(groups, groupKey)
		if finding.RemediationStatus != "open" {
			coveredCount++
			group.CoverageRate++
		}
		if finding.RemediationStatus == "open" || finding.RemediationStatus == "in_progress" || finding.RemediationStatus == "ready_for_recheck" {
			openCount++
			group.OpenFindingsCount++
			if finding.RiskLevel == "high" || finding.RiskLevel == "critical" || finding.RiskLevel == "" {
				highRiskOpen++
				group.HighRiskOpenCount++
			}
		}
	}

	for _, task := range s.store.remediations {
		if !matchTaskOverviewFilter(task, filter) {
			continue
		}
		totalTasks++
		groupKey := overviewGroupKey(groupBy, &findingRecord{ClusterID: task.ClusterID, WorkspaceID: task.WorkspaceID, ProjectID: task.ProjectID, BaselineID: task.BaselineID})
		group := ensureOverviewGroup(groups, groupKey)
		if task.Status == "done" {
			doneTasks++
			group.RemediationCompletionRate++
		}
	}

	response := &ComplianceOverview{Groups: make([]ComplianceOverviewGroup, 0, len(groups))}
	if total > 0 {
		response.CoverageRate = float64(coveredCount) * 100 / float64(total)
		response.OpenFindingsCount = openCount
		response.HighRiskOpenCount = highRiskOpen
	}
	if totalTasks > 0 {
		response.RemediationCompletionRate = float64(doneTasks) * 100 / float64(totalTasks)
	}
	for _, group := range groups {
		memberCount := countFindingsForGroupLocked(s.store, filter, groupBy, group.GroupKey)
		if memberCount > 0 {
			group.CoverageRate = group.CoverageRate * 100 / memberCount
			group.ScoreAvg = maxFloat(0, 100-float64(group.OpenFindingsCount*10))
		}
		taskCount := countTasksForGroupLocked(s.store, filter, groupBy, group.GroupKey)
		if taskCount > 0 {
			group.RemediationCompletionRate = group.RemediationCompletionRate * 100 / taskCount
		}
		response.Groups = append(response.Groups, *group)
	}
	sort.Slice(response.Groups, func(i, j int) bool { return response.Groups[i].GroupKey < response.Groups[j].GroupKey })
	return response, nil
}

func normalizeOverviewGroupBy(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "cluster", "workspace", "project", "baseline":
		return strings.ToLower(strings.TrimSpace(raw))
	default:
		return "cluster"
	}
}

func overviewGroupKey(groupBy string, finding *findingRecord) string {
	if finding == nil {
		return "unknown"
	}
	switch groupBy {
	case "workspace":
		return formatUintKey("workspace", finding.WorkspaceID)
	case "project":
		return formatUintKey("project", finding.ProjectID)
	case "baseline":
		if id := strings.TrimSpace(finding.BaselineID); id != "" {
			return "baseline:" + id
		}
		return "baseline:unknown"
	default:
		return formatUintKey("cluster", finding.ClusterID)
	}
}

func ensureOverviewGroup(groups map[string]*ComplianceOverviewGroup, groupKey string) *ComplianceOverviewGroup {
	group := groups[groupKey]
	if group == nil {
		group = &ComplianceOverviewGroup{GroupKey: groupKey}
		groups[groupKey] = group
	}
	return group
}

func formatUintKey(prefix string, value uint64) string {
	if value == 0 {
		return prefix + ":unknown"
	}
	return prefix + ":" + strconv.FormatUint(value, 10)
}

func countFindingsForGroupLocked(store *complianceStore, filter OverviewFilter, groupBy, groupKey string) float64 {
	count := 0.0
	for _, finding := range store.findings {
		if !matchFindingOverviewFilter(finding, filter) {
			continue
		}
		if overviewGroupKey(groupBy, finding) == groupKey {
			count++
		}
	}
	return count
}

func countTasksForGroupLocked(store *complianceStore, filter OverviewFilter, groupBy, groupKey string) float64 {
	count := 0.0
	for _, task := range store.remediations {
		if !matchTaskOverviewFilter(task, filter) {
			continue
		}
		finding := &findingRecord{ClusterID: task.ClusterID, WorkspaceID: task.WorkspaceID, ProjectID: task.ProjectID, BaselineID: task.BaselineID}
		if overviewGroupKey(groupBy, finding) == groupKey {
			count++
		}
	}
	return count
}

func matchFindingOverviewFilter(finding *findingRecord, filter OverviewFilter) bool {
	if finding == nil {
		return false
	}
	if filter.WorkspaceID != 0 && finding.WorkspaceID != filter.WorkspaceID {
		return false
	}
	if filter.ProjectID != 0 && finding.ProjectID != filter.ProjectID {
		return false
	}
	if filter.TimeFrom != nil && finding.UpdatedAt.Before(*filter.TimeFrom) {
		return false
	}
	if filter.TimeTo != nil && finding.UpdatedAt.After(*filter.TimeTo) {
		return false
	}
	return true
}

func matchTaskOverviewFilter(task *RemediationTask, filter OverviewFilter) bool {
	if task == nil {
		return false
	}
	if filter.WorkspaceID != 0 && task.WorkspaceID != filter.WorkspaceID {
		return false
	}
	if filter.ProjectID != 0 && task.ProjectID != filter.ProjectID {
		return false
	}
	if filter.TimeFrom != nil && task.CreatedAt.Before(*filter.TimeFrom) {
		return false
	}
	if filter.TimeTo != nil && task.CreatedAt.After(*filter.TimeTo) {
		return false
	}
	return true
}

func maxFloat(left, right float64) float64 {
	if left > right {
		return left
	}
	return right
}
