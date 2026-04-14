package gitops

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

type ReleaseRecordInput struct {
	DeliveryUnitID uint64
	OperatorID     uint64
	SourceRevision string
	AppVersion     string
	ConfigVersion  string
	Environment    string
	Notes          string
}

type RevisionService struct {
	revisions *repository.ReleaseRevisionRepository
	units     *repository.DeliveryUnitRepository
}

func NewRevisionService(
	revisions *repository.ReleaseRevisionRepository,
	units *repository.DeliveryUnitRepository,
) *RevisionService {
	return &RevisionService{revisions: revisions, units: units}
}

func (s *RevisionService) List(ctx context.Context, unitID uint64, environment string) ([]domain.ReleaseRevision, error) {
	if s == nil || s.revisions == nil {
		return nil, ErrGitOpsNotConfigured
	}
	items, err := s.revisions.ListByDeliveryUnit(ctx, unitID)
	if err != nil {
		return nil, err
	}
	env := strings.ToLower(strings.TrimSpace(environment))
	if env == "" {
		return items, nil
	}
	filtered := make([]domain.ReleaseRevision, 0, len(items))
	for i := range items {
		scope := strings.ToLower(strings.TrimSpace(items[i].EffectiveScopeJSON))
		if scope == "" || strings.Contains(scope, env) {
			filtered = append(filtered, items[i])
		}
	}
	return filtered, nil
}

func (s *RevisionService) RecordRelease(ctx context.Context, in ReleaseRecordInput) (*domain.ReleaseRevision, error) {
	if s == nil || s.revisions == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if in.DeliveryUnitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if in.OperatorID == 0 {
		return nil, errors.New("operator id is required")
	}
	if err := s.revisions.MarkOthersHistorical(ctx, in.DeliveryUnitID, 0); err != nil {
		return nil, err
	}
	revision := &domain.ReleaseRevision{
		DeliveryUnitID:      in.DeliveryUnitID,
		SourceRevision:      strings.TrimSpace(in.SourceRevision),
		AppVersion:          strings.TrimSpace(in.AppVersion),
		ConfigVersion:       strings.TrimSpace(in.ConfigVersion),
		EffectiveScopeJSON:  normalizeRevisionScope(in.Environment),
		ReleaseNotesSummary: strings.TrimSpace(in.Notes),
		CreatedBy:           in.OperatorID,
		RollbackAvailable:   true,
		Status:              domain.ReleaseRevisionStatusActive,
	}
	if revision.SourceRevision == "" {
		revision.SourceRevision = "main"
	}
	if err := s.revisions.Create(ctx, revision); err != nil {
		return nil, err
	}
	return revision, nil
}

func (s *RevisionService) ResolveRollbackTarget(ctx context.Context, unitID uint64, targetReleaseID *uint64) (*domain.ReleaseRevision, error) {
	if s == nil || s.revisions == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if unitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if targetReleaseID != nil && *targetReleaseID > 0 {
		target, err := s.revisions.GetByID(ctx, *targetReleaseID)
		if err != nil {
			return nil, err
		}
		if target.DeliveryUnitID != unitID {
			return nil, errors.New("target release does not belong to delivery unit")
		}
		return target, nil
	}
	items, err := s.revisions.ListByDeliveryUnit(ctx, unitID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].Status == domain.ReleaseRevisionStatusActive {
			continue
		}
		if items[i].RollbackAvailable {
			return &items[i], nil
		}
	}
	if len(items) > 0 {
		return &items[0], nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *RevisionService) RollbackToRevision(
	ctx context.Context,
	unitID uint64,
	targetReleaseID *uint64,
	operatorID uint64,
) (*domain.ReleaseRevision, error) {
	if s == nil || s.revisions == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if unitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	target, err := s.ResolveRollbackTarget(ctx, unitID, targetReleaseID)
	if err != nil {
		return nil, err
	}

	items, err := s.revisions.ListByDeliveryUnit(ctx, unitID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		status := domain.ReleaseRevisionStatusHistorical
		switch {
		case items[i].ID == target.ID:
			status = domain.ReleaseRevisionStatusActive
		case items[i].Status == domain.ReleaseRevisionStatusActive:
			status = domain.ReleaseRevisionStatusRolledBack
		}
		if err := s.revisions.UpdateStatus(ctx, items[i].ID, status, true); err != nil {
			return nil, err
		}
	}
	if err := s.units.UpdateFields(ctx, unitID, map[string]any{
		"desired_revision":       target.SourceRevision,
		"desired_app_version":    target.AppVersion,
		"desired_config_version": target.ConfigVersion,
		"last_release_id":        target.ID,
	}); err != nil {
		return nil, err
	}
	_ = operatorID
	return s.revisions.GetByID(ctx, target.ID)
}

func normalizeRevisionScope(environment string) string {
	env := strings.TrimSpace(environment)
	if env == "" {
		return "all"
	}
	return fmt.Sprintf("env:%s", env)
}
