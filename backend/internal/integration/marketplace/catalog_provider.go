package marketplace

import (
	"context"
	"encoding/json"
	"strings"

	"kbmanage/backend/internal/domain"
)

type TemplateSeed struct {
	Name                string                `json:"name"`
	Slug                string                `json:"slug"`
	Category            string                `json:"category"`
	Summary             string                `json:"summary"`
	PublishStatus       string                `json:"publishStatus"`
	SupportedScopes     []string              `json:"supportedScopes"`
	ReleaseNotesSummary string                `json:"releaseNotesSummary"`
	Versions            []TemplateVersionSeed `json:"versions"`
}

type TemplateVersionSeed struct {
	Version                     string   `json:"version"`
	Status                      string   `json:"status"`
	Dependencies                []string `json:"dependencies"`
	ParameterSchemaSummary      string   `json:"parameterSchemaSummary"`
	DeploymentConstraintSummary string   `json:"deploymentConstraintSummary"`
	ReleaseNotes                string   `json:"releaseNotes"`
	IsUpgradeable               *bool    `json:"isUpgradeable"`
	SupersedesVersion           string   `json:"supersedesVersion"`
}

type CatalogSeedEnvelope struct {
	Templates []TemplateSeed `json:"templates"`
}

type CatalogSyncResult struct {
	Templates []TemplateSeed `json:"templates"`
}

type CatalogProvider interface {
	Sync(ctx context.Context, source domain.CatalogSource) (*CatalogSyncResult, error)
}

type StaticCatalogProvider struct{}

func NewStaticCatalogProvider() *StaticCatalogProvider {
	return &StaticCatalogProvider{}
}

func (p *StaticCatalogProvider) Sync(_ context.Context, source domain.CatalogSource) (*CatalogSyncResult, error) {
	trimmed := strings.TrimSpace(source.ConfigSummary)
	if trimmed == "" {
		return &CatalogSyncResult{Templates: []TemplateSeed{}}, nil
	}
	var envelope CatalogSeedEnvelope
	if err := json.Unmarshal([]byte(trimmed), &envelope); err != nil {
		return nil, err
	}
	if envelope.Templates == nil {
		envelope.Templates = []TemplateSeed{}
	}
	return &CatalogSyncResult{Templates: envelope.Templates}, nil
}
