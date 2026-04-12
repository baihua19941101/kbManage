package logs

import (
	"context"
	"time"
)

type Provider interface {
	Query(ctx context.Context, req QueryRequest) (QueryResult, error)
}

type QueryRequest struct {
	ClusterIDs []uint64
	Namespace  string
	Workload   string
	Pod        string
	Container  string
	Keyword    string
	StartAt    string
	EndAt      string
	Limit      int
}

type QueryResult struct {
	Items []LogEntry
}

type LogEntry struct {
	Timestamp string
	Message   string
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Query(_ context.Context, req QueryRequest) (QueryResult, error) {
	now := time.Now().UTC()
	return QueryResult{
		Items: []LogEntry{
			{
				Timestamp: now.Add(-2 * time.Minute).Format(time.RFC3339),
				Message:   "probe succeeded",
			},
			{
				Timestamp: now.Add(-1 * time.Minute).Format(time.RFC3339),
				Message:   "latency within threshold",
			},
		},
	}, nil
}
