package gitops

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type OperationProgressSnapshot struct {
	Percent   int       `json:"percent"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProgressCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewProgressCache(client *redis.Client, ttl time.Duration) *ProgressCache {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &ProgressCache{client: client, ttl: ttl}
}

func (c *ProgressCache) SetOperationProgress(ctx context.Context, operationID uint64, snapshot OperationProgressSnapshot) error {
	if c == nil || c.client == nil {
		return nil
	}
	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now()
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, progressKey(operationID), payload, c.ttl).Err()
}

func (c *ProgressCache) GetOperationProgress(ctx context.Context, operationID uint64) (OperationProgressSnapshot, error) {
	if c == nil || c.client == nil {
		return OperationProgressSnapshot{}, nil
	}
	value, err := c.client.Get(ctx, progressKey(operationID)).Bytes()
	if err == redis.Nil {
		return OperationProgressSnapshot{}, nil
	}
	if err != nil {
		return OperationProgressSnapshot{}, err
	}
	var snapshot OperationProgressSnapshot
	if err := json.Unmarshal(value, &snapshot); err != nil {
		return OperationProgressSnapshot{}, err
	}
	return snapshot, nil
}

func (c *ProgressCache) DeleteOperationProgress(ctx context.Context, operationID uint64) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Del(ctx, progressKey(operationID)).Err()
}

func progressKey(operationID uint64) string {
	return repository.GitOpsProgressKey("operation", strconv.FormatUint(operationID, 10))
}
