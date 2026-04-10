package operation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const defaultOperationQueueKey = "operation:queue"

// QueueService defines queue primitives independent from storage backend.
type QueueService interface {
	Enqueue(ctx context.Context, operationID uint64) error
	Dequeue(ctx context.Context) (uint64, error)
}

func NewQueueService(rdb *redis.Client) QueueService {
	if rdb != nil {
		return &redisQueueService{rdb: rdb, queueKey: defaultOperationQueueKey}
	}
	return &memoryQueueService{queue: make(chan uint64, 1024)}
}

type memoryQueueService struct {
	queue chan uint64
}

func (q *memoryQueueService) Enqueue(ctx context.Context, operationID uint64) error {
	select {
	case q.queue <- operationID:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *memoryQueueService) Dequeue(ctx context.Context) (uint64, error) {
	select {
	case id := <-q.queue:
		return id, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

type redisQueueService struct {
	rdb      *redis.Client
	queueKey string
}

func (q *redisQueueService) Enqueue(ctx context.Context, operationID uint64) error {
	return q.rdb.RPush(ctx, q.queueKey, operationID).Err()
}

func (q *redisQueueService) Dequeue(ctx context.Context) (uint64, error) {
	result, err := q.rdb.BLPop(ctx, 0, q.queueKey).Result()
	if err != nil {
		return 0, err
	}
	if len(result) < 2 {
		return 0, fmt.Errorf("invalid queue payload")
	}
	id, err := strconv.ParseUint(result[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid operation id in queue: %w", err)
	}
	return id, nil
}
