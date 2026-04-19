package marketplace

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CatalogCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewCatalogCache(rdb *redis.Client) *CatalogCache {
	return &CatalogCache{rdb: rdb, ttl: 10 * time.Minute}
}

func (c *CatalogCache) Store(ctx context.Context, key string, value any) error {
	if c == nil || c.rdb == nil || key == "" {
		return nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, "marketplace:catalog:"+key, payload, c.ttl).Err()
}
