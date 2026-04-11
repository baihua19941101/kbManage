package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type AuditExportStatus string

const (
	AuditExportStatusPending   AuditExportStatus = "pending"
	AuditExportStatusRunning   AuditExportStatus = "running"
	AuditExportStatusSucceeded AuditExportStatus = "succeeded"
	AuditExportStatusFailed    AuditExportStatus = "failed"
)

type AuditExportTask struct {
	ID         string            `json:"id"`
	OperatorID uint64            `json:"operatorId"`
	Status     AuditExportStatus `json:"status"`
	Filters    AuditQuery        `json:"-"`

	ResultTotal  int        `json:"resultTotal"`
	DownloadURL  string     `json:"downloadUrl,omitempty"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
}

type AuditExportArtifact struct {
	TaskID      string
	FileName    string
	ContentType string
	Data        []byte
	GeneratedAt time.Time
}

type AuditExportRepository struct {
	db *gorm.DB

	mu        sync.RWMutex
	nextID    uint64
	tasks     map[string]AuditExportTask
	artifacts map[string]AuditExportArtifact
	queue     chan string
}

func NewAuditExportRepository(db *gorm.DB) *AuditExportRepository {
	return &AuditExportRepository{
		db:        db,
		nextID:    1,
		tasks:     make(map[string]AuditExportTask),
		artifacts: make(map[string]AuditExportArtifact),
		queue:     make(chan string, 256),
	}
}

func (r *AuditExportRepository) Create(ctx context.Context, operatorID uint64, filters AuditQuery) (*AuditExportTask, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	taskID := fmt.Sprintf("aexp-%d", r.nextID)
	r.nextID++

	task := AuditExportTask{
		ID:         taskID,
		OperatorID: operatorID,
		Status:     AuditExportStatusPending,
		Filters:    filters,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	r.tasks[task.ID] = task

	copyTask := task
	return &copyTask, nil
}

func (r *AuditExportRepository) Get(ctx context.Context, taskID string) (*AuditExportTask, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[taskID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyTask := task
	return &copyTask, nil
}

func (r *AuditExportRepository) Enqueue(ctx context.Context, taskID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case r.queue <- taskID:
		return nil
	}
}

func (r *AuditExportRepository) Dequeue(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case taskID := <-r.queue:
		return taskID, nil
	}
}

func (r *AuditExportRepository) MarkRunning(ctx context.Context, taskID string) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[taskID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	task.Status = AuditExportStatusRunning
	task.UpdatedAt = time.Now()
	r.tasks[taskID] = task
	return nil
}

func (r *AuditExportRepository) MarkSucceeded(ctx context.Context, taskID string, total int, downloadURL string) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[taskID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	now := time.Now()
	task.Status = AuditExportStatusSucceeded
	task.ResultTotal = total
	task.DownloadURL = downloadURL
	task.ErrorMessage = ""
	task.CompletedAt = &now
	task.UpdatedAt = now
	r.tasks[taskID] = task
	return nil
}

func (r *AuditExportRepository) SaveArtifact(ctx context.Context, taskID, fileName, contentType string, data []byte) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[taskID]; !ok {
		return gorm.ErrRecordNotFound
	}
	r.artifacts[taskID] = AuditExportArtifact{
		TaskID:      taskID,
		FileName:    fileName,
		ContentType: contentType,
		Data:        append([]byte(nil), data...),
		GeneratedAt: time.Now(),
	}
	return nil
}

func (r *AuditExportRepository) GetArtifact(ctx context.Context, taskID string) (*AuditExportArtifact, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.artifacts[taskID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyItem := item
	copyItem.Data = append([]byte(nil), item.Data...)
	return &copyItem, nil
}

func (r *AuditExportRepository) MarkFailed(ctx context.Context, taskID string, message string) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[taskID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	now := time.Now()
	task.Status = AuditExportStatusFailed
	task.ErrorMessage = message
	task.CompletedAt = &now
	task.UpdatedAt = now
	r.tasks[taskID] = task
	return nil
}
