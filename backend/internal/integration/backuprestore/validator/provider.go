package validator

import (
	"context"
	"strings"
)

type Request struct {
	JobType           string
	ScopeSelection    map[string]any
	TargetEnvironment string
}

type Result struct {
	Status            string   `json:"status"`
	Blockers          []string `json:"blockers"`
	Warnings          []string `json:"warnings"`
	ConsistencyNotice string   `json:"consistencyNotice"`
}

type Provider interface {
	Validate(context.Context, Request) (Result, error)
}

type StaticProvider struct{}

func NewStaticProvider() Provider {
	return &StaticProvider{}
}

func (p *StaticProvider) Validate(_ context.Context, req Request) (Result, error) {
	result := Result{
		Status:            "passed",
		Blockers:          []string{},
		Warnings:          []string{},
		ConsistencyNotice: "恢复前校验通过，仍需在目标环境完成业务一致性核验",
	}
	if strings.TrimSpace(req.TargetEnvironment) == "" {
		result.Status = "blocked"
		result.Blockers = append(result.Blockers, "目标环境不能为空")
	}
	if len(req.ScopeSelection) == 0 {
		result.Status = "blocked"
		result.Blockers = append(result.Blockers, "至少选择一个恢复范围")
	}
	if req.JobType == "cross-cluster-restore" {
		result.Warnings = append(result.Warnings, "跨集群恢复前请确认网络、存储和密钥材料已预置")
	}
	if len(result.Blockers) > 0 {
		result.ConsistencyNotice = "存在阻断项，不能继续执行恢复或迁移"
	}
	return result, nil
}
