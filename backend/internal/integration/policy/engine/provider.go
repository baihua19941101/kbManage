package engine

import (
	"context"

	"kbmanage/backend/internal/domain"
)

// EvaluationProvider defines policy rule evaluation abstraction.
type EvaluationProvider interface {
	Evaluate(ctx context.Context, policy domain.SecurityPolicy, object map[string]any) (EvaluationResult, error)
}

type EvaluationResult struct {
	Hit        bool
	HitResult  domain.PolicyHitResult
	Message    string
	RiskLevel  domain.PolicyRiskLevel
	OccurredAt int64
}

// NoopProvider is a safe default for bootstrap/testing paths.
type NoopProvider struct{}

func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

func (p *NoopProvider) Evaluate(_ context.Context, _ domain.SecurityPolicy, _ map[string]any) (EvaluationResult, error) {
	return EvaluationResult{Hit: false, HitResult: domain.PolicyHitResultPass}, nil
}
