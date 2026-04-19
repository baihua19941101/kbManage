package identitytenancy

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewSessionCache(rdb *redis.Client) *SessionCache {
	return &SessionCache{rdb: rdb, ttl: 10 * time.Minute}
}

func (c *SessionCache) Store(ctx context.Context, userID uint64, sessions any) error {
	if c == nil || c.rdb == nil || userID == 0 {
		return nil
	}
	payload, err := json.Marshal(sessions)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, repositoryKey("session", strconv.FormatUint(userID, 10)), payload, c.ttl).Err()
}
