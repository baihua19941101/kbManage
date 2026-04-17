package compliance

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
)

var (
	ErrComplianceNotConfigured = errors.New("compliance service is not configured")
	ErrFindingIDRequired       = errors.New("finding id is required")
	ErrTaskIDRequired          = errors.New("task id is required")
	ErrExceptionIDRequired     = errors.New("exception id is required")
	ErrRecheckIDRequired       = errors.New("recheck id is required")
	ErrExportIDRequired        = errors.New("export id is required")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

type ScopeSnapshot struct {
	ClusterIDs   []uint64 `json:"clusterIds,omitempty"`
	WorkspaceIDs []uint64 `json:"workspaceIds,omitempty"`
	ProjectIDs   []uint64 `json:"projectIds,omitempty"`
	ScopeType    string   `json:"scopeType,omitempty"`
	ScopeRef     string   `json:"scopeRef,omitempty"`
	Namespace    string   `json:"namespace,omitempty"`
	NodeName     string   `json:"nodeName,omitempty"`
	ResourceKind string   `json:"resourceKind,omitempty"`
	ResourceName string   `json:"resourceName,omitempty"`
}

type BaselineSnapshot struct {
	BaselineID   string `json:"baselineId,omitempty"`
	Name         string `json:"name,omitempty"`
	StandardType string `json:"standardType,omitempty"`
	Version      string `json:"version,omitempty"`
}

type findingRecord struct {
	ID                string
	BaselineID        string
	ClusterID         uint64
	WorkspaceID       uint64
	ProjectID         uint64
	RiskLevel         string
	RemediationStatus string
	ScopeSnapshot     ScopeSnapshot
	UpdatedAt         time.Time
}

type RemediationTask struct {
	ID                string         `json:"id"`
	FindingID         string         `json:"findingId"`
	Title             string         `json:"title"`
	Owner             string         `json:"owner"`
	Priority          string         `json:"priority"`
	Status            string         `json:"status"`
	DueAt             *time.Time     `json:"dueAt,omitempty"`
	ResolutionSummary string         `json:"resolutionSummary,omitempty"`
	Overdue           bool           `json:"overdue"`
	CreatedBy         uint64         `json:"createdBy"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	CompletedAt       *time.Time     `json:"completedAt,omitempty"`
	ScopeSnapshot     *ScopeSnapshot `json:"scopeSnapshot,omitempty"`
	ClusterID         uint64         `json:"clusterId,omitempty"`
	WorkspaceID       uint64         `json:"workspaceId,omitempty"`
	ProjectID         uint64         `json:"projectId,omitempty"`
	BaselineID        string         `json:"baselineId,omitempty"`
}

type CreateRemediationTaskInput struct {
	Title         string
	Owner         string
	Priority      string
	DueAt         *time.Time
	Summary       string
	ScopeSnapshot *ScopeSnapshot
	ClusterID     uint64
	WorkspaceID   uint64
	ProjectID     uint64
	BaselineID    string
}

type UpdateRemediationTaskInput struct {
	Status            string
	ResolutionSummary string
}

type RemediationTaskFilter struct {
	WorkspaceID uint64
	ProjectID   uint64
	Owner       string
	Status      string
	Priority    string
	TimeFrom    *time.Time
	TimeTo      *time.Time
}

type complianceStore struct {
	mu           sync.RWMutex
	findings     map[string]*findingRecord
	remediations map[string]*RemediationTask
	exceptions   map[string]*ComplianceExceptionRequest
	rechecks     map[string]*RecheckTask
	snapshots    []*ComplianceTrendPoint
	exports      map[string]*ArchiveExportTask
}

var defaultComplianceStore = &complianceStore{
	findings:     make(map[string]*findingRecord),
	remediations: make(map[string]*RemediationTask),
	exceptions:   make(map[string]*ComplianceExceptionRequest),
	rechecks:     make(map[string]*RecheckTask),
	snapshots:    make([]*ComplianceTrendPoint, 0),
	exports:      make(map[string]*ArchiveExportTask),
}

type RemediationService struct {
	store *complianceStore
	audit *auditSvc.EventWriter
	now   func() time.Time
}

func NewRemediationService(auditWriter ...*auditSvc.EventWriter) *RemediationService {
	var writer *auditSvc.EventWriter
	if len(auditWriter) > 0 {
		writer = auditWriter[0]
	}
	return &RemediationService{store: defaultComplianceStore, audit: writer, now: time.Now}
}

func (s *RemediationService) CreateTask(ctx context.Context, operatorID uint64, findingID string, input CreateRemediationTaskInput) (*RemediationTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	findingID = strings.TrimSpace(findingID)
	if findingID == "" {
		return nil, ErrFindingIDRequired
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, errors.New("title is required")
	}
	owner := strings.TrimSpace(input.Owner)
	if owner == "" {
		return nil, errors.New("owner is required")
	}
	priority, err := normalizePriority(input.Priority)
	if err != nil {
		return nil, err
	}
	now := s.now()

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	finding := s.store.ensureFindingLocked(findingID, input.ClusterID, input.WorkspaceID, input.ProjectID, input.BaselineID, input.ScopeSnapshot)
	task := &RemediationTask{
		ID:            uuid.NewString(),
		FindingID:     findingID,
		Title:         title,
		Owner:         owner,
		Priority:      priority,
		Status:        "todo",
		DueAt:         cloneTimePtr(input.DueAt),
		CreatedBy:     operatorID,
		CreatedAt:     now,
		UpdatedAt:     now,
		ScopeSnapshot: cloneScopeSnapshotPtr(input.ScopeSnapshot),
		ClusterID:     nonZero(input.ClusterID, finding.ClusterID),
		WorkspaceID:   nonZero(input.WorkspaceID, finding.WorkspaceID),
		ProjectID:     nonZero(input.ProjectID, finding.ProjectID),
		BaselineID:    firstNonEmpty(strings.TrimSpace(input.BaselineID), finding.BaselineID),
	}
	if summary := strings.TrimSpace(input.Summary); summary != "" {
		task.ResolutionSummary = summary
	}
	task.Overdue = isTaskOverdue(task, now)
	s.store.remediations[task.ID] = cloneRemediationTask(task)
	finding.RemediationStatus = "in_progress"
	finding.UpdatedAt = now

	s.writeAudit(ctx, operatorID, auditSvc.ComplianceAuditActionRemediationCreate, task.ID, domain.AuditOutcomeSuccess, map[string]any{
		"findingId":     findingID,
		"owner":         task.Owner,
		"priority":      task.Priority,
		"status":        task.Status,
		"workspaceId":   task.WorkspaceID,
		"projectId":     task.ProjectID,
		"clusterId":     task.ClusterID,
		"baselineId":    task.BaselineID,
		"scopeSnapshot": task.ScopeSnapshot,
	})
	return cloneRemediationTask(task), nil
}

func (s *RemediationService) UpdateTask(ctx context.Context, operatorID uint64, taskID string, input UpdateRemediationTaskInput) (*RemediationTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, ErrTaskIDRequired
	}
	status, err := normalizeRemediationStatus(input.Status, true)
	if err != nil {
		return nil, err
	}
	now := s.now()

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	task, ok := s.store.remediations[taskID]
	if !ok {
		return nil, errors.New("remediation task not found")
	}
	updated := cloneRemediationTask(task)
	if status != "" {
		if !canTransitRemediationStatus(updated.Status, status) {
			return nil, fmt.Errorf("%w: %s -> %s", ErrInvalidStatusTransition, updated.Status, status)
		}
		updated.Status = status
	}
	if summary := strings.TrimSpace(input.ResolutionSummary); summary != "" {
		updated.ResolutionSummary = summary
	}
	if updated.Status == "done" {
		if strings.TrimSpace(updated.ResolutionSummary) == "" {
			return nil, errors.New("resolution summary is required when status is done")
		}
		updated.CompletedAt = ptrTime(now)
	} else {
		updated.CompletedAt = nil
	}
	updated.UpdatedAt = now
	updated.Overdue = isTaskOverdue(updated, now)
	s.store.remediations[taskID] = cloneRemediationTask(updated)
	s.store.refreshFindingStatusLocked(updated.FindingID, now)

	s.writeAudit(ctx, operatorID, auditSvc.ComplianceAuditActionRemediationUpdate, updated.ID, domain.AuditOutcomeSuccess, map[string]any{
		"findingId":         updated.FindingID,
		"status":            updated.Status,
		"resolutionSummary": updated.ResolutionSummary,
		"workspaceId":       updated.WorkspaceID,
		"projectId":         updated.ProjectID,
		"clusterId":         updated.ClusterID,
		"baselineId":        updated.BaselineID,
		"completedAt":       updated.CompletedAt,
		"scopeSnapshot":     updated.ScopeSnapshot,
	})
	return cloneRemediationTask(updated), nil
}

func (s *RemediationService) ListTasks(_ context.Context, filter RemediationTaskFilter) ([]RemediationTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	now := s.now()
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()

	items := make([]RemediationTask, 0, len(s.store.remediations))
	for _, task := range s.store.remediations {
		if filter.WorkspaceID != 0 && task.WorkspaceID != filter.WorkspaceID {
			continue
		}
		if filter.ProjectID != 0 && task.ProjectID != filter.ProjectID {
			continue
		}
		if filter.Owner != "" && !strings.EqualFold(strings.TrimSpace(task.Owner), strings.TrimSpace(filter.Owner)) {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(task.Status, filter.Status) {
			continue
		}
		if filter.Priority != "" && !strings.EqualFold(task.Priority, filter.Priority) {
			continue
		}
		if filter.TimeFrom != nil && task.CreatedAt.Before(*filter.TimeFrom) {
			continue
		}
		if filter.TimeTo != nil && task.CreatedAt.After(*filter.TimeTo) {
			continue
		}
		copyItem := cloneRemediationTask(task)
		copyItem.Overdue = isTaskOverdue(copyItem, now)
		items = append(items, *copyItem)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

func (s *RemediationService) GetTask(_ context.Context, taskID string) (*RemediationTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, ErrTaskIDRequired
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	task, ok := s.store.remediations[taskID]
	if !ok {
		return nil, errors.New("remediation task not found")
	}
	copyItem := cloneRemediationTask(task)
	copyItem.Overdue = isTaskOverdue(copyItem, s.now())
	return copyItem, nil
}

func (s *RemediationService) writeAudit(ctx context.Context, operatorID uint64, action, resourceID string, outcome domain.AuditOutcome, details map[string]any) {
	if s == nil || s.audit == nil {
		return
	}
	var actorID *uint64
	if operatorID != 0 {
		actorID = &operatorID
	}
	_ = s.audit.WriteComplianceEvent(ctx, "", actorID, action, resourceID, outcome, details)
}

func (store *complianceStore) ensureFindingLocked(findingID string, clusterID, workspaceID, projectID uint64, baselineID string, scope *ScopeSnapshot) *findingRecord {
	item, ok := store.findings[findingID]
	if !ok {
		item = &findingRecord{ID: findingID, RemediationStatus: "open"}
		store.findings[findingID] = item
	}
	if clusterID != 0 {
		item.ClusterID = clusterID
	}
	if workspaceID != 0 {
		item.WorkspaceID = workspaceID
	}
	if projectID != 0 {
		item.ProjectID = projectID
	}
	if baselineID = strings.TrimSpace(baselineID); baselineID != "" {
		item.BaselineID = baselineID
	}
	if scope != nil {
		item.ScopeSnapshot = *cloneScopeSnapshotPtr(scope)
	}
	item.UpdatedAt = time.Now()
	return item
}

func (store *complianceStore) refreshFindingStatusLocked(findingID string, now time.Time) {
	finding, ok := store.findings[findingID]
	if !ok {
		return
	}
	status := "open"
	for _, exception := range store.exceptions {
		if exception.FindingID != findingID {
			continue
		}
		if exception.Status == "active" || exception.Status == "approved" {
			status = "exception_active"
			break
		}
	}
	if status != "exception_active" {
		hasReadyForRecheck := false
		hasInProgress := false
		for _, recheck := range store.rechecks {
			if recheck.FindingID != findingID {
				continue
			}
			switch recheck.Status {
			case "passed":
				status = "closed"
			case "pending", "running":
				hasReadyForRecheck = true
			case "failed":
				hasInProgress = true
			}
		}
		if status != "closed" {
			for _, task := range store.remediations {
				if task.FindingID != findingID {
					continue
				}
				switch task.Status {
				case "done":
					hasReadyForRecheck = true
				case "todo", "in_progress", "blocked":
					hasInProgress = true
				}
			}
			if hasReadyForRecheck {
				status = "ready_for_recheck"
			} else if hasInProgress {
				status = "in_progress"
			}
		}
	}
	finding.RemediationStatus = status
	finding.UpdatedAt = now
}

func cloneRemediationTask(task *RemediationTask) *RemediationTask {
	if task == nil {
		return nil
	}
	copyItem := *task
	copyItem.DueAt = cloneTimePtr(task.DueAt)
	copyItem.CompletedAt = cloneTimePtr(task.CompletedAt)
	copyItem.ScopeSnapshot = cloneScopeSnapshotPtr(task.ScopeSnapshot)
	return &copyItem
}

func cloneScopeSnapshotPtr(item *ScopeSnapshot) *ScopeSnapshot {
	if item == nil {
		return nil
	}
	copyItem := *item
	copyItem.ClusterIDs = append([]uint64(nil), item.ClusterIDs...)
	copyItem.WorkspaceIDs = append([]uint64(nil), item.WorkspaceIDs...)
	copyItem.ProjectIDs = append([]uint64(nil), item.ProjectIDs...)
	return &copyItem
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func ptrTime(value time.Time) *time.Time { return &value }

func normalizePriority(raw string) (string, error) {
	switch value := strings.ToLower(strings.TrimSpace(raw)); value {
	case "low", "medium", "high", "critical":
		return value, nil
	default:
		return "", errors.New("priority must be one of low, medium, high, critical")
	}
}

func normalizeRemediationStatus(raw string, allowEmpty bool) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" && allowEmpty {
		return "", nil
	}
	switch value {
	case "todo", "in_progress", "blocked", "done", "canceled":
		return value, nil
	default:
		return "", errors.New("status must be one of todo, in_progress, blocked, done, canceled")
	}
}

func canTransitRemediationStatus(from, to string) bool {
	if from == to || to == "" {
		return true
	}
	switch from {
	case "todo":
		return to == "in_progress" || to == "blocked" || to == "canceled"
	case "in_progress":
		return to == "blocked" || to == "done" || to == "canceled"
	case "blocked":
		return to == "in_progress" || to == "canceled"
	case "done", "canceled":
		return false
	default:
		return false
	}
}

func isTaskOverdue(task *RemediationTask, now time.Time) bool {
	return task != nil && task.DueAt != nil && task.CompletedAt == nil && task.Status != "canceled" && task.DueAt.Before(now)
}

func nonZero(left, right uint64) uint64 {
	if left != 0 {
		return left
	}
	return right
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
