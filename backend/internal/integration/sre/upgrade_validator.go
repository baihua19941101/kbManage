package sre

import (
	"context"
	"fmt"
	"strings"
)

type UpgradePrecheckInput struct {
	CurrentVersion string
	TargetVersion  string
	WorkspaceID    uint64
	ProjectID      uint64
}

type UpgradePrecheckResult struct {
	Decision             string   `json:"decision"`
	CompatibilitySummary string   `json:"compatibilitySummary"`
	Blockers             []string `json:"blockers"`
	Warnings             []string `json:"warnings"`
	AllowConditions      []string `json:"allowConditions"`
}

type UpgradeValidator interface {
	Validate(ctx context.Context, input UpgradePrecheckInput) UpgradePrecheckResult
}

type StaticUpgradeValidator struct{}

func NewStaticUpgradeValidator() UpgradeValidator { return StaticUpgradeValidator{} }

func (StaticUpgradeValidator) Validate(_ context.Context, input UpgradePrecheckInput) UpgradePrecheckResult {
	current := strings.TrimSpace(input.CurrentVersion)
	target := strings.TrimSpace(input.TargetVersion)
	result := UpgradePrecheckResult{
		Decision:             "allow",
		CompatibilitySummary: fmt.Sprintf("平台版本可从 %s 升级到 %s", current, target),
		AllowConditions: []string{
			"升级前确认维护窗口已开启或业务低峰已确认",
			"升级前确认容量基线未处于 critical",
		},
	}
	if current == "" || target == "" || current == target {
		result.Decision = "block"
		result.Blockers = append(result.Blockers, "当前版本与目标版本必须同时提供且不能相同")
	}
	if strings.Contains(target, "unsupported") {
		result.Decision = "block"
		result.Blockers = append(result.Blockers, "目标版本不在受支持升级路径中")
	}
	if strings.Contains(target, "rc") || strings.Contains(target, "beta") {
		result.Warnings = append(result.Warnings, "目标版本为预发布版本，需要额外验证")
	}
	return result
}
