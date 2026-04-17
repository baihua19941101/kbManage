package compliance

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
)

type ArchiveExportTask struct {
	ID            string         `json:"id"`
	WorkspaceID   uint64         `json:"workspaceId,omitempty"`
	ProjectID     uint64         `json:"projectId,omitempty"`
	BaselineID    string         `json:"baselineId,omitempty"`
	ExportScope   string         `json:"exportScope"`
	Status        string         `json:"status"`
	ArtifactRef   string         `json:"artifactRef,omitempty"`
	RequestedBy   uint64         `json:"requestedBy"`
	StartedAt     *time.Time     `json:"startedAt,omitempty"`
	CompletedAt   *time.Time     `json:"completedAt,omitempty"`
	FailureReason string         `json:"failureReason,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	Filters       map[string]any `json:"filters,omitempty"`
	ArtifactData  []byte         `json:"-"`
}

type CreateArchiveExportInput struct {
	WorkspaceID uint64
	ProjectID   uint64
	BaselineID  string
	ExportScope string
	TimeFrom    *time.Time
	TimeTo      *time.Time
	Filters     map[string]any
}

type ArchiveExportFilter struct {
	WorkspaceID uint64
	ProjectID   uint64
	ExportScope string
	Status      string
	TimeFrom    *time.Time
	TimeTo      *time.Time
}

type ArchiveExportService struct {
	store *complianceStore
	audit *auditSvc.EventWriter
	now   func() time.Time
}

func NewArchiveExportService(auditWriter ...*auditSvc.EventWriter) *ArchiveExportService {
	var writer *auditSvc.EventWriter
	if len(auditWriter) > 0 {
		writer = auditWriter[0]
	}
	return &ArchiveExportService{store: defaultComplianceStore, audit: writer, now: time.Now}
}

func (s *ArchiveExportService) CreateExport(ctx context.Context, operatorID uint64, input CreateArchiveExportInput) (*ArchiveExportTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	exportScope, err := normalizeExportScope(input.ExportScope)
	if err != nil {
		return nil, err
	}
	now := s.now()
	filters := make(map[string]any, len(input.Filters)+2)
	for key, value := range input.Filters {
		filters[key] = value
	}
	if input.TimeFrom != nil {
		filters["timeFrom"] = input.TimeFrom.UTC().Format(time.RFC3339)
	}
	if input.TimeTo != nil {
		filters["timeTo"] = input.TimeTo.UTC().Format(time.RFC3339)
	}
	task := &ArchiveExportTask{
		ID:          uuid.NewString(),
		WorkspaceID: input.WorkspaceID,
		ProjectID:   input.ProjectID,
		BaselineID:  strings.TrimSpace(input.BaselineID),
		ExportScope: exportScope,
		Status:      "pending",
		RequestedBy: operatorID,
		CreatedAt:   now,
		Filters:     filters,
	}
	s.store.mu.Lock()
	s.store.exports[task.ID] = cloneArchiveExport(task)
	s.store.mu.Unlock()
	s.writeAudit(ctx, operatorID, auditSvc.ComplianceAuditActionArchiveExport, task.ID, domain.AuditOutcomeSuccess, map[string]any{
		"exportScope": exportScope,
		"workspaceId": task.WorkspaceID,
		"projectId":   task.ProjectID,
		"baselineId":  task.BaselineID,
	})
	return cloneArchiveExport(task), nil
}

func (s *ArchiveExportService) ListExports(_ context.Context, filter ArchiveExportFilter) ([]ArchiveExportTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	items := make([]ArchiveExportTask, 0, len(s.store.exports))
	for _, item := range s.store.exports {
		if filter.WorkspaceID != 0 && item.WorkspaceID != filter.WorkspaceID {
			continue
		}
		if filter.ProjectID != 0 && item.ProjectID != filter.ProjectID {
			continue
		}
		if filter.ExportScope != "" && !strings.EqualFold(item.ExportScope, filter.ExportScope) {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(item.Status, filter.Status) {
			continue
		}
		if filter.TimeFrom != nil && item.CreatedAt.Before(*filter.TimeFrom) {
			continue
		}
		if filter.TimeTo != nil && item.CreatedAt.After(*filter.TimeTo) {
			continue
		}
		items = append(items, *cloneArchiveExport(item))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (s *ArchiveExportService) GetExport(_ context.Context, exportID string) (*ArchiveExportTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	exportID = strings.TrimSpace(exportID)
	if exportID == "" {
		return nil, ErrExportIDRequired
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	item, ok := s.store.exports[exportID]
	if !ok {
		return nil, errors.New("archive export task not found")
	}
	return cloneArchiveExport(item), nil
}

func (s *ArchiveExportService) ProcessExport(ctx context.Context, exportID string) (*ArchiveExportTask, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	exportID = strings.TrimSpace(exportID)
	if exportID == "" {
		return nil, ErrExportIDRequired
	}
	now := s.now()
	s.store.mu.Lock()
	item, ok := s.store.exports[exportID]
	if !ok {
		s.store.mu.Unlock()
		return nil, errors.New("archive export task not found")
	}
	updated := cloneArchiveExport(item)
	updated.Status = "running"
	updated.StartedAt = ptrTime(now)
	s.store.exports[exportID] = cloneArchiveExport(updated)
	s.store.mu.Unlock()

	payload, err := s.buildPayload(updated)
	if err != nil {
		updated.Status = "failed"
		updated.FailureReason = err.Error()
		updated.CompletedAt = ptrTime(s.now())
		s.store.mu.Lock()
		s.store.exports[exportID] = cloneArchiveExport(updated)
		s.store.mu.Unlock()
		return nil, err
	}
	updated.Status = "succeeded"
	updated.ArtifactRef = "/api/v1/compliance/archive-exports/" + updated.ID + "/download"
	updated.ArtifactData = payload
	updated.CompletedAt = ptrTime(s.now())
	s.store.mu.Lock()
	s.store.exports[exportID] = cloneArchiveExport(updated)
	s.store.mu.Unlock()
	return cloneArchiveExport(updated), nil
}

func (s *ArchiveExportService) buildPayload(task *ArchiveExportTask) ([]byte, error) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	bundle := map[string]any{
		"task": task,
	}
	switch task.ExportScope {
	case "scans", "findings":
		bundle["findings"] = s.store.findings
	case "trends":
		bundle["trends"] = s.store.snapshots
	case "audit":
		bundle["audit"] = task.Filters
	case "bundle":
		bundle["findings"] = s.store.findings
		bundle["remediationTasks"] = s.store.remediations
		bundle["exceptions"] = s.store.exceptions
		bundle["rechecks"] = s.store.rechecks
		bundle["trends"] = s.store.snapshots
	}
	return json.Marshal(bundle)
}

func cloneArchiveExport(item *ArchiveExportTask) *ArchiveExportTask {
	if item == nil {
		return nil
	}
	copyItem := *item
	copyItem.StartedAt = cloneTimePtr(item.StartedAt)
	copyItem.CompletedAt = cloneTimePtr(item.CompletedAt)
	copyItem.ArtifactData = append([]byte(nil), item.ArtifactData...)
	if item.Filters != nil {
		copyItem.Filters = make(map[string]any, len(item.Filters))
		for key, value := range item.Filters {
			copyItem.Filters[key] = value
		}
	}
	return &copyItem
}

func normalizeExportScope(raw string) (string, error) {
	switch value := strings.ToLower(strings.TrimSpace(raw)); value {
	case "scans", "findings", "trends", "audit", "bundle":
		return value, nil
	default:
		return "", errors.New("exportScope must be one of scans, findings, trends, audit, bundle")
	}
}

func (s *ArchiveExportService) writeAudit(ctx context.Context, operatorID uint64, action, resourceID string, outcome domain.AuditOutcome, details map[string]any) {
	if s == nil || s.audit == nil {
		return
	}
	var actorID *uint64
	if operatorID != 0 {
		actorID = &operatorID
	}
	_ = s.audit.WriteComplianceEvent(ctx, "", actorID, action, resourceID, outcome, details)
}
