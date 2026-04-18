package validator

import (
	"context"
	"fmt"
	"strings"
)

type CheckSeverity string

const (
	CheckSeverityInfo    CheckSeverity = "info"
	CheckSeverityWarning CheckSeverity = "warning"
	CheckSeverityBlocker CheckSeverity = "blocker"
)

type CheckResult struct {
	Code     string        `json:"code"`
	Message  string        `json:"message"`
	Severity CheckSeverity `json:"severity"`
	Passed   bool          `json:"passed"`
}

type Request struct {
	InfrastructureType string         `json:"infrastructureType"`
	DriverKey          string         `json:"driverKey"`
	DriverVersion      string         `json:"driverVersion"`
	RequiredDomains    []string       `json:"requiredDomains"`
	Parameters         map[string]any `json:"parameters"`
}

type Result struct {
	Status      string        `json:"status"`
	CanContinue bool          `json:"canContinue"`
	Summary     string        `json:"summary"`
	Checks      []CheckResult `json:"checks"`
}

type Provider interface {
	Validate(ctx context.Context, req Request) (Result, error)
}

type StaticProvider struct{}

func NewStaticProvider() Provider {
	return &StaticProvider{}
}

func (p *StaticProvider) Validate(_ context.Context, req Request) (Result, error) {
	if strings.TrimSpace(req.DriverKey) == "" {
		return Result{}, fmt.Errorf("driverKey is required")
	}
	checks := []CheckResult{
		{Code: "driver.available", Message: "驱动可用", Severity: CheckSeverityInfo, Passed: true},
	}
	status := "passed"
	canContinue := true
	if strings.TrimSpace(req.InfrastructureType) == "" {
		checks = append(checks, CheckResult{Code: "infra.missing", Message: "基础设施类型缺失", Severity: CheckSeverityBlocker, Passed: false})
		status = "failed"
		canContinue = false
	}
	for _, domain := range req.RequiredDomains {
		trimmed := strings.TrimSpace(domain)
		if trimmed == "" {
			continue
		}
		checks = append(checks, CheckResult{
			Code:     "capability." + trimmed,
			Message:  "能力域 " + trimmed + " 已声明",
			Severity: CheckSeverityInfo,
			Passed:   true,
		})
	}
	return Result{
		Status:      status,
		CanContinue: canContinue,
		Summary:     "validation completed",
		Checks:      checks,
	}, nil
}
