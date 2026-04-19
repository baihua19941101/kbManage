package backuprestore

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type OperationLock struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewOperationLock(rdb *redis.Client) *OperationLock {
	return &OperationLock{rdb: rdb, ttl: 2 * time.Minute}
}

func (l *OperationLock) Acquire(ctx context.Context, targetType string, targetID uint64) (bool, error) {
	if l == nil || l.rdb == nil || targetID == 0 {
		return true, nil
	}
	return l.rdb.SetNX(ctx, backupRestoreLockKey(targetType, strconv.FormatUint(targetID, 10)), "locked", l.ttl).Result()
}

func (l *OperationLock) Release(ctx context.Context, targetType string, targetID uint64) error {
	if l == nil || l.rdb == nil || targetID == 0 {
		return nil
	}
	return l.rdb.Del(ctx, backupRestoreLockKey(targetType, strconv.FormatUint(targetID, 10))).Err()
}
