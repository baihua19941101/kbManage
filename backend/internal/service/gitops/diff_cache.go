package gitops

import (
	"context"
	"strconv"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type DiffCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewDiffCache(client *redis.Client, ttl time.Duration) *DiffCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &DiffCache{client: client, ttl: ttl}
}

func (c *DiffCache) SetDeliveryUnitDiff(ctx context.Context, deliveryUnitID uint64, stageID uint64, payload string) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Set(ctx, diffKey(deliveryUnitID, stageID), payload, c.ttl).Err()
}

func (c *DiffCache) GetDeliveryUnitDiff(ctx context.Context, deliveryUnitID uint64, stageID uint64) (string, error) {
	if c == nil || c.client == nil {
		return "", nil
	}
	value, err := c.client.Get(ctx, diffKey(deliveryUnitID, stageID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

func (c *DiffCache) DeleteDeliveryUnitDiff(ctx context.Context, deliveryUnitID uint64, stageID uint64) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Del(ctx, diffKey(deliveryUnitID, stageID)).Err()
}

func diffKey(deliveryUnitID uint64, stageID uint64) string {
	return repository.GitOpsDiffKey(
		"unit",
		strconv.FormatUint(deliveryUnitID, 10),
		"stage",
		strconv.FormatUint(stageID, 10),
	)
}
