package gitops

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const defaultGitOpsOperationQueueKey = "gitops:operation:queue"

// OperationQueue abstracts delivery operation queueing across memory/redis.
type OperationQueue interface {
	Enqueue(ctx context.Context, operationID uint64) error
	Dequeue(ctx context.Context) (uint64, error)
}

func NewOperationQueue(rdb *redis.Client) OperationQueue {
	if rdb != nil {
		return &redisOperationQueue{rdb: rdb, queueKey: defaultGitOpsOperationQueueKey}
	}
	return &memoryOperationQueue{queue: make(chan uint64, 1024)}
}

type memoryOperationQueue struct {
	queue chan uint64
}

func (q *memoryOperationQueue) Enqueue(ctx context.Context, operationID uint64) error {
	if q == nil || q.queue == nil {
		return nil
	}
	select {
	case q.queue <- operationID:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *memoryOperationQueue) Dequeue(ctx context.Context) (uint64, error) {
	if q == nil || q.queue == nil {
		return 0, fmt.Errorf("queue is not initialized")
	}
	select {
	case id := <-q.queue:
		return id, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

type redisOperationQueue struct {
	rdb      *redis.Client
	queueKey string
}

func (q *redisOperationQueue) Enqueue(ctx context.Context, operationID uint64) error {
	if q == nil || q.rdb == nil {
		return nil
	}
	return q.rdb.RPush(ctx, q.queueKey, operationID).Err()
}

func (q *redisOperationQueue) Dequeue(ctx context.Context) (uint64, error) {
	if q == nil || q.rdb == nil {
		return 0, fmt.Errorf("queue is not configured")
	}
	res, err := q.rdb.BLPop(ctx, 0, q.queueKey).Result()
	if err != nil {
		return 0, err
	}
	if len(res) < 2 {
		return 0, fmt.Errorf("invalid queue payload")
	}
	id, err := strconv.ParseUint(res[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid operation id in queue: %w", err)
	}
	return id, nil
}
