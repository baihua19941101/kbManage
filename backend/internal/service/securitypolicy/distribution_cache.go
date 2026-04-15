package securitypolicy

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type DistributionTaskSnapshot struct {
	Status         string    `json:"status"`
	TargetCount    int       `json:"targetCount"`
	SucceededCount int       `json:"succeededCount"`
	FailedCount    int       `json:"failedCount"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type DistributionCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewDistributionCache(client *redis.Client, ttl time.Duration) *DistributionCache {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &DistributionCache{client: client, ttl: ttl}
}

func (c *DistributionCache) SetTaskSnapshot(ctx context.Context, taskID uint64, snapshot DistributionTaskSnapshot) error {
	if c == nil || c.client == nil || taskID == 0 {
		return nil
	}
	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now()
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, distributionTaskKey(taskID), payload, c.ttl).Err()
}

func (c *DistributionCache) GetTaskSnapshot(ctx context.Context, taskID uint64) (DistributionTaskSnapshot, error) {
	if c == nil || c.client == nil || taskID == 0 {
		return DistributionTaskSnapshot{}, nil
	}
	payload, err := c.client.Get(ctx, distributionTaskKey(taskID)).Bytes()
	if err == redis.Nil {
		return DistributionTaskSnapshot{}, nil
	}
	if err != nil {
		return DistributionTaskSnapshot{}, err
	}
	var snapshot DistributionTaskSnapshot
	if err := json.Unmarshal(payload, &snapshot); err != nil {
		return DistributionTaskSnapshot{}, err
	}
	return snapshot, nil
}

func distributionTaskKey(taskID uint64) string {
	return repository.PolicyDistributionKey("task", strconv.FormatUint(taskID, 10))
}
