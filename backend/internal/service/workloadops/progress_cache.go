package workloadops

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProgressCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewProgressCache(client *redis.Client, ttl time.Duration) *ProgressCache {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &ProgressCache{client: client, ttl: ttl}
}

func (c *ProgressCache) SetActionProgress(ctx context.Context, actionID uint64, message string) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Set(ctx, redisKeyActionProgress(actionID), message, c.ttl).Err()
}

func (c *ProgressCache) GetActionProgress(ctx context.Context, actionID uint64) (string, error) {
	if c == nil || c.client == nil {
		return "", nil
	}
	value, err := c.client.Get(ctx, redisKeyActionProgress(actionID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

func redisKeyActionProgress(actionID uint64) string {
	return "workloadops:progress:action:" + strconv.FormatUint(actionID, 10)
}
