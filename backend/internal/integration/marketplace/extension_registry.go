package marketplace

import (
	"context"
	"encoding/json"
	"strings"

	"kbmanage/backend/internal/domain"
)

type CompatibilitySeed struct {
	TargetType string `json:"targetType"`
	TargetRef  string `json:"targetRef"`
	Result     string `json:"result"`
	Summary    string `json:"summary"`
}

type ExtensionRegistry interface {
	NormalizePermissionDeclaration(ctx context.Context, declaration string) (string, error)
	NormalizeCompatibilityPolicy(ctx context.Context, policy string) (string, error)
}

type StaticExtensionRegistry struct{}

func NewStaticExtensionRegistry() *StaticExtensionRegistry {
	return &StaticExtensionRegistry{}
}

func (r *StaticExtensionRegistry) NormalizePermissionDeclaration(_ context.Context, declaration string) (string, error) {
	return normalizeJSONText(declaration)
}

func (r *StaticExtensionRegistry) NormalizeCompatibilityPolicy(_ context.Context, policy string) (string, error) {
	return normalizeJSONText(policy)
}

func normalizeJSONText(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil
	}
	var payload any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return "", err
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func CompatibilityOwnerRefForExtension(item *domain.ExtensionPackage) string {
	if item == nil {
		return ""
	}
	return strings.TrimSpace(item.Name) + ":" + strings.TrimSpace(item.Version)
}
