package identitytenancy

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RevocationCoordinator struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRevocationCoordinator(rdb *redis.Client) *RevocationCoordinator {
	return &RevocationCoordinator{rdb: rdb, ttl: 15 * time.Minute}
}

func (c *RevocationCoordinator) Mark(ctx context.Context, userID uint64, reason string) error {
	if c == nil || c.rdb == nil || userID == 0 {
		return nil
	}
	return c.rdb.Set(ctx, repositoryKey("revocation", strconv.FormatUint(userID, 10)), reason, c.ttl).Err()
}

func repositoryKey(parts ...string) string {
	key := "identitytenancy"
	for _, part := range parts {
		if part == "" {
			continue
		}
		key += ":" + part
	}
	return key
}
