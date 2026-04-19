package sre

import (
	"context"
	"encoding/json"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type HealthCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewHealthCache(rdb *redis.Client) *HealthCache {
	return &HealthCache{rdb: rdb, ttl: 30 * time.Second}
}

func (c *HealthCache) Store(ctx context.Context, key string, value any) error {
	if c == nil || c.rdb == nil || key == "" {
		return nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, repository.PlatformSREHealthKey(key), payload, c.ttl).Err()
}
