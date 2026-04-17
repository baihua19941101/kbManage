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

type ComplianceExceptionRequest struct {
	ID            string         `json:"id"`
	FindingID     string         `json:"findingId"`
	Status        string         `json:"status"`
	Reason        string         `json:"reason"`
	StartsAt      time.Time      `json:"startsAt"`
	ExpiresAt     time.Time      `json:"expiresAt"`
	ReviewComment string         `json:"reviewComment,omitempty"`
	RequestedBy   uint64         `json:"requestedBy"`
	ReviewedBy    uint64         `json:"reviewedBy,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	ActivatedAt   *time.Time     `json:"activatedAt,omitempty"`
	ScopeSnapshot *ScopeSnapshot `json:"scopeSnapshot,omitempty"`
	ClusterID     uint64         `json:"clusterId,omitempty"`
	WorkspaceID   uint64         `json:"workspaceId,omitempty"`
	ProjectID     uint64         `json:"projectId,omitempty"`
	BaselineID    string         `json:"baselineId,omitempty"`
}

type CreateExceptionInput struct {
	Reason        string
	StartsAt      time.Time
	ExpiresAt     time.Time
	ScopeSnapshot *ScopeSnapshot
	ClusterID     uint64
	WorkspaceID   uint64
	ProjectID     uint64
	BaselineID    string
}

type ReviewExceptionInput struct {
	Decision      string
	ReviewComment string
}

type ExceptionFilter struct {
	WorkspaceID uint64
	ProjectID   uint64
	Status      string
	BaselineID  string
}

type ExceptionService struct {
	store *complianceStore
	audit *auditSvc.EventWriter
	now   func() time.Time
}

func NewExceptionService(auditWriter ...*auditSvc.EventWriter) *ExceptionService {
	var writer *auditSvc.EventWriter
	if len(auditWriter) > 0 {
		writer = auditWriter[0]
	}
	return &ExceptionService{store: defaultComplianceStore, audit: writer, now: time.Now}
}

func (s *ExceptionService) CreateException(ctx context.Context, operatorID uint64, findingID string, input CreateExceptionInput) (*ComplianceExceptionRequest, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	findingID = strings.TrimSpace(findingID)
	if findingID == "" {
		return nil, ErrFindingIDRequired
	}
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, errors.New("reason is required")
	}
	if !input.ExpiresAt.After(input.StartsAt) {
		return nil, errors.New("expiresAt must be later than startsAt")
	}
	now := s.now()
	status := "pending"

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	finding := s.store.ensureFindingLocked(findingID, input.ClusterID, input.WorkspaceID, input.ProjectID, input.BaselineID, input.ScopeSnapshot)
	item := &ComplianceExceptionRequest{
		ID:            uuid.NewString(),
		FindingID:     findingID,
		Status:        status,
		Reason:        reason,
		StartsAt:      input.StartsAt,
		ExpiresAt:     input.ExpiresAt,
		RequestedBy:   operatorID,
		CreatedAt:     now,
		UpdatedAt:     now,
		ScopeSnapshot: cloneScopeSnapshotPtr(input.ScopeSnapshot),
		ClusterID:     nonZero(input.ClusterID, finding.ClusterID),
		WorkspaceID:   nonZero(input.WorkspaceID, finding.WorkspaceID),
		ProjectID:     nonZero(input.ProjectID, finding.ProjectID),
		BaselineID:    firstNonEmpty(input.BaselineID, finding.BaselineID),
	}
	s.store.exceptions[item.ID] = cloneException(item)
	s.store.refreshFindingStatusLocked(findingID, now)
	s.writeAudit(ctx, operatorID, auditSvc.ComplianceAuditActionExceptionRequest, item.ID, domain.AuditOutcomeSuccess, map[string]any{
		"findingId":     findingID,
		"status":        item.Status,
		"workspaceId":   item.WorkspaceID,
		"projectId":     item.ProjectID,
		"clusterId":     item.ClusterID,
		"baselineId":    item.BaselineID,
		"scopeSnapshot": item.ScopeSnapshot,
	})
	return cloneException(item), nil
}

func (s *ExceptionService) ReviewException(ctx context.Context, reviewerID uint64, exceptionID string, input ReviewExceptionInput) (*ComplianceExceptionRequest, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	exceptionID = strings.TrimSpace(exceptionID)
	if exceptionID == "" {
		return nil, ErrExceptionIDRequired
	}
	decision, err := normalizeExceptionDecision(input.Decision)
	if err != nil {
		return nil, err
	}
	now := s.now()

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	item, ok := s.store.exceptions[exceptionID]
	if !ok {
		return nil, errors.New("compliance exception not found")
	}
	updated := cloneException(item)
	updated.ReviewedBy = reviewerID
	updated.ReviewComment = strings.TrimSpace(input.ReviewComment)
	updated.UpdatedAt = now
	switch decision {
	case "approve":
		updated.Status = "approved"
		if !now.Before(updated.StartsAt) && now.Before(updated.ExpiresAt) {
			updated.Status = "active"
			updated.ActivatedAt = ptrTime(now)
		}
	case "reject":
		updated.Status = "rejected"
		updated.ActivatedAt = nil
	case "revoke":
		if updated.Status != "approved" && updated.Status != "active" {
			return nil, errors.New("only approved or active exceptions can be revoked")
		}
		updated.Status = "revoked"
		updated.ActivatedAt = nil
	}
	s.store.exceptions[exceptionID] = cloneException(updated)
	s.store.refreshFindingStatusLocked(updated.FindingID, now)
	s.writeAudit(ctx, reviewerID, auditSvc.ComplianceAuditActionExceptionReview, updated.ID, domain.AuditOutcomeSuccess, map[string]any{
		"findingId":     updated.FindingID,
		"decision":      decision,
		"status":        updated.Status,
		"reviewComment": updated.ReviewComment,
		"workspaceId":   updated.WorkspaceID,
		"projectId":     updated.ProjectID,
		"clusterId":     updated.ClusterID,
		"baselineId":    updated.BaselineID,
		"scopeSnapshot": updated.ScopeSnapshot,
	})
	return cloneException(updated), nil
}

func (s *ExceptionService) ListExceptions(_ context.Context, filter ExceptionFilter) ([]ComplianceExceptionRequest, error) {
	if s == nil || s.store == nil {
		return nil, ErrComplianceNotConfigured
	}
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	items := make([]ComplianceExceptionRequest, 0, len(s.store.exceptions))
	for _, item := range s.store.exceptions {
		if filter.WorkspaceID != 0 && item.WorkspaceID != filter.WorkspaceID {
			continue
		}
		if filter.ProjectID != 0 && item.ProjectID != filter.ProjectID {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(item.Status, filter.Status) {
			continue
		}
		if filter.BaselineID != "" && item.BaselineID != strings.TrimSpace(filter.BaselineID) {
			continue
		}
		items = append(items, *cloneException(item))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

func (s *ExceptionService) ExpireDueExceptions(ctx context.Context, now time.Time) (int, error) {
	if s == nil || s.store == nil {
		return 0, ErrComplianceNotConfigured
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	count := 0
	for id, item := range s.store.exceptions {
		updated := cloneException(item)
		changed := false
		if updated.Status == "approved" && !now.Before(updated.StartsAt) && now.Before(updated.ExpiresAt) {
			updated.Status = "active"
			updated.ActivatedAt = ptrTime(now)
			changed = true
		}
		if (updated.Status == "active" || updated.Status == "approved") && !now.Before(updated.ExpiresAt) {
			updated.Status = "expired"
			updated.UpdatedAt = now
			updated.ActivatedAt = nil
			changed = true
			count++
			s.writeAudit(ctx, 0, auditSvc.ComplianceAuditActionExceptionReview, updated.ID, domain.AuditOutcomeSuccess, map[string]any{
				"findingId":   updated.FindingID,
				"decision":    "expire",
				"status":      updated.Status,
				"workspaceId": updated.WorkspaceID,
				"projectId":   updated.ProjectID,
				"clusterId":   updated.ClusterID,
			})
		}
		if changed {
			updated.UpdatedAt = now
			s.store.exceptions[id] = cloneException(updated)
			s.store.refreshFindingStatusLocked(updated.FindingID, now)
		}
	}
	return count, nil
}

func (s *ExceptionService) writeAudit(ctx context.Context, operatorID uint64, action, resourceID string, outcome domain.AuditOutcome, details map[string]any) {
	if s == nil || s.audit == nil {
		return
	}
	var actorID *uint64
	if operatorID != 0 {
		actorID = &operatorID
	}
	_ = s.audit.WriteComplianceEvent(ctx, "", actorID, action, resourceID, outcome, details)
}

func cloneException(item *ComplianceExceptionRequest) *ComplianceExceptionRequest {
	if item == nil {
		return nil
	}
	copyItem := *item
	copyItem.ScopeSnapshot = cloneScopeSnapshotPtr(item.ScopeSnapshot)
	copyItem.ActivatedAt = cloneTimePtr(item.ActivatedAt)
	return &copyItem
}

func normalizeExceptionDecision(raw string) (string, error) {
	switch value := strings.ToLower(strings.TrimSpace(raw)); value {
	case "approve", "reject", "revoke":
		return value, nil
	default:
		return "", errors.New("decision must be one of approve, reject, revoke")
	}
}
