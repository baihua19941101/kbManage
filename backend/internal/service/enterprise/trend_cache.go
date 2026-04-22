package enterprise

import (
	"context"
	"encoding/json"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type TrendCache struct{ rdb *redis.Client }

func NewTrendCache(rdb *redis.Client) *TrendCache { return &TrendCache{rdb: rdb} }

func (c *TrendCache) Store(ctx context.Context, key string, value any) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, repository.PrefixedRedisKey("enterprise:trend", key), data, 10*time.Minute).Err()
}
