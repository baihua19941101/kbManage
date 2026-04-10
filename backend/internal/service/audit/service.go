package audit

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

var (
	errInvalidTimeRange      = errors.New("startAt must be earlier than endAt")
	errOperatorIDRequired    = errors.New("operator id is required")
	errAuditExportIDRequired = errors.New("task id is required")
)

type QueryEventsRequest struct {
	StartAt *time.Time
	EndAt   *time.Time
	ActorID *uint64
	Action  string
	Outcome string
	Limit   int
}

type SubmitExportRequest struct {
	StartAt *time.Time
	EndAt   *time.Time
	ActorID *uint64
	Action  string
	Outcome string
}

type Service struct {
	auditRepo       *repository.AuditRepository
	auditExportRepo *repository.AuditExportRepository
}

func NewService(auditRepo *repository.AuditRepository, auditExportRepo *repository.AuditExportRepository) *Service {
	return &Service{
		auditRepo:       auditRepo,
		auditExportRepo: auditExportRepo,
	}
}

func (s *Service) QueryEvents(ctx context.Context, req QueryEventsRequest) ([]domain.AuditEvent, error) {
	if req.StartAt != nil && req.EndAt != nil && req.StartAt.After(*req.EndAt) {
		return nil, errInvalidTimeRange
	}
	if s.auditRepo == nil {
		return []domain.AuditEvent{}, nil
	}

	outcome, err := normalizeOutcome(req.Outcome)
	if err != nil {
		return nil, err
	}

	return s.auditRepo.Query(ctx, repository.AuditQuery{
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
		ActorID: req.ActorID,
		Action:  strings.TrimSpace(req.Action),
		Outcome: outcome,
		Limit:   req.Limit,
	})
}

func (s *Service) SubmitExport(ctx context.Context, operatorID uint64, req SubmitExportRequest) (*repository.AuditExportTask, error) {
	if operatorID == 0 {
		return nil, errOperatorIDRequired
	}
	if req.StartAt != nil && req.EndAt != nil && req.StartAt.After(*req.EndAt) {
		return nil, errInvalidTimeRange
	}
	if s.auditExportRepo == nil {
		return nil, errors.New("audit export repository is not initialized")
	}

	outcome, err := normalizeOutcome(req.Outcome)
	if err != nil {
		return nil, err
	}

	task, err := s.auditExportRepo.Create(ctx, operatorID, repository.AuditQuery{
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
		ActorID: req.ActorID,
		Action:  strings.TrimSpace(req.Action),
		Outcome: outcome,
		Limit:   10000, // export uses a higher upper bound than list query.
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditExportRepo.Enqueue(ctx, task.ID); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) GetExportTask(ctx context.Context, taskID string) (*repository.AuditExportTask, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, errAuditExportIDRequired
	}
	if s.auditExportRepo == nil {
		return nil, errors.New("audit export repository is not initialized")
	}
	return s.auditExportRepo.Get(ctx, taskID)
}

func (s *Service) ProcessExportTask(ctx context.Context, taskID string) error {
	if strings.TrimSpace(taskID) == "" {
		return errAuditExportIDRequired
	}
	if s.auditExportRepo == nil {
		return errors.New("audit export repository is not initialized")
	}

	task, err := s.auditExportRepo.Get(ctx, taskID)
	if err != nil {
		return err
	}
	if err := s.auditExportRepo.MarkRunning(ctx, taskID); err != nil {
		return err
	}

	result, err := s.QueryEvents(ctx, QueryEventsRequest{
		StartAt: task.Filters.StartAt,
		EndAt:   task.Filters.EndAt,
		ActorID: task.Filters.ActorID,
		Action:  task.Filters.Action,
		Outcome: task.Filters.Outcome,
		Limit:   task.Filters.Limit,
	})
	if err != nil {
		_ = s.auditExportRepo.MarkFailed(ctx, taskID, err.Error())
		return err
	}

	downloadURL := fmt.Sprintf("/api/v1/audits/exports/%s/download", taskID)
	return s.auditExportRepo.MarkSucceeded(ctx, taskID, len(result), downloadURL)
}

func normalizeOutcome(raw string) (string, error) {
	outcome := strings.TrimSpace(strings.ToLower(raw))
	if outcome == "" {
		return "", nil
	}

	switch domain.AuditOutcome(outcome) {
	case domain.AuditOutcomeSuccess, domain.AuditOutcomeDenied, domain.AuditOutcomeFailed:
		return outcome, nil
	default:
		return "", fmt.Errorf("invalid outcome: %s", raw)
	}
}
