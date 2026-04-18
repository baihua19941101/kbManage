package clusterlifecycle

import (
	"context"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type OperationLock struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewOperationLock(rdb *redis.Client) *OperationLock {
	return &OperationLock{rdb: rdb, ttl: 2 * time.Minute}
}

func (l *OperationLock) Acquire(ctx context.Context, clusterID uint64, action string) (bool, error) {
	if l == nil || l.rdb == nil || clusterID == 0 {
		return true, nil
	}
	return l.rdb.SetNX(ctx, repository.ClusterLifecycleLockKey(strconv.FormatUint(clusterID, 10)), action, l.ttl).Result()
}

func (l *OperationLock) Release(ctx context.Context, clusterID uint64) error {
	if l == nil || l.rdb == nil || clusterID == 0 {
		return nil
	}
	return l.rdb.Del(ctx, repository.ClusterLifecycleLockKey(strconv.FormatUint(clusterID, 10))).Err()
}
