package marketplace

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CompatibilityCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewCompatibilityCache(rdb *redis.Client) *CompatibilityCache {
	return &CompatibilityCache{rdb: rdb, ttl: 10 * time.Minute}
}

func (c *CompatibilityCache) Mark(ctx context.Context, key, value string) error {
	if c == nil || c.rdb == nil || key == "" {
		return nil
	}
	return c.rdb.Set(ctx, "marketplace:compatibility:"+key, value, c.ttl).Err()
}
