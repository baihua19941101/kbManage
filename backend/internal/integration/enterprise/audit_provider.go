package enterprise

import "context"

type AuditSummaryInput struct {
	SubjectType string
	SubjectRef  string
	ChangeType  string
}

type AuditProvider interface {
	BuildTrailSummary(ctx context.Context, input AuditSummaryInput) string
}

type StaticAuditProvider struct{}

func NewStaticAuditProvider() *StaticAuditProvider { return &StaticAuditProvider{} }

func (p *StaticAuditProvider) BuildTrailSummary(_ context.Context, input AuditSummaryInput) string {
	return input.SubjectType + ":" + input.SubjectRef + ":" + input.ChangeType
}
