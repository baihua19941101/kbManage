package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type BatchOperationRepository struct {
	db *gorm.DB

	mu       sync.RWMutex
	nextTask uint64
	nextItem uint64
	tasks    map[uint64]domain.BatchOperationTask
	items    map[uint64][]domain.BatchOperationItem
}

func NewBatchOperationRepository(db *gorm.DB) *BatchOperationRepository {
	return &BatchOperationRepository{
		db:       db,
		nextTask: 1,
		nextItem: 1,
		tasks:    make(map[uint64]domain.BatchOperationTask),
		items:    make(map[uint64][]domain.BatchOperationItem),
	}
}

func (r *BatchOperationRepository) CreateTask(ctx context.Context, task *domain.BatchOperationTask) error {
	if task == nil {
		return errors.New("batch task is required")
	}
	if r.db != nil {
		return r.db.WithContext(ctx).Create(task).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	copyTask := *task
	if copyTask.ID == 0 {
		copyTask.ID = r.nextTask
		r.nextTask++
	}
	now := time.Now()
	if copyTask.CreatedAt.IsZero() {
		copyTask.CreatedAt = now
	}
	copyTask.UpdatedAt = now
	r.tasks[copyTask.ID] = copyTask
	*task = copyTask
	return nil
}

func (r *BatchOperationRepository) CreateItems(ctx context.Context, taskID uint64, items []domain.BatchOperationItem) error {
	if r.db != nil {
		if len(items) == 0 {
			return nil
		}
		return r.db.WithContext(ctx).Create(&items).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	stored := r.items[taskID]
	for _, item := range items {
		copyItem := item
		copyItem.BatchID = taskID
		if copyItem.ID == 0 {
			copyItem.ID = r.nextItem
			r.nextItem++
		}
		if copyItem.CreatedAt.IsZero() {
			copyItem.CreatedAt = now
		}
		copyItem.UpdatedAt = now
		stored = append(stored, copyItem)
	}
	r.items[taskID] = stored
	return nil
}

func (r *BatchOperationRepository) GetTaskByID(ctx context.Context, id uint64) (*domain.BatchOperationTask, []domain.BatchOperationItem, error) {
	if r.db != nil {
		var task domain.BatchOperationTask
		if err := r.db.WithContext(ctx).First(&task, id).Error; err != nil {
			return nil, nil, err
		}
		var items []domain.BatchOperationItem
		if err := r.db.WithContext(ctx).Where("batch_id = ?", id).Order("id asc").Find(&items).Error; err != nil {
			return nil, nil, err
		}
		return &task, items, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return nil, nil, gorm.ErrRecordNotFound
	}
	copyTask := task
	copyItems := append([]domain.BatchOperationItem(nil), r.items[id]...)
	return &copyTask, copyItems, nil
}

func (r *BatchOperationRepository) UpdateItemResult(
	ctx context.Context,
	batchID uint64,
	itemID uint64,
	status domain.BatchOperationItemStatus,
	actionRequestID *uint64,
	resultMessage string,
	failureReason string,
) error {
	now := time.Now()
	if r.db != nil {
		updates := map[string]any{
			"status":            status,
			"action_request_id": actionRequestID,
			"result_message":    resultMessage,
			"failure_reason":    failureReason,
			"updated_at":        now,
		}
		if status == domain.BatchOperationItemStatusRunning {
			updates["started_at"] = now
		}
		if status == domain.BatchOperationItemStatusSucceeded || status == domain.BatchOperationItemStatusFailed || status == domain.BatchOperationItemStatusCanceled {
			updates["completed_at"] = now
		}
		res := r.db.WithContext(ctx).Model(&domain.BatchOperationItem{}).Where("id = ? AND batch_id = ?", itemID, batchID).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	items := r.items[batchID]
	for idx := range items {
		if items[idx].ID != itemID {
			continue
		}
		items[idx].Status = status
		items[idx].ActionRequestID = actionRequestID
		items[idx].ResultMessage = resultMessage
		items[idx].FailureReason = failureReason
		items[idx].UpdatedAt = now
		if status == domain.BatchOperationItemStatusRunning {
			items[idx].StartedAt = &now
		}
		if status == domain.BatchOperationItemStatusSucceeded || status == domain.BatchOperationItemStatusFailed || status == domain.BatchOperationItemStatusCanceled {
			items[idx].CompletedAt = &now
		}
		r.items[batchID] = items
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (r *BatchOperationRepository) UpdateTaskSummary(
	ctx context.Context,
	taskID uint64,
	status domain.BatchOperationStatus,
	succeeded int,
	failed int,
	canceled int,
	progress int,
) error {
	now := time.Now()
	if r.db != nil {
		updates := map[string]any{
			"status":            status,
			"succeeded_targets": succeeded,
			"failed_targets":    failed,
			"canceled_targets":  canceled,
			"progress_percent":  progress,
			"updated_at":        now,
		}
		if status == domain.BatchOperationStatusRunning {
			updates["started_at"] = now
		}
		if status == domain.BatchOperationStatusSucceeded || status == domain.BatchOperationStatusPartiallySucceeded || status == domain.BatchOperationStatusFailed || status == domain.BatchOperationStatusCanceled {
			updates["completed_at"] = now
		}
		res := r.db.WithContext(ctx).Model(&domain.BatchOperationTask{}).Where("id = ?", taskID).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	task.Status = status
	task.SucceededTargets = succeeded
	task.FailedTargets = failed
	task.CanceledTargets = canceled
	task.ProgressPercent = progress
	task.UpdatedAt = now
	if status == domain.BatchOperationStatusRunning {
		task.StartedAt = &now
	}
	if status == domain.BatchOperationStatusSucceeded || status == domain.BatchOperationStatusPartiallySucceeded || status == domain.BatchOperationStatusFailed || status == domain.BatchOperationStatusCanceled {
		task.CompletedAt = &now
	}
	r.tasks[taskID] = task
	return nil
}
