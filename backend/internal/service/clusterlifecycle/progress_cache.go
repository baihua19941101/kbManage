package clusterlifecycle

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ProgressCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewProgressCache(rdb *redis.Client) *ProgressCache {
	return &ProgressCache{rdb: rdb, ttl: 10 * time.Minute}
}

func (c *ProgressCache) SetOperation(ctx context.Context, clusterID uint64, operationType, status string) error {
	if c == nil || c.rdb == nil || clusterID == 0 {
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"clusterId":      clusterID,
		"operationType":  operationType,
		"status":         status,
		"updatedAtEpoch": time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, repository.ClusterLifecycleProgressKey(strconv.FormatUint(clusterID, 10)), payload, c.ttl).Err()
}

func (c *ProgressCache) Clear(ctx context.Context, clusterID uint64) error {
	if c == nil || c.rdb == nil || clusterID == 0 {
		return nil
	}
	return c.rdb.Del(ctx, repository.ClusterLifecycleProgressKey(strconv.FormatUint(clusterID, 10))).Err()
}
