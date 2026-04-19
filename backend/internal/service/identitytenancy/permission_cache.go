package identitytenancy

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type PermissionCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewPermissionCache(rdb *redis.Client) *PermissionCache {
	return &PermissionCache{rdb: rdb, ttl: 5 * time.Minute}
}

func (c *PermissionCache) Mark(ctx context.Context, userID uint64, version string) error {
	if c == nil || c.rdb == nil || userID == 0 {
		return nil
	}
	return c.rdb.Set(ctx, repositoryKey("permission", strconv.FormatUint(userID, 10)), version, c.ttl).Err()
}
