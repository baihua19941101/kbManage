package enterprise

import (
	"context"
	"encoding/json"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ReportCache struct{ rdb *redis.Client }

func NewReportCache(rdb *redis.Client) *ReportCache { return &ReportCache{rdb: rdb} }

func (c *ReportCache) Store(ctx context.Context, key string, value any) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, repository.PrefixedRedisKey("enterprise:report", key), data, 10*time.Minute).Err()
}
