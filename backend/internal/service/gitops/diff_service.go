package gitops

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	deliverydiff "kbmanage/backend/internal/integration/delivery/diff"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

// DiffService provides a minimal desired/live diff chain with cache fallback.
type DiffService struct {
	units      *repository.DeliveryUnitRepository
	cache      *DiffCache
	comparator deliverydiff.Comparator
}

func NewDiffService(
	units *repository.DeliveryUnitRepository,
	cache *DiffCache,
	comparator deliverydiff.Comparator,
) *DiffService {
	if comparator == nil {
		comparator = deliverydiff.NewNoopComparator()
	}
	return &DiffService{
		units:      units,
		cache:      cache,
		comparator: comparator,
	}
}

func (s *DiffService) GetOrBuild(ctx context.Context, unitID uint64, stageID uint64) (string, error) {
	if unitID == 0 {
		return "", errors.New("delivery unit id is required")
	}
	if s == nil || s.units == nil || s.comparator == nil {
		return "", ErrGitOpsNotConfigured
	}

	if s.cache != nil {
		cached, err := s.cache.GetDeliveryUnitDiff(ctx, unitID, stageID)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(cached) != "" {
			return cached, nil
		}
	}

	payload, err := s.buildDiffPayload(ctx, unitID, stageID)
	if err != nil {
		return "", err
	}
	if s.cache != nil {
		_ = s.cache.SetDeliveryUnitDiff(ctx, unitID, stageID, payload)
	}
	return payload, nil
}

func (s *DiffService) buildDiffPayload(ctx context.Context, unitID uint64, stageID uint64) (string, error) {
	unit, err := s.units.GetByID(ctx, unitID)
	if err != nil {
		return "", err
	}

	stageName := ""
	stageStatus := ""
	if stageID > 0 {
		stage, stageErr := s.units.GetEnvironmentStageByID(ctx, stageID)
		switch {
		case stageErr == nil && stage != nil:
			if stage.DeliveryUnitID == unitID {
				stageName = strings.TrimSpace(stage.Name)
				stageStatus = strings.TrimSpace(string(stage.Status))
			}
		case stageErr != nil && !errors.Is(stageErr, gorm.ErrRecordNotFound):
			return "", stageErr
		}
	}

	desired := map[string]any{
		"revision":      strings.TrimSpace(unit.DesiredRevision),
		"appVersion":    strings.TrimSpace(unit.DesiredAppVersion),
		"configVersion": strings.TrimSpace(unit.DesiredConfigVersion),
		"paused":        unit.Paused,
		"syncMode":      strings.TrimSpace(string(unit.SyncMode)),
	}
	live := map[string]any{
		"deliveryStatus": strings.TrimSpace(string(unit.DeliveryStatus)),
		"lastSyncedAt":   unit.LastSyncedAt,
		"stage": map[string]any{
			"id":     stageID,
			"name":   stageName,
			"status": stageStatus,
		},
	}
	desiredRaw, err := json.Marshal(desired)
	if err != nil {
		return "", err
	}
	liveRaw, err := json.Marshal(live)
	if err != nil {
		return "", err
	}

	summary, err := s.comparator.Compare(ctx, desiredRaw, liveRaw)
	if err != nil {
		return "", err
	}
	driftStatus := string(domain.DeliveryUnitStatusReady)
	if summary.HasChanges {
		driftStatus = string(domain.DeliveryUnitStatusOutOfSync)
	}

	payload := map[string]any{
		"summary":     summary,
		"driftStatus": driftStatus,
		"items":       []any{},
		"generatedAt": time.Now().UTC(),
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
