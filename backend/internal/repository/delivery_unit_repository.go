package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DeliveryUnitDetail struct {
	Unit         domain.ApplicationDeliveryUnit
	Environments []domain.EnvironmentStage
	Overlays     []domain.ConfigurationOverlay
}

type DeliveryUnitStore interface {
	Create(ctx context.Context, item *domain.ApplicationDeliveryUnit) error
	CreateWithDetails(
		ctx context.Context,
		unit *domain.ApplicationDeliveryUnit,
		environments []domain.EnvironmentStage,
		overlays []domain.ConfigurationOverlay,
	) (*DeliveryUnitDetail, error)
	GetByID(ctx context.Context, id uint64) (*domain.ApplicationDeliveryUnit, error)
	GetDetailByID(ctx context.Context, id uint64) (*DeliveryUnitDetail, error)
	ListByScope(ctx context.Context, workspaceID uint64, projectID *uint64) ([]domain.ApplicationDeliveryUnit, error)
	ReplaceDetails(
		ctx context.Context,
		unitID uint64,
		updates map[string]any,
		environments []domain.EnvironmentStage,
		overlays []domain.ConfigurationOverlay,
	) (*DeliveryUnitDetail, error)
	Update(ctx context.Context, item *domain.ApplicationDeliveryUnit) error
	UpdateFields(ctx context.Context, unitID uint64, updates map[string]any) error
	ListEnvironmentStages(ctx context.Context, unitID uint64) ([]domain.EnvironmentStage, error)
	GetEnvironmentStageByID(ctx context.Context, stageID uint64) (*domain.EnvironmentStage, error)
	UpdateEnvironmentStage(ctx context.Context, stageID uint64, updates map[string]any) error
}

type DeliveryUnitRepository struct {
	db *gorm.DB
}

func NewDeliveryUnitRepository(db *gorm.DB) *DeliveryUnitRepository {
	return &DeliveryUnitRepository{db: db}
}

func (r *DeliveryUnitRepository) Create(ctx context.Context, item *domain.ApplicationDeliveryUnit) error {
	if item == nil {
		return errors.New("delivery unit is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery unit repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DeliveryUnitRepository) CreateWithDetails(
	ctx context.Context,
	unit *domain.ApplicationDeliveryUnit,
	environments []domain.EnvironmentStage,
	overlays []domain.ConfigurationOverlay,
) (*DeliveryUnitDetail, error) {
	if unit == nil {
		return nil, errors.New("delivery unit is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}

	result := &DeliveryUnitDetail{}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(unit).Error; err != nil {
			return err
		}

		stages, err := createEnvironmentStages(tx, unit.ID, environments)
		if err != nil {
			return err
		}
		savedOverlays, err := createOverlays(tx, unit.ID, stages, overlays)
		if err != nil {
			return err
		}

		result.Unit = *unit
		result.Environments = stages
		result.Overlays = savedOverlays
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DeliveryUnitRepository) GetByID(ctx context.Context, id uint64) (*domain.ApplicationDeliveryUnit, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}
	var item domain.ApplicationDeliveryUnit
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DeliveryUnitRepository) GetDetailByID(ctx context.Context, id uint64) (*DeliveryUnitDetail, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}
	unit, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	stages, err := r.ListEnvironmentStages(ctx, id)
	if err != nil {
		return nil, err
	}
	overlays := make([]domain.ConfigurationOverlay, 0)
	if err := r.db.WithContext(ctx).
		Where("delivery_unit_id = ?", id).
		Order("precedence ASC, id ASC").
		Find(&overlays).Error; err != nil {
		return nil, err
	}
	return &DeliveryUnitDetail{
		Unit:         *unit,
		Environments: stages,
		Overlays:     overlays,
	}, nil
}

func (r *DeliveryUnitRepository) ListByScope(ctx context.Context, workspaceID uint64, projectID *uint64) ([]domain.ApplicationDeliveryUnit, error) {
	if workspaceID == 0 {
		return nil, errors.New("workspace id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.ApplicationDeliveryUnit{}).Where("workspace_id = ?", workspaceID)
	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}
	items := make([]domain.ApplicationDeliveryUnit, 0)
	if err := query.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DeliveryUnitRepository) ReplaceDetails(
	ctx context.Context,
	unitID uint64,
	updates map[string]any,
	environments []domain.EnvironmentStage,
	overlays []domain.ConfigurationOverlay,
) (*DeliveryUnitDetail, error) {
	if unitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}

	result := &DeliveryUnitDetail{}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			res := tx.Model(&domain.ApplicationDeliveryUnit{}).Where("id = ?", unitID).Updates(updates)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}

		if err := tx.Where("delivery_unit_id = ?", unitID).Delete(&domain.ConfigurationOverlay{}).Error; err != nil {
			return err
		}
		if err := tx.Where("delivery_unit_id = ?", unitID).Delete(&domain.EnvironmentStage{}).Error; err != nil {
			return err
		}

		stages, err := createEnvironmentStages(tx, unitID, environments)
		if err != nil {
			return err
		}
		savedOverlays, err := createOverlays(tx, unitID, stages, overlays)
		if err != nil {
			return err
		}

		var unit domain.ApplicationDeliveryUnit
		if err := tx.First(&unit, unitID).Error; err != nil {
			return err
		}

		result.Unit = unit
		result.Environments = stages
		result.Overlays = savedOverlays
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DeliveryUnitRepository) Update(ctx context.Context, item *domain.ApplicationDeliveryUnit) error {
	if item == nil || item.ID == 0 {
		return errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery unit repository is not configured")
	}
	return r.db.WithContext(ctx).Model(&domain.ApplicationDeliveryUnit{}).Where("id = ?", item.ID).Updates(item).Error
}

func (r *DeliveryUnitRepository) UpdateFields(ctx context.Context, unitID uint64, updates map[string]any) error {
	if unitID == 0 {
		return errors.New("delivery unit id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("delivery unit repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.ApplicationDeliveryUnit{}).Where("id = ?", unitID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *DeliveryUnitRepository) ListEnvironmentStages(ctx context.Context, unitID uint64) ([]domain.EnvironmentStage, error) {
	if unitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}
	stages := make([]domain.EnvironmentStage, 0)
	if err := r.db.WithContext(ctx).
		Where("delivery_unit_id = ?", unitID).
		Order("order_index ASC, id ASC").
		Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}

func (r *DeliveryUnitRepository) GetEnvironmentStageByID(ctx context.Context, stageID uint64) (*domain.EnvironmentStage, error) {
	if stageID == 0 {
		return nil, errors.New("environment stage id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("delivery unit repository is not configured")
	}
	var item domain.EnvironmentStage
	if err := r.db.WithContext(ctx).First(&item, stageID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DeliveryUnitRepository) UpdateEnvironmentStage(ctx context.Context, stageID uint64, updates map[string]any) error {
	if stageID == 0 {
		return errors.New("environment stage id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("delivery unit repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.EnvironmentStage{}).Where("id = ?", stageID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func createEnvironmentStages(
	tx *gorm.DB,
	unitID uint64,
	environments []domain.EnvironmentStage,
) ([]domain.EnvironmentStage, error) {
	stages := make([]domain.EnvironmentStage, 0, len(environments))
	for i := range environments {
		stage := environments[i]
		stage.ID = 0
		stage.DeliveryUnitID = unitID
		stage.Name = strings.TrimSpace(stage.Name)
		if stage.Name == "" {
			return nil, errors.New("environment name is required")
		}
		if stage.TargetGroupID == 0 {
			return nil, errors.New("target group id is required")
		}
		if stage.OrderIndex == 0 {
			stage.OrderIndex = i + 1
		}
		if strings.TrimSpace(string(stage.PromotionMode)) == "" {
			stage.PromotionMode = domain.PromotionModeManual
		}
		if strings.TrimSpace(string(stage.Status)) == "" {
			stage.Status = domain.EnvironmentStageStatusIdle
		}
		stages = append(stages, stage)
	}
	if len(stages) == 0 {
		return stages, nil
	}
	if err := tx.Create(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}

func createOverlays(
	tx *gorm.DB,
	unitID uint64,
	environments []domain.EnvironmentStage,
	overlays []domain.ConfigurationOverlay,
) ([]domain.ConfigurationOverlay, error) {
	saved := make([]domain.ConfigurationOverlay, 0, len(overlays))
	if len(overlays) == 0 {
		return saved, nil
	}

	environmentIDByName := make(map[string]uint64, len(environments))
	for i := range environments {
		environmentIDByName[strings.ToLower(strings.TrimSpace(environments[i].Name))] = environments[i].ID
	}

	for i := range overlays {
		item := overlays[i]
		item.ID = 0
		item.DeliveryUnitID = unitID
		item.OverlayRef = strings.TrimSpace(item.OverlayRef)
		if item.OverlayRef == "" {
			return nil, errors.New("overlay ref is required")
		}
		if strings.TrimSpace(string(item.OverlayType)) == "" {
			item.OverlayType = domain.ConfigurationOverlayTypeValues
		}
		if item.EnvironmentStageID == nil {
			if envName, ok := parseEnvironmentNameFromScope(item.EffectiveScopeJSON); ok {
				if stageID, exists := environmentIDByName[envName]; exists {
					item.EnvironmentStageID = &stageID
				}
			}
		}
		saved = append(saved, item)
	}

	if err := tx.Create(&saved).Error; err != nil {
		return nil, err
	}
	return saved, nil
}

func parseEnvironmentNameFromScope(scope string) (string, bool) {
	trimmed := strings.TrimSpace(scope)
	if trimmed == "" {
		return "", false
	}
	lower := strings.ToLower(trimmed)
	if !strings.HasPrefix(lower, "env:") {
		return "", false
	}
	env := strings.TrimSpace(trimmed[4:])
	if env == "" {
		return "", false
	}
	return strings.ToLower(env), true
}
