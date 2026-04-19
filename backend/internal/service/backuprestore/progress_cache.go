package backuprestore

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProgressCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewProgressCache(rdb *redis.Client) *ProgressCache {
	return &ProgressCache{rdb: rdb, ttl: 10 * time.Minute}
}

func (c *ProgressCache) Set(ctx context.Context, targetType string, targetID uint64, operation, status string) error {
	if c == nil || c.rdb == nil || targetID == 0 {
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"targetType": targetType,
		"targetId":   targetID,
		"operation":  operation,
		"status":     status,
	})
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, backupRestoreProgressKey(targetType, strconv.FormatUint(targetID, 10)), payload, c.ttl).Err()
}
