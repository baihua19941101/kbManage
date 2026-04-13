package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}

func ObservabilityCacheKey(parts ...string) string {
	if len(parts) == 0 {
		return "observability"
	}
	key := "observability"
	for _, part := range parts {
		if part == "" {
			continue
		}
		key += ":" + part
	}
	return key
}

func WorkloadOpsProgressKey(parts ...string) string {
	return prefixedRedisKey("workloadops:progress", parts...)
}

func WorkloadOpsSessionKey(parts ...string) string {
	return prefixedRedisKey("workloadops:session", parts...)
}

func WorkloadOpsBatchKey(parts ...string) string {
	return prefixedRedisKey("workloadops:batch", parts...)
}

func prefixedRedisKey(prefix string, parts ...string) string {
	if len(parts) == 0 {
		return prefix
	}
	key := prefix
	for _, part := range parts {
		if part == "" {
			continue
		}
		key += ":" + part
	}
	return key
}
