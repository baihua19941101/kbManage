package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type ArchiveExportCacheSnapshot struct {
	ExportID    string    `json:"exportId"`
	Status      string    `json:"status"`
	DownloadURL string    `json:"downloadUrl,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ExportCache struct {
	client *redis.Client
	ttl    time.Duration
	mu     sync.RWMutex
	mem    map[string]ArchiveExportCacheSnapshot
}

func NewExportCache(client *redis.Client, ttl time.Duration) *ExportCache {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &ExportCache{client: client, ttl: ttl, mem: make(map[string]ArchiveExportCacheSnapshot)}
}

func (c *ExportCache) Set(ctx context.Context, exportID string, snapshot ArchiveExportCacheSnapshot) error {
	if c == nil || exportID == "" {
		return nil
	}
	snapshot.ExportID = exportID
	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now().UTC()
	}
	if c.client == nil {
		c.mu.Lock()
		c.mem[exportID] = snapshot
		c.mu.Unlock()
		return nil
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("compliance:archive-export:%s", exportID), payload, c.ttl).Err()
}
