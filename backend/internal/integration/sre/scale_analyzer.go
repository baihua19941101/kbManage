package sre

import (
	"context"
	"fmt"
	"strings"
)

type ScaleEvidenceInput struct {
	EvidenceType    string
	Summary         string
	ForecastSummary string
	ConfidenceLevel string
}

type ScaleAnalysis struct {
	Summary         string `json:"summary"`
	Bottleneck      string `json:"bottleneck"`
	ForecastSummary string `json:"forecastSummary"`
	ConfidenceLevel string `json:"confidenceLevel"`
}

type ScaleAnalyzer interface {
	Analyze(ctx context.Context, input ScaleEvidenceInput) ScaleAnalysis
}

type StaticScaleAnalyzer struct{}

func NewStaticScaleAnalyzer() ScaleAnalyzer { return StaticScaleAnalyzer{} }

func (StaticScaleAnalyzer) Analyze(_ context.Context, input ScaleEvidenceInput) ScaleAnalysis {
	confidence := strings.TrimSpace(input.ConfidenceLevel)
	if confidence == "" {
		confidence = "medium"
	}
	summary := strings.TrimSpace(input.Summary)
	if summary == "" {
		summary = "已记录规模化治理证据"
	}
	forecast := strings.TrimSpace(input.ForecastSummary)
	if forecast == "" {
		forecast = "未来 7 天内暂无明显容量透支风险"
	}
	return ScaleAnalysis{
		Summary:         summary,
		Bottleneck:      fmt.Sprintf("%s 场景下关注 API 并发与任务积压", strings.TrimSpace(input.EvidenceType)),
		ForecastSummary: forecast,
		ConfidenceLevel: confidence,
	}
}
