package enterprise

import (
	"context"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ExportCoordinator struct{ rdb *redis.Client }

func NewExportCoordinator(rdb *redis.Client) *ExportCoordinator { return &ExportCoordinator{rdb: rdb} }

func (c *ExportCoordinator) Lock(ctx context.Context, key string) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Set(ctx, repository.PrefixedRedisKey("enterprise:export", key), "1", 15*time.Minute).Err()
}
