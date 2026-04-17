package baseline

import (
	"context"
	"encoding/json"

	"kbmanage/backend/internal/domain"
)

type Provider interface {
	BuildSnapshot(ctx context.Context, baseline *domain.ComplianceBaseline) (domain.ComplianceBaselineSnapshot, error)
}

type StaticProvider struct{}

func NewStaticProvider() Provider {
	return StaticProvider{}
}

func (StaticProvider) BuildSnapshot(_ context.Context, item *domain.ComplianceBaseline) (domain.ComplianceBaselineSnapshot, error) {
	if item == nil {
		return domain.ComplianceBaselineSnapshot{}, nil
	}
	var targetLevels []string
	_ = json.Unmarshal([]byte(item.TargetLevelsJSON), &targetLevels)
	var rules map[string]any
	_ = json.Unmarshal([]byte(item.RulesJSON), &rules)
	return domain.ComplianceBaselineSnapshot{
		BaselineID:   item.ID,
		Name:         item.Name,
		StandardType: item.StandardType,
		Version:      item.Version,
		VersionLabel: string(item.StandardType) + ":" + item.Version,
		RuleCount:    item.RuleCount,
		TargetLevels: targetLevels,
		Rules:        rules,
	}, nil
}
