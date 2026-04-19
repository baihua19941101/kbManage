package sre

import (
	"context"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type UpgradeCoordinator struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewUpgradeCoordinator(rdb *redis.Client) *UpgradeCoordinator {
	return &UpgradeCoordinator{rdb: rdb, ttl: 15 * time.Minute}
}

func (c *UpgradeCoordinator) Mark(ctx context.Context, key, value string) error {
	if c == nil || c.rdb == nil || key == "" {
		return nil
	}
	return c.rdb.Set(ctx, repository.PlatformSREUpgradeKey(key), value, c.ttl).Err()
}
