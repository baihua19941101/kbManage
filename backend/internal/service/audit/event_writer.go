package audit

import (
	"context"
	"encoding/json"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type EventWriter struct {
	repo *repository.AuditRepository
}

func NewEventWriter(repo *repository.AuditRepository) *EventWriter {
	return &EventWriter{repo: repo}
}

func (w *EventWriter) Write(
	ctx context.Context,
	requestID string,
	actorID *uint64,
	action string,
	resourceType string,
	resourceID string,
	outcome domain.AuditOutcome,
	details map[string]any,
) error {
	payload, err := json.Marshal(details)
	if err != nil {
		return err
	}

	event := &domain.AuditEvent{
		RequestID:    requestID,
		ActorID:      actorID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Outcome:      outcome,
		Details:      payload,
	}
	return w.repo.Create(ctx, event)
}
