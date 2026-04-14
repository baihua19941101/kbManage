package gitops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	deliverydiff "kbmanage/backend/internal/integration/delivery/diff"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

const (
	PermissionGitOpsRead         = "gitops:read"
	PermissionGitOpsManageSource = "gitops:manage-source"
	PermissionGitOpsSync         = "gitops:sync"
	PermissionGitOpsPromote      = "gitops:promote"
	PermissionGitOpsRollback     = "gitops:rollback"
	PermissionGitOpsOverride     = "gitops:override"
)

var (
	ErrGitOpsNotConfigured = errors.New("gitops service is not configured")
	ErrGitOpsScopeDenied   = errors.New("gitops scope access denied")
)

type SourceListFilter struct {
	SourceType string
	Status     string
}

type CreateSourceInput struct {
	Name          string
	SourceType    string
	Endpoint      string
	DefaultRef    string
	CredentialRef string
	WorkspaceID   uint64
	ProjectID     uint64
}

type UpdateSourceInput struct {
	Name          *string
	DefaultRef    *string
	CredentialRef *string
	Disabled      *bool
}

type CreateTargetGroupInput struct {
	Name            string
	WorkspaceID     uint64
	ProjectID       uint64
	ClusterRefs     []uint64
	SelectorSummary string
	Description     string
}

type UpdateTargetGroupInput struct {
	Name            *string
	ClusterRefs     *[]uint64
	SelectorSummary *string
	Description     *string
	Disabled        *bool
}

type EnvironmentStageInput struct {
	Name          string
	OrderIndex    int
	TargetGroupID uint64
	PromotionMode string
	Paused        bool
}

type ConfigurationOverlayInput struct {
	OverlayType    string
	OverlayRef     string
	Precedence     int
	EffectiveScope string
}

type CreateDeliveryUnitInput struct {
	Name                 string
	WorkspaceID          uint64
	ProjectID            uint64
	SourceID             uint64
	SourcePath           string
	DefaultNamespace     string
	SyncMode             string
	DesiredRevision      string
	DesiredAppVersion    string
	DesiredConfigVersion string
	Environments         []EnvironmentStageInput
	Overlays             []ConfigurationOverlayInput
}

type UpdateDeliveryUnitInput struct {
	Name                 *string
	SourcePath           *string
	DefaultNamespace     *string
	SyncMode             *string
	DesiredRevision      *string
	DesiredAppVersion    *string
	DesiredConfigVersion *string
	Environments         *[]EnvironmentStageInput
	Overlays             *[]ConfigurationOverlayInput
}

type DeliveryUnitDetail struct {
	Unit         domain.ApplicationDeliveryUnit
	Environments []domain.EnvironmentStage
	Overlays     []domain.ConfigurationOverlay
}

type ActionRequest struct {
	RequestID          string
	DeliveryUnitID     uint64
	EnvironmentStageID *uint64
	ActionType         domain.DeliveryActionType
	TargetReleaseID    *uint64
	PayloadJSON        string
}

type Service struct {
	sources     *repository.DeliverySourceRepository
	targets     *repository.ClusterTargetGroupRepository
	units       *repository.DeliveryUnitRepository
	revisions   *repository.ReleaseRevisionRepository
	operations  *repository.DeliveryOperationRepository
	scope       *ScopeService
	progress    *ProgressCache
	diffs       *DiffCache
	diffService *DiffService
	locks       *LockService
	queue       OperationQueue
	defaultLock time.Duration
}

func NewService(
	sources *repository.DeliverySourceRepository,
	targets *repository.ClusterTargetGroupRepository,
	units *repository.DeliveryUnitRepository,
	revisions *repository.ReleaseRevisionRepository,
	operations *repository.DeliveryOperationRepository,
	scope *ScopeService,
	progress *ProgressCache,
	diffs *DiffCache,
	locks *LockService,
	queue OperationQueue,
) *Service {
	svc := &Service{
		sources:     sources,
		targets:     targets,
		units:       units,
		revisions:   revisions,
		operations:  operations,
		scope:       scope,
		progress:    progress,
		diffs:       diffs,
		locks:       locks,
		queue:       queue,
		defaultLock: 30 * time.Second,
	}
	svc.diffService = NewDiffService(units, diffs, deliverydiff.NewNoopComparator())
	return svc
}

func (s *Service) ListSources(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	projectID uint64,
	filter SourceListFilter,
) ([]domain.DeliverySource, error) {
	if s == nil || s.sources == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if err := s.validateScope(ctx, userID, workspaceID, projectID, PermissionGitOpsRead); err != nil {
		return nil, err
	}
	sourceType, err := normalizeSourceType(filter.SourceType, true)
	if err != nil {
		return nil, err
	}
	status, err := normalizeSourceStatus(filter.Status, true)
	if err != nil {
		return nil, err
	}
	return s.sources.ListByScope(ctx, uint64PtrOrNil(workspaceID), uint64PtrOrNil(projectID), sourceType, status)
}

func (s *Service) CreateSource(ctx context.Context, userID uint64, input CreateSourceInput) (*domain.DeliverySource, error) {
	if s == nil || s.sources == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if input.WorkspaceID == 0 {
		return nil, errors.New("workspace id is required")
	}
	if err := s.validateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionGitOpsManageSource); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	endpoint := strings.TrimSpace(input.Endpoint)
	if endpoint == "" {
		return nil, errors.New("endpoint is required")
	}
	sourceType, err := normalizeSourceType(input.SourceType, false)
	if err != nil {
		return nil, err
	}

	item := &domain.DeliverySource{
		Name:          name,
		SourceType:    sourceType,
		Endpoint:      endpoint,
		DefaultRef:    strings.TrimSpace(input.DefaultRef),
		CredentialRef: strings.TrimSpace(input.CredentialRef),
		WorkspaceID:   uint64PtrOrNil(input.WorkspaceID),
		ProjectID:     uint64PtrOrNil(input.ProjectID),
		Status:        domain.DeliverySourceStatusPending,
	}
	if err := s.sources.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetSource(ctx context.Context, userID uint64, sourceID uint64) (*domain.DeliverySource, error) {
	if s == nil || s.sources == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.sources.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, derefUint64(item.WorkspaceID), derefUint64(item.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateSource(
	ctx context.Context,
	userID uint64,
	sourceID uint64,
	input UpdateSourceInput,
) (*domain.DeliverySource, error) {
	if s == nil || s.sources == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.sources.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, derefUint64(item.WorkspaceID), derefUint64(item.ProjectID), PermissionGitOpsManageSource); err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, errors.New("name is required")
		}
		item.Name = name
	}
	if input.DefaultRef != nil {
		item.DefaultRef = strings.TrimSpace(*input.DefaultRef)
	}
	if input.CredentialRef != nil {
		item.CredentialRef = strings.TrimSpace(*input.CredentialRef)
	}
	if input.Disabled != nil {
		if *input.Disabled {
			item.Status = domain.DeliverySourceStatusDisabled
		} else if item.Status == domain.DeliverySourceStatusDisabled {
			item.Status = domain.DeliverySourceStatusPending
		}
	}

	if err := s.sources.Update(ctx, item); err != nil {
		return nil, err
	}
	return s.sources.GetByID(ctx, sourceID)
}

func (s *Service) VerifySource(
	ctx context.Context,
	userID uint64,
	sourceID uint64,
	requestID string,
) (*domain.DeliveryOperation, error) {
	if s == nil || s.sources == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.sources.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, derefUint64(item.WorkspaceID), derefUint64(item.ProjectID), PermissionGitOpsManageSource); err != nil {
		return nil, err
	}

	now := time.Now()
	item.Status = domain.DeliverySourceStatusReady
	item.LastVerifiedAt = &now
	item.LastErrorMessage = ""
	if err := s.sources.Update(ctx, item); err != nil {
		return nil, err
	}

	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		requestID = fmt.Sprintf("verify-source-%d-%d", sourceID, now.UnixNano())
	}
	operation := &domain.DeliveryOperation{
		RequestID:       requestID,
		OperatorID:      userID,
		DeliveryUnitID:  0,
		ActionType:      domain.DeliveryActionType("verify-source"),
		Status:          domain.DeliveryOperationStatusPending,
		ProgressPercent: 0,
	}
	if s.operations != nil {
		if err := s.operations.Create(ctx, operation); err != nil {
			return nil, err
		}
		_ = s.SetProgress(ctx, operation.ID, OperationProgressSnapshot{
			Percent:   5,
			Message:   "queued",
			UpdatedAt: now,
		})
	}
	return operation, nil
}

func (s *Service) ListTargetGroups(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	projectID uint64,
) ([]domain.ClusterTargetGroup, error) {
	if s == nil || s.targets == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if err := s.validateScope(ctx, userID, workspaceID, projectID, PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return s.targets.ListByScope(ctx, workspaceID, uint64PtrOrNil(projectID))
}

func (s *Service) CreateTargetGroup(
	ctx context.Context,
	userID uint64,
	input CreateTargetGroupInput,
) (*domain.ClusterTargetGroup, error) {
	if s == nil || s.targets == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if input.WorkspaceID == 0 {
		return nil, errors.New("workspace id is required")
	}
	if err := s.validateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionGitOpsOverride); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	item := &domain.ClusterTargetGroup{
		Name:                    name,
		WorkspaceID:             input.WorkspaceID,
		ProjectID:               uint64PtrOrNil(input.ProjectID),
		ClusterRefsJSON:         encodeUint64Array(sanitizeClusterRefs(input.ClusterRefs)),
		ClusterSelectorSnapshot: strings.TrimSpace(input.SelectorSummary),
		Description:             strings.TrimSpace(input.Description),
		Status:                  domain.ClusterTargetGroupStatusActive,
	}
	if err := s.targets.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetTargetGroup(ctx context.Context, userID uint64, targetGroupID uint64) (*domain.ClusterTargetGroup, error) {
	if s == nil || s.targets == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.targets.GetByID(ctx, targetGroupID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, item.WorkspaceID, derefUint64(item.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateTargetGroup(
	ctx context.Context,
	userID uint64,
	targetGroupID uint64,
	input UpdateTargetGroupInput,
) (*domain.ClusterTargetGroup, error) {
	if s == nil || s.targets == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.targets.GetByID(ctx, targetGroupID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, item.WorkspaceID, derefUint64(item.ProjectID), PermissionGitOpsOverride); err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, errors.New("name is required")
		}
		item.Name = name
	}
	if input.ClusterRefs != nil {
		item.ClusterRefsJSON = encodeUint64Array(sanitizeClusterRefs(*input.ClusterRefs))
	}
	if input.SelectorSummary != nil {
		item.ClusterSelectorSnapshot = strings.TrimSpace(*input.SelectorSummary)
	}
	if input.Description != nil {
		item.Description = strings.TrimSpace(*input.Description)
	}
	if input.Disabled != nil {
		if *input.Disabled {
			item.Status = domain.ClusterTargetGroupStatusDisabled
		} else if item.Status == domain.ClusterTargetGroupStatusDisabled {
			item.Status = domain.ClusterTargetGroupStatusActive
		}
	}

	if err := s.targets.Update(ctx, item); err != nil {
		return nil, err
	}
	return s.targets.GetByID(ctx, targetGroupID)
}

func (s *Service) ListDeliveryUnits(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	projectID uint64,
) ([]domain.ApplicationDeliveryUnit, error) {
	if s == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if err := s.validateScope(ctx, userID, workspaceID, projectID, PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return s.units.ListByScope(ctx, workspaceID, uint64PtrOrNil(projectID))
}

func (s *Service) CreateDeliveryUnit(
	ctx context.Context,
	userID uint64,
	input CreateDeliveryUnitInput,
) (*DeliveryUnitDetail, error) {
	if s == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if input.WorkspaceID == 0 {
		return nil, errors.New("workspace id is required")
	}
	if input.SourceID == 0 {
		return nil, errors.New("source id is required")
	}
	if err := s.validateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionGitOpsOverride); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if len(input.Environments) == 0 {
		return nil, errors.New("environments is required")
	}
	if s.sources != nil {
		source, err := s.sources.GetByID(ctx, input.SourceID)
		if err != nil {
			return nil, err
		}
		if derefUint64(source.WorkspaceID) != 0 && derefUint64(source.WorkspaceID) != input.WorkspaceID {
			return nil, errors.New("source workspace mismatch")
		}
	}

	syncMode, err := normalizeSyncMode(input.SyncMode, true)
	if err != nil {
		return nil, err
	}
	unit := &domain.ApplicationDeliveryUnit{
		Name:                 name,
		WorkspaceID:          input.WorkspaceID,
		ProjectID:            uint64PtrOrNil(input.ProjectID),
		SourceID:             input.SourceID,
		SourcePath:           strings.TrimSpace(input.SourcePath),
		DefaultNamespace:     strings.TrimSpace(input.DefaultNamespace),
		SyncMode:             syncMode,
		DesiredRevision:      strings.TrimSpace(input.DesiredRevision),
		DesiredAppVersion:    strings.TrimSpace(input.DesiredAppVersion),
		DesiredConfigVersion: strings.TrimSpace(input.DesiredConfigVersion),
		DeliveryStatus:       domain.DeliveryUnitStatusUnknown,
	}

	detail, err := s.units.CreateWithDetails(
		ctx,
		unit,
		toEnvironmentStages(input.Environments),
		toConfigurationOverlays(input.Overlays),
	)
	if err != nil {
		return nil, err
	}
	return convertDeliveryUnitDetail(detail), nil
}

func (s *Service) GetDeliveryUnit(ctx context.Context, userID uint64, unitID uint64) (*DeliveryUnitDetail, error) {
	if s == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	detail, err := s.units.GetDetailByID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, detail.Unit.WorkspaceID, derefUint64(detail.Unit.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return convertDeliveryUnitDetail(detail), nil
}

func (s *Service) UpdateDeliveryUnit(
	ctx context.Context,
	userID uint64,
	unitID uint64,
	input UpdateDeliveryUnitInput,
) (*DeliveryUnitDetail, error) {
	if s == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if unitID == 0 {
		return nil, errors.New("unit id is required")
	}
	existing, err := s.units.GetDetailByID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, existing.Unit.WorkspaceID, derefUint64(existing.Unit.ProjectID), PermissionGitOpsOverride); err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, errors.New("name is required")
		}
		updates["name"] = name
	}
	if input.SourcePath != nil {
		updates["source_path"] = strings.TrimSpace(*input.SourcePath)
	}
	if input.DefaultNamespace != nil {
		updates["default_namespace"] = strings.TrimSpace(*input.DefaultNamespace)
	}
	if input.SyncMode != nil {
		syncMode, err := normalizeSyncMode(*input.SyncMode, true)
		if err != nil {
			return nil, err
		}
		updates["sync_mode"] = syncMode
	}
	if input.DesiredRevision != nil {
		updates["desired_revision"] = strings.TrimSpace(*input.DesiredRevision)
	}
	if input.DesiredAppVersion != nil {
		updates["desired_app_version"] = strings.TrimSpace(*input.DesiredAppVersion)
	}
	if input.DesiredConfigVersion != nil {
		updates["desired_config_version"] = strings.TrimSpace(*input.DesiredConfigVersion)
	}

	environments := existing.Environments
	if input.Environments != nil {
		environments = toEnvironmentStages(*input.Environments)
	}
	overlays := existing.Overlays
	if input.Overlays != nil {
		overlays = toConfigurationOverlays(*input.Overlays)
	}

	updated, err := s.units.ReplaceDetails(ctx, unitID, updates, environments, overlays)
	if err != nil {
		return nil, err
	}
	return convertDeliveryUnitDetail(updated), nil
}

func (s *Service) GetDeliveryUnitStatus(
	ctx context.Context,
	userID uint64,
	unitID uint64,
	environment string,
) (map[string]any, error) {
	if s == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	detail, err := s.units.GetDetailByID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, detail.Unit.WorkspaceID, derefUint64(detail.Unit.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}

	deliveryStatus := strings.TrimSpace(string(detail.Unit.DeliveryStatus))
	if deliveryStatus == "" {
		deliveryStatus = string(domain.DeliveryUnitStatusUnknown)
	}
	envFilter := strings.TrimSpace(environment)
	envs := make([]map[string]any, 0, len(detail.Environments))
	for i := range detail.Environments {
		envName := strings.TrimSpace(detail.Environments[i].Name)
		if envFilter != "" && !strings.EqualFold(envName, envFilter) {
			continue
		}
		envs = append(envs, map[string]any{
			"environment":    envName,
			"syncStatus":     deliveryStatus,
			"driftStatus":    "unknown",
			"targetCount":    1,
			"succeededCount": 0,
			"failedCount":    0,
		})
	}

	return map[string]any{
		"deliveryUnitId": detail.Unit.ID,
		"deliveryStatus": deliveryStatus,
		"driftStatus":    "unknown",
		"lastSyncedAt":   detail.Unit.LastSyncedAt,
		"environments":   envs,
	}, nil
}

func (s *Service) GetDeliveryUnitDiff(ctx context.Context, userID uint64, unitID uint64, stageID uint64) (map[string]any, error) {
	if s == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.units.GetByID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, item.WorkspaceID, derefUint64(item.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}

	diffPayload := "{}"
	placeholder := true
	if s.diffService != nil {
		resolved, resolveErr := s.diffService.GetOrBuild(ctx, unitID, stageID)
		if resolveErr != nil {
			return nil, resolveErr
		}
		diffPayload = resolved
		placeholder = false
	} else {
		cached, cacheErr := s.GetDiff(ctx, unitID, stageID)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if strings.TrimSpace(cached) != "" {
			diffPayload = cached
			placeholder = false
		}
	}
	if strings.TrimSpace(diffPayload) == "" {
		diffPayload = "{}"
		placeholder = true
	}
	return map[string]any{
		"deliveryUnitId": unitID,
		"stageId":        stageID,
		"diff":           diffPayload,
		"placeholder":    placeholder,
	}, nil
}

func (s *Service) ListReleaseRevisions(ctx context.Context, userID uint64, unitID uint64) ([]domain.ReleaseRevision, error) {
	if s == nil || s.units == nil || s.revisions == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.units.GetByID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, item.WorkspaceID, derefUint64(item.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return s.revisions.ListByDeliveryUnit(ctx, unitID)
}

func (s *Service) SubmitAction(ctx context.Context, userID uint64, req ActionRequest) (*domain.DeliveryOperation, error) {
	return s.Execute(ctx, userID, req)
}

func (s *Service) Execute(ctx context.Context, userID uint64, req ActionRequest) (*domain.DeliveryOperation, error) {
	if s == nil || s.units == nil || s.operations == nil {
		return nil, ErrGitOpsNotConfigured
	}
	if req.DeliveryUnitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}

	unit, err := s.units.GetByID(ctx, req.DeliveryUnitID)
	if err != nil {
		return nil, err
	}
	permission := permissionByAction(req.ActionType)
	if err := s.validateScope(ctx, userID, unit.WorkspaceID, derefUint64(unit.ProjectID), permission); err != nil {
		return nil, err
	}

	lockScope := fmt.Sprintf("unit:%d", req.DeliveryUnitID)
	token, locked, lockErr := s.AcquireOperationLock(ctx, lockScope, s.defaultLock)
	if lockErr != nil {
		return nil, lockErr
	}
	if !locked {
		return nil, errors.New("delivery action is locked")
	}
	defer func() {
		_ = s.ReleaseOperationLock(ctx, lockScope, token)
	}()

	requestID := strings.TrimSpace(req.RequestID)
	if requestID == "" {
		requestID = fmt.Sprintf("gitops-%d-%d", req.DeliveryUnitID, time.Now().UnixNano())
	} else {
		existing, lookupErr := s.operations.GetByRequestID(ctx, requestID)
		switch {
		case lookupErr == nil:
			return existing, nil
		case !errors.Is(lookupErr, gorm.ErrRecordNotFound):
			return nil, lookupErr
		}
	}
	actionType := req.ActionType
	if actionType == "" {
		actionType = domain.DeliveryActionTypeSync
	}

	operation := &domain.DeliveryOperation{
		RequestID:          requestID,
		OperatorID:         userID,
		DeliveryUnitID:     req.DeliveryUnitID,
		EnvironmentStageID: req.EnvironmentStageID,
		ActionType:         actionType,
		TargetReleaseID:    req.TargetReleaseID,
		PayloadJSON:        strings.TrimSpace(req.PayloadJSON),
		Status:             domain.DeliveryOperationStatusPending,
		ProgressPercent:    0,
	}
	if err := s.operations.Create(ctx, operation); err != nil {
		existing, lookupErr := s.operations.GetByRequestID(ctx, requestID)
		if lookupErr == nil {
			return existing, nil
		}
		return nil, err
	}
	_ = s.SetProgress(ctx, operation.ID, OperationProgressSnapshot{
		Percent:   5,
		Message:   "queued",
		UpdatedAt: time.Now(),
	})
	if s.queue != nil {
		if err := s.queue.Enqueue(ctx, operation.ID); err != nil {
			_ = s.operations.UpdateStatus(ctx, operation.ID, domain.DeliveryOperationStatusFailed, 0, "queue enqueue failed", err.Error())
			return nil, err
		}
	}
	return operation, nil
}

func (s *Service) GetOperation(ctx context.Context, userID uint64, operationID uint64) (*domain.DeliveryOperation, error) {
	if s == nil || s.operations == nil || s.units == nil {
		return nil, ErrGitOpsNotConfigured
	}
	item, err := s.operations.GetByID(ctx, operationID)
	if err != nil {
		return nil, err
	}
	unit, err := s.units.GetByID(ctx, item.DeliveryUnitID)
	if err != nil {
		return nil, err
	}
	if err := s.validateScope(ctx, userID, unit.WorkspaceID, derefUint64(unit.ProjectID), PermissionGitOpsRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) AcquireOperationLock(ctx context.Context, scope string, ttl time.Duration) (string, bool, error) {
	if strings.TrimSpace(scope) == "" {
		return "", false, errors.New("lock scope is required")
	}
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	if s == nil || s.locks == nil {
		return token, true, nil
	}
	ok, err := s.locks.Acquire(ctx, scope, token, ttl)
	return token, ok, err
}

func (s *Service) ReleaseOperationLock(ctx context.Context, scope string, token string) error {
	if s == nil || s.locks == nil {
		return nil
	}
	return s.locks.Release(ctx, scope, token)
}

func (s *Service) SetProgress(ctx context.Context, operationID uint64, snapshot OperationProgressSnapshot) error {
	if s == nil || s.progress == nil {
		return nil
	}
	return s.progress.SetOperationProgress(ctx, operationID, snapshot)
}

func (s *Service) GetProgress(ctx context.Context, operationID uint64) (OperationProgressSnapshot, error) {
	if s == nil || s.progress == nil {
		return OperationProgressSnapshot{}, nil
	}
	return s.progress.GetOperationProgress(ctx, operationID)
}

func (s *Service) SetDiff(ctx context.Context, deliveryUnitID uint64, stageID uint64, payload string) error {
	if s == nil || s.diffs == nil {
		return nil
	}
	return s.diffs.SetDeliveryUnitDiff(ctx, deliveryUnitID, stageID, payload)
}

func (s *Service) GetDiff(ctx context.Context, deliveryUnitID uint64, stageID uint64) (string, error) {
	if s == nil || s.diffs == nil {
		return "", nil
	}
	return s.diffs.GetDeliveryUnitDiff(ctx, deliveryUnitID, stageID)
}

func (s *Service) validateScope(ctx context.Context, userID uint64, workspaceID uint64, projectID uint64, permission string) error {
	if s == nil || s.scope == nil {
		return nil
	}
	return s.scope.ValidateScope(ctx, userID, workspaceID, projectID, permission)
}

func permissionByAction(action domain.DeliveryActionType) string {
	switch action {
	case domain.DeliveryActionTypePromote:
		return PermissionGitOpsPromote
	case domain.DeliveryActionTypeRollback:
		return PermissionGitOpsRollback
	case domain.DeliveryActionTypeInstall,
		domain.DeliveryActionTypeSync,
		domain.DeliveryActionTypeResync,
		domain.DeliveryActionTypeUpgrade,
		domain.DeliveryActionTypePause,
		domain.DeliveryActionTypeResume,
		domain.DeliveryActionTypeUninstall:
		return PermissionGitOpsSync
	default:
		return PermissionGitOpsRead
	}
}

func convertDeliveryUnitDetail(detail *repository.DeliveryUnitDetail) *DeliveryUnitDetail {
	if detail == nil {
		return nil
	}
	return &DeliveryUnitDetail{
		Unit:         detail.Unit,
		Environments: detail.Environments,
		Overlays:     detail.Overlays,
	}
}

func normalizeSourceType(raw string, allowEmpty bool) (domain.DeliverySourceType, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" {
		if allowEmpty {
			return "", nil
		}
		return "", errors.New("source type is required")
	}
	if trimmed != string(domain.DeliverySourceTypeGit) && trimmed != string(domain.DeliverySourceTypePackage) {
		return "", errors.New("invalid source type")
	}
	return domain.DeliverySourceType(trimmed), nil
}

func normalizeSourceStatus(raw string, allowEmpty bool) (domain.DeliverySourceStatus, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" {
		if allowEmpty {
			return "", nil
		}
		return "", errors.New("status is required")
	}
	switch domain.DeliverySourceStatus(trimmed) {
	case domain.DeliverySourceStatusPending,
		domain.DeliverySourceStatusReady,
		domain.DeliverySourceStatusFailed,
		domain.DeliverySourceStatusDisabled:
		return domain.DeliverySourceStatus(trimmed), nil
	default:
		return "", errors.New("invalid status")
	}
}

func normalizeSyncMode(raw string, allowEmpty bool) (domain.DeliverySyncMode, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" {
		if allowEmpty {
			return domain.DeliverySyncModeManual, nil
		}
		return "", errors.New("sync mode is required")
	}
	if trimmed != string(domain.DeliverySyncModeManual) && trimmed != string(domain.DeliverySyncModeAuto) {
		return "", errors.New("invalid sync mode")
	}
	return domain.DeliverySyncMode(trimmed), nil
}

func normalizePromotionMode(raw string) domain.PromotionMode {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == string(domain.PromotionModeAutomatic) {
		return domain.PromotionModeAutomatic
	}
	return domain.PromotionModeManual
}

func normalizeOverlayType(raw string) domain.ConfigurationOverlayType {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	switch domain.ConfigurationOverlayType(trimmed) {
	case domain.ConfigurationOverlayTypePatch, domain.ConfigurationOverlayTypeManifestSnippet:
		return domain.ConfigurationOverlayType(trimmed)
	default:
		return domain.ConfigurationOverlayTypeValues
	}
}

func toEnvironmentStages(inputs []EnvironmentStageInput) []domain.EnvironmentStage {
	stages := make([]domain.EnvironmentStage, 0, len(inputs))
	for i := range inputs {
		stage := inputs[i]
		stages = append(stages, domain.EnvironmentStage{
			Name:          strings.TrimSpace(stage.Name),
			OrderIndex:    stage.OrderIndex,
			TargetGroupID: stage.TargetGroupID,
			PromotionMode: normalizePromotionMode(stage.PromotionMode),
			Paused:        stage.Paused,
			Status:        domain.EnvironmentStageStatusIdle,
		})
	}
	return stages
}

func toConfigurationOverlays(inputs []ConfigurationOverlayInput) []domain.ConfigurationOverlay {
	overlays := make([]domain.ConfigurationOverlay, 0, len(inputs))
	for i := range inputs {
		overlay := inputs[i]
		overlays = append(overlays, domain.ConfigurationOverlay{
			OverlayType:        normalizeOverlayType(overlay.OverlayType),
			OverlayRef:         strings.TrimSpace(overlay.OverlayRef),
			Precedence:         overlay.Precedence,
			EffectiveScopeJSON: strings.TrimSpace(overlay.EffectiveScope),
		})
	}
	return overlays
}

func sanitizeClusterRefs(clusterRefs []uint64) []uint64 {
	if len(clusterRefs) == 0 {
		return []uint64{}
	}
	uniq := make(map[uint64]struct{}, len(clusterRefs))
	res := make([]uint64, 0, len(clusterRefs))
	for i := range clusterRefs {
		if clusterRefs[i] == 0 {
			continue
		}
		if _, exists := uniq[clusterRefs[i]]; exists {
			continue
		}
		uniq[clusterRefs[i]] = struct{}{}
		res = append(res, clusterRefs[i])
	}
	return res
}

func encodeUint64Array(values []uint64) string {
	data, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func uint64PtrOrNil(value uint64) *uint64 {
	if value == 0 {
		return nil
	}
	v := value
	return &v
}

func derefUint64(value *uint64) uint64 {
	if value == nil {
		return 0
	}
	return *value
}
