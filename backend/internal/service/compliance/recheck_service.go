package compliance

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
)

type RecheckTask struct {
	ID                    string         `json:"id"`
	FindingID             string         `json:"findingId"`
	TriggerSource         string         `json:"triggerSource"`
	Status                string         `json:"status"`
	ResultScanExecutionID string         `json:"resultScanExecutionId,omitempty"`
	Summary               string         `json:"summary,omitempty"`
	RequestedBy           uint64         `json:"requestedBy"`
	CreatedAt             time.Time      `json:"createdAt"`
	StartedAt             *time.Time     `json:"startedAt,omitempty"`
	CompletedAt           *time.Time     `json:"completedAt,omitempty"`
	ScopeSnapshot         *ScopeSnapshot `json:"scopeSnapshot,omitempty"`
	ClusterID             uint64         `json:"clusterId,omitempty"`
	WorkspaceID           uint64         `json:"workspaceId,omitempty"`
	ProjectID             uint64         `json:"projectId,omitempty"`
	BaselineID            string         `json:"baselineId,omitempty"`
}

type CreateRecheckInput struct {
	TriggerSource string
	Summary       string
	ScopeSnapshot *ScopeSnapshot
	ClusterID     uint64
	WorkspaceID   uint64
	ProjectID     uint64
	BaselineID    string
}

type CompleteRecheckInput struct {
	Passed                bool
	ResultScanExecutionID string
	Summary               string
}

type RecheckFilter struct {
	WorkspaceID   uint64
	ProjectID     uint64
	Status        string
	TriggerSource string
	TimeFrom      *time.Time
	TimeTo        *time.Time
}

type RecheckService struct {
	store *complianceStore
	audit *auditSvc.EventWriter
	now   func() time.Time
}

func NewRecheckService(auditWriter ...*auditSvc.EventWriter) *RecheckService {
	var writer *auditSvc.EventWriter
	if len(auditWriter) > 0 {
		writer = auditWriter[0]
	}
	return &RecheckService{store: defaultComplianceStore, audit: writer, now: time.Now}
}

func (s *RecheckService) CreateTask(ctx context.Context, operatorID uint64, findingID string, input CreateRecheckInput) (*RecheckTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	findingID = strings.TrimSpace(findingID)
	if findingID == "" {
		return nil, ErrFindingIDRequired
	}
	triggerSource, err := normalizeRecheckTriggerSource(input.TriggerSource, true)
	if err != nil {
		return nil, err
	}
	if triggerSource == "" {
		triggerSource = "manual"
	}
	now := s.now()

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	finding := s.store.ensureFindingLocked(findingID, input.ClusterID, input.WorkspaceID, input.ProjectID, input.BaselineID, input.ScopeSnapshot)
	item := &RecheckTask{
		ID:            uuid.NewString(),
		FindingID:     findingID,
		TriggerSource: triggerSource,
		Status:        "pending",
		Summary:       strings.TrimSpace(input.Summary),
		RequestedBy:   operatorID,
		CreatedAt:     now,
		ScopeSnapshot: cloneScopeSnapshotPtr(input.ScopeSnapshot),
		ClusterID:     nonZero(input.ClusterID, finding.ClusterID),
		WorkspaceID:   nonZero(input.WorkspaceID, finding.WorkspaceID),
		ProjectID:     nonZero(input.ProjectID, finding.ProjectID),
		BaselineID:    firstNonEmpty(input.BaselineID, finding.BaselineID),
	}
	s.store.rechecks[item.ID] = cloneRecheck(item)
	finding.RemediationStatus = "ready_for_recheck"
	finding.UpdatedAt = now
	s.writeAudit(ctx, operatorID, auditSvc.ComplianceAuditActionRecheckCreate, item.ID, domain.AuditOutcomeSuccess, map[string]any{
		"findingId":     item.FindingID,
		"status":        item.Status,
		"triggerSource": item.TriggerSource,
		"workspaceId":   item.WorkspaceID,
		"projectId":     item.ProjectID,
		"clusterId":     item.ClusterID,
		"baselineId":    item.BaselineID,
		"scopeSnapshot": item.ScopeSnapshot,
	})
	return cloneRecheck(item), nil
}

func (s *RecheckService) CompleteTask(ctx context.Context, operatorID uint64, recheckID string, input CompleteRecheckInput) (*RecheckTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	recheckID = strings.TrimSpace(recheckID)
	if recheckID == "" {
		return nil, ErrRecheckIDRequired
	}
	now := s.now()

	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	item, ok := s.store.rechecks[recheckID]
	if !ok {
		return nil, errors.New("recheck task not found")
	}
	updated := cloneRecheck(item)
	if updated.Status == "passed" || updated.Status == "failed" || updated.Status == "canceled" {
		return nil, fmtRecheckTerminal(updated.Status)
	}
	if updated.StartedAt == nil {
		updated.StartedAt = ptrTime(now)
	}
	updated.CompletedAt = ptrTime(now)
	updated.ResultScanExecutionID = strings.TrimSpace(input.ResultScanExecutionID)
	updated.Summary = strings.TrimSpace(input.Summary)
	if input.Passed {
		updated.Status = "passed"
	} else {
		updated.Status = "failed"
	}
	s.store.rechecks[recheckID] = cloneRecheck(updated)
	s.store.refreshFindingStatusLocked(updated.FindingID, now)
	s.writeAudit(ctx, operatorID, auditSvc.ComplianceAuditActionRecheckComplete, updated.ID, domain.AuditOutcomeSuccess, map[string]any{
		"findingId":             updated.FindingID,
		"status":                updated.Status,
		"resultScanExecutionId": updated.ResultScanExecutionID,
		"summary":               updated.Summary,
		"workspaceId":           updated.WorkspaceID,
		"projectId":             updated.ProjectID,
		"clusterId":             updated.ClusterID,
		"baselineId":            updated.BaselineID,
	})
	return cloneRecheck(updated), nil
}

func (s *RecheckService) RunPending(ctx context.Context, limit int) (int, error) {
	if s == nil || s.store == nil {
		return 0, ErrComplianceNotConfigured
	}
	if limit <= 0 {
		limit = 20
	}
	processed := 0
	now := s.now()

	s.store.mu.Lock()
	pendingIDs := make([]string, 0, limit)
	for id, item := range s.store.rechecks {
		if item.Status == "pending" {
			pendingIDs = append(pendingIDs, id)
			if len(pendingIDs) >= limit {
				break
			}
		}
	}
	for _, id := range pendingIDs {
		item := cloneRecheck(s.store.rechecks[id])
		item.Status = "running"
		item.StartedAt = ptrTime(now)
		s.store.rechecks[id] = cloneRecheck(item)
	}
	s.store.mu.Unlock()

	for _, id := range pendingIDs {
		_, _ = s.CompleteTask(ctx, 0, id, CompleteRecheckInput{Passed: true, Summary: "recheck worker completed"})
		processed++
	}
	return processed, nil
}

func (s *RecheckService) GetTask(_ context.Context, recheckID string) (*RecheckTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	recheckID = strings.TrimSpace(recheckID)
	if recheckID == "" {
		return nil, ErrRecheckIDRequired
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	item, ok := s.store.rechecks[recheckID]
	if !ok {
		return nil, errors.New("recheck task not found")
	}
	return cloneRecheck(item), nil
}

func (s *RecheckService) ListTasks(_ context.Context, filter RecheckFilter) ([]RecheckTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	items := make([]RecheckTask, 0, len(s.store.rechecks))
	for _, item := range s.store.rechecks {
		if filter.WorkspaceID != 0 && item.WorkspaceID != filter.WorkspaceID {
			continue
		}
		if filter.ProjectID != 0 && item.ProjectID != filter.ProjectID {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(item.Status, filter.Status) {
			continue
		}
		if filter.TriggerSource != "" && !strings.EqualFold(item.TriggerSource, filter.TriggerSource) {
			continue
		}
		if filter.TimeFrom != nil && item.CreatedAt.Before(*filter.TimeFrom) {
			continue
		}
		if filter.TimeTo != nil && item.CreatedAt.After(*filter.TimeTo) {
			continue
		}
		items = append(items, *cloneRecheck(item))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (s *RecheckService) writeAudit(ctx context.Context, operatorID uint64, action, resourceID string, outcome domain.AuditOutcome, details map[string]any) {
	if s == nil || s.audit == nil {
		return
	}
	var actorID *uint64
	if operatorID != 0 {
		actorID = &operatorID
	}
	_ = s.audit.WriteComplianceEvent(ctx, "", actorID, action, resourceID, outcome, details)
}

func cloneRecheck(item *RecheckTask) *RecheckTask {
	if item == nil {
		return nil
	}
	copyItem := *item
	copyItem.ScopeSnapshot = cloneScopeSnapshotPtr(item.ScopeSnapshot)
	copyItem.StartedAt = cloneTimePtr(item.StartedAt)
	copyItem.CompletedAt = cloneTimePtr(item.CompletedAt)
	return &copyItem
}

func normalizeRecheckTriggerSource(raw string, allowEmpty bool) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" && allowEmpty {
		return "", nil
	}
	switch value {
	case "manual", "remediation_done", "exception_expired":
		return value, nil
	default:
		return "", errors.New("triggerSource must be one of manual, remediation_done, exception_expired")
	}
}

func fmtRecheckTerminal(status string) error {
	return errors.New("recheck task is already " + strings.TrimSpace(status))
}
