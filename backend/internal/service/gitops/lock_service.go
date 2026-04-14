package gitops

import (
	"context"
	"time"

	"kbmanage/backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type LockService struct {
	client *redis.Client
}

func NewLockService(client *redis.Client) *LockService {
	return &LockService{client: client}
}

func (s *LockService) Acquire(ctx context.Context, scope string, token string, ttl time.Duration) (bool, error) {
	if s == nil || s.client == nil {
		return true, nil
	}
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return s.client.SetNX(ctx, lockKey(scope), token, ttl).Result()
}

func (s *LockService) Release(ctx context.Context, scope string, token string) error {
	if s == nil || s.client == nil {
		return nil
	}
	const lua = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`
	return s.client.Eval(ctx, lua, []string{lockKey(scope)}, token).Err()
}

func lockKey(scope string) string {
	return repository.GitOpsLockKey(scope)
}
