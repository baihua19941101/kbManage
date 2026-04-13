package workloadops

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionCache(client *redis.Client, ttl time.Duration) *SessionCache {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return &SessionCache{client: client, ttl: ttl}
}

func (c *SessionCache) SetSessionToken(ctx context.Context, sessionKey, token string) error {
	if c == nil || c.client == nil || sessionKey == "" {
		return nil
	}
	return c.client.Set(ctx, "workloadops:terminal:session:"+sessionKey, token, c.ttl).Err()
}

func (c *SessionCache) GetSessionToken(ctx context.Context, sessionKey string) (string, error) {
	if c == nil || c.client == nil || sessionKey == "" {
		return "", nil
	}
	value, err := c.client.Get(ctx, "workloadops:terminal:session:"+sessionKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}
