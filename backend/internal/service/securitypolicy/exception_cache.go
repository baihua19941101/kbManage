package securitypolicy

import (
	"context"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ExceptionCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewExceptionCache(client *redis.Client, ttl time.Duration) *ExceptionCache {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return &ExceptionCache{client: client, ttl: ttl}
}

func (c *ExceptionCache) SetExceptionStatus(ctx context.Context, exceptionID uint64, status string) error {
	if c == nil || c.client == nil || exceptionID == 0 {
		return nil
	}
	return c.client.Set(ctx, repository.PolicyExceptionKey("status", strconv.FormatUint(exceptionID, 10)), status, c.ttl).Err()
}

func (c *ExceptionCache) GetExceptionStatus(ctx context.Context, exceptionID uint64) (string, error) {
	if c == nil || c.client == nil || exceptionID == 0 {
		return "", nil
	}
	value, err := c.client.Get(ctx, repository.PolicyExceptionKey("status", strconv.FormatUint(exceptionID, 10))).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}
