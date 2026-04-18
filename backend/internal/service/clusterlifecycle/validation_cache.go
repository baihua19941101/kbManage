package clusterlifecycle

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ValidationCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewValidationCache(rdb *redis.Client) *ValidationCache {
	return &ValidationCache{rdb: rdb, ttl: 15 * time.Minute}
}

func (c *ValidationCache) Store(ctx context.Context, clusterID uint64, result *ValidationResult) error {
	if c == nil || c.rdb == nil || result == nil {
		return nil
	}
	cacheKey := "global"
	if clusterID != 0 {
		cacheKey = strconv.FormatUint(clusterID, 10)
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, repository.ClusterLifecycleValidationKey(cacheKey), payload, c.ttl).Err()
}
