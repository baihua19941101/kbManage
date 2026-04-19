package identity

import (
	"context"
	"strings"
	"time"
)

type HealthCheckRequest struct {
	Name       string
	SourceType string
}

type HealthCheckResult struct {
	Available     bool       `json:"available"`
	Status        string     `json:"status"`
	LastCheckedAt *time.Time `json:"lastCheckedAt,omitempty"`
	Message       string     `json:"message,omitempty"`
}

type Provider interface {
	CheckHealth(context.Context, HealthCheckRequest) (HealthCheckResult, error)
}

type StaticProvider struct{}

func NewStaticProvider() Provider {
	return &StaticProvider{}
}

func (p *StaticProvider) CheckHealth(_ context.Context, req HealthCheckRequest) (HealthCheckResult, error) {
	now := time.Now()
	return HealthCheckResult{
		Available:     true,
		Status:        "active",
		LastCheckedAt: &now,
		Message:       strings.TrimSpace(req.SourceType) + " identity source is reachable",
	}, nil
}
