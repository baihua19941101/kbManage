package sre

import (
	"context"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ScaleCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewScaleCache(rdb *redis.Client) *ScaleCache {
	return &ScaleCache{rdb: rdb, ttl: 10 * time.Minute}
}

func (c *ScaleCache) Mark(ctx context.Context, key, value string) error {
	if c == nil || c.rdb == nil || key == "" {
		return nil
	}
	return c.rdb.Set(ctx, repository.PlatformSREScaleKey(key), value, c.ttl).Err()
}
