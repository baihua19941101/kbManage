package backuprestore

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type PrecheckCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewPrecheckCache(rdb *redis.Client) *PrecheckCache {
	return &PrecheckCache{rdb: rdb, ttl: 15 * time.Minute}
}

func (c *PrecheckCache) Store(ctx context.Context, restorePointID uint64, result *PrecheckResult) error {
	if c == nil || c.rdb == nil || restorePointID == 0 || result == nil {
		return nil
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, backupRestorePrecheckKey(strconv.FormatUint(restorePointID, 10)), payload, c.ttl).Err()
}
