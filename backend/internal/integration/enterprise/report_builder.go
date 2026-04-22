package enterprise

import "context"

type ReportBuildInput struct {
	ReportType   string
	AudienceType string
	TimeRange    string
}

type ReportBuildOutput struct {
	SummarySection    string
	DetailSection     string
	AttachmentCatalog []string
}

type ReportBuilder interface {
	Build(ctx context.Context, input ReportBuildInput) ReportBuildOutput
}

type StaticReportBuilder struct{}

func NewStaticReportBuilder() *StaticReportBuilder { return &StaticReportBuilder{} }

func (b *StaticReportBuilder) Build(_ context.Context, input ReportBuildInput) ReportBuildOutput {
	return ReportBuildOutput{
		SummarySection:    "报表类型：" + input.ReportType + "；对象：" + input.AudienceType,
		DetailSection:     "时间范围：" + input.TimeRange,
		AttachmentCatalog: []string{"summary.pdf", "details.csv"},
	}
}
