package marketplace

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type DistributionCoordinator struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewDistributionCoordinator(rdb *redis.Client) *DistributionCoordinator {
	return &DistributionCoordinator{rdb: rdb, ttl: 15 * time.Minute}
}

func (c *DistributionCoordinator) Mark(ctx context.Context, key, value string) error {
	if c == nil || c.rdb == nil || key == "" {
		return nil
	}
	return c.rdb.Set(ctx, "marketplace:distribution:"+key, value, c.ttl).Err()
}
