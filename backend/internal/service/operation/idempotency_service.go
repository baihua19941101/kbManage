package operation

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultIdempotencyTTL = 10 * time.Second

// IdempotencyService provides a minimal distributed/local idempotency lock.
type IdempotencyService struct {
	rdb       *redis.Client
	keyPrefix string
	ttl       time.Duration

	mu    sync.Mutex
	locks map[string]time.Time
}

func NewIdempotencyService(rdb *redis.Client) *IdempotencyService {
	return &IdempotencyService{
		rdb:       rdb,
		keyPrefix: "operation:idempotency:",
		ttl:       defaultIdempotencyTTL,
		locks:     make(map[string]time.Time),
	}
}

func (s *IdempotencyService) BuildRequestID(operatorID uint64, idempotencyKey string) string {
	cleanKey := strings.TrimSpace(idempotencyKey)
	if cleanKey == "" {
		return fmt.Sprintf("op:%d:%d", operatorID, time.Now().UnixNano())
	}
	return fmt.Sprintf("op:%d:%s", operatorID, cleanKey)
}

func (s *IdempotencyService) Acquire(ctx context.Context, requestID string) (bool, error) {
	if strings.TrimSpace(requestID) == "" {
		return true, nil
	}

	if s.rdb != nil {
		return s.rdb.SetNX(ctx, s.lockKey(requestID), "1", s.ttl).Result()
	}

	now := time.Now()
	expiresAt := now.Add(s.ttl)

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, expiry := range s.locks {
		if now.After(expiry) {
			delete(s.locks, key)
		}
	}
	if expiry, ok := s.locks[requestID]; ok && now.Before(expiry) {
		return false, nil
	}
	s.locks[requestID] = expiresAt
	return true, nil
}

func (s *IdempotencyService) Release(ctx context.Context, requestID string) error {
	if strings.TrimSpace(requestID) == "" {
		return nil
	}

	if s.rdb != nil {
		return s.rdb.Del(ctx, s.lockKey(requestID)).Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.locks, requestID)
	return nil
}

func (s *IdempotencyService) lockKey(requestID string) string {
	return s.keyPrefix + requestID
}
