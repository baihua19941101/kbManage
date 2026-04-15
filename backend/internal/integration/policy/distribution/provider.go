package distribution

import "context"

// DistributionProvider defines policy distribution abstraction for cluster targets.
type DistributionProvider interface {
	Apply(ctx context.Context, policyID uint64, assignmentID uint64) error
}

// NoopProvider is a safe default for bootstrap/testing paths.
type NoopProvider struct{}

func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

func (p *NoopProvider) Apply(_ context.Context, _ uint64, _ uint64) error {
	return nil
}
