package observability

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrQueryCacheUnavailable = errors.New("observability query cache is not configured")

type QueryCache struct {
	client *redis.Client
	prefix string
}

func NewQueryCache(client *redis.Client, prefix string) *QueryCache {
	if prefix == "" {
		prefix = "observability:query"
	}
	return &QueryCache{
		client: client,
		prefix: prefix,
	}
}

func (c *QueryCache) Set(ctx context.Context, key string, payload string, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return ErrQueryCacheUnavailable
	}
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return c.client.Set(ctx, c.buildKey(key), payload, ttl).Err()
}

func (c *QueryCache) Get(ctx context.Context, key string) (string, bool, error) {
	if c == nil || c.client == nil {
		return "", false, ErrQueryCacheUnavailable
	}
	value, err := c.client.Get(ctx, c.buildKey(key)).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (c *QueryCache) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}
