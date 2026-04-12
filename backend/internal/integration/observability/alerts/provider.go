package alerts

import (
	"context"
	"time"
)

type Provider interface {
	ListAlerts(ctx context.Context, req AlertQuery) ([]AlertItem, error)
}

type AlertQuery struct {
	ClusterIDs []uint64
	StartAt    string
	EndAt      string
	Status     string
	Severity   string
}

type AlertItem struct {
	IncidentKey string
	Status      string
	Severity    string
	Summary     string
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) ListAlerts(_ context.Context, req AlertQuery) ([]AlertItem, error) {
	_ = req
	return []AlertItem{
		{
			IncidentKey: "inc-"+time.Now().UTC().Format("20060102150405"),
			Status:      "firing",
			Severity:    "warning",
			Summary:     "mock alert for observability center",
		},
	}, nil
}
