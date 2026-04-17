package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProfileScheduleSnapshot struct {
	CronExpression string    `json:"cronExpression,omitempty"`
	ScheduleMode   string    `json:"scheduleMode,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ScheduleCache struct {
	client *redis.Client
	ttl    time.Duration
	mu     sync.RWMutex
	mem    map[uint64]ProfileScheduleSnapshot
}

func NewScheduleCache(client *redis.Client, ttl time.Duration) *ScheduleCache {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &ScheduleCache{client: client, ttl: ttl, mem: make(map[uint64]ProfileScheduleSnapshot)}
}

func (c *ScheduleCache) SetProfileSchedule(ctx context.Context, profileID uint64, snapshot ProfileScheduleSnapshot) error {
	if c == nil || profileID == 0 {
		return nil
	}
	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now().UTC()
	}
	if c.client == nil {
		c.mu.Lock()
		c.mem[profileID] = snapshot
		c.mu.Unlock()
		return nil
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("compliance:schedule:%d", profileID), payload, c.ttl).Err()
}
