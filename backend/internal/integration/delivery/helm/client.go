package helm

import (
	"context"
	"errors"
	"strings"
)

var ErrHelmClientNotConfigured = errors.New("helm delivery client is not configured")

type ChartVersion struct {
	Version    string `json:"version"`
	AppVersion string `json:"appVersion"`
	CreatedAt  string `json:"createdAt,omitempty"`
	Deprecated bool   `json:"deprecated"`
}

type ListChartVersionsRequest struct {
	RepositoryURL string
	ChartName     string
	CredentialRef string
}

type RenderRequest struct {
	RepositoryURL string
	ChartName     string
	ChartVersion  string
	Namespace     string
	ValuesYAML    string
	CredentialRef string
}

type RenderResult struct {
	Manifest string `json:"manifest"`
	Digest   string `json:"digest"`
}

type Client interface {
	ListChartVersions(ctx context.Context, req ListChartVersionsRequest) ([]ChartVersion, error)
	Render(ctx context.Context, req RenderRequest) (*RenderResult, error)
}

type NoopClient struct{}

func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

func (c *NoopClient) ListChartVersions(ctx context.Context, req ListChartVersionsRequest) ([]ChartVersion, error) {
	_ = ctx
	if strings.TrimSpace(req.RepositoryURL) == "" || strings.TrimSpace(req.ChartName) == "" {
		return nil, errors.New("helm repositoryURL and chartName are required")
	}
	return []ChartVersion{}, ErrHelmClientNotConfigured
}

func (c *NoopClient) Render(ctx context.Context, req RenderRequest) (*RenderResult, error) {
	_ = ctx
	if strings.TrimSpace(req.RepositoryURL) == "" || strings.TrimSpace(req.ChartName) == "" {
		return nil, errors.New("helm repositoryURL and chartName are required")
	}
	if strings.TrimSpace(req.ChartVersion) == "" {
		return nil, errors.New("helm chartVersion is required")
	}
	return nil, ErrHelmClientNotConfigured
}
