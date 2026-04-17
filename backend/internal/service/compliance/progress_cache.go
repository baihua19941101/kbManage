package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type ScanProgressSnapshot struct {
	ExecutionID uint64    `json:"executionId"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	Message     string    `json:"message,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProgressCache struct {
	client *redis.Client
	ttl    time.Duration
	mu     sync.RWMutex
	mem    map[uint64]ScanProgressSnapshot
}

func NewProgressCache(client *redis.Client, ttl time.Duration) *ProgressCache {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &ProgressCache{client: client, ttl: ttl, mem: make(map[uint64]ScanProgressSnapshot)}
}

func (c *ProgressCache) Set(ctx context.Context, executionID uint64, snapshot ScanProgressSnapshot) error {
	if c == nil || executionID == 0 {
		return nil
	}
	snapshot.ExecutionID = executionID
	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now().UTC()
	}
	if c.client == nil {
		c.mu.Lock()
		c.mem[executionID] = snapshot
		c.mu.Unlock()
		return nil
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("compliance:scan-progress:%d", executionID), payload, c.ttl).Err()
}

func (c *ProgressCache) Get(ctx context.Context, executionID uint64) (ScanProgressSnapshot, error) {
	if c == nil || executionID == 0 {
		return ScanProgressSnapshot{}, nil
	}
	if c.client == nil {
		c.mu.RLock()
		item := c.mem[executionID]
		c.mu.RUnlock()
		return item, nil
	}
	raw, err := c.client.Get(ctx, fmt.Sprintf("compliance:scan-progress:%d", executionID)).Bytes()
	if err != nil {
		return ScanProgressSnapshot{}, nil
	}
	var item ScanProgressSnapshot
	if err := json.Unmarshal(raw, &item); err != nil {
		return ScanProgressSnapshot{}, err
	}
	return item, nil
}
