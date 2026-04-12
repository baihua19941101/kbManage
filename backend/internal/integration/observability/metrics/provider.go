package metrics

import (
	"context"
	"time"
)

type Provider interface {
	QuerySeries(ctx context.Context, req SeriesQuery) (SeriesResult, error)
}

type SeriesQuery struct {
	ClusterIDs []uint64
	Subject    string
	MetricKey  string
	StartAt    string
	EndAt      string
	Step       string
}

type SeriesResult struct {
	Points []Point
}

type Point struct {
	Timestamp string
	Value     float64
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) QuerySeries(_ context.Context, req SeriesQuery) (SeriesResult, error) {
	now := time.Now().UTC()
	return SeriesResult{
		Points: []Point{
			{Timestamp: now.Add(-10 * time.Minute).Format(time.RFC3339), Value: 0.32},
			{Timestamp: now.Add(-5 * time.Minute).Format(time.RFC3339), Value: 0.41},
			{Timestamp: now.Format(time.RFC3339), Value: 0.38},
		},
	}, nil
}
