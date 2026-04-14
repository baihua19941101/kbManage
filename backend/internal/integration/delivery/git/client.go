package git

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrGitClientNotConfigured = errors.New("git delivery client is not configured")

type Reference struct {
	Name      string     `json:"name"`
	Hash      string     `json:"hash"`
	IsTag     bool       `json:"isTag"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type ListRefsRequest struct {
	Endpoint      string
	CredentialRef string
	Prefix        string
}

type FetchRevisionRequest struct {
	Endpoint      string
	CredentialRef string
	Revision      string
	Path          string
}

type RevisionSnapshot struct {
	Revision  string            `json:"revision"`
	Files     map[string]string `json:"files"`
	FetchedAt time.Time         `json:"fetchedAt"`
}

type Client interface {
	ListRefs(ctx context.Context, req ListRefsRequest) ([]Reference, error)
	FetchRevision(ctx context.Context, req FetchRevisionRequest) (*RevisionSnapshot, error)
}

type NoopClient struct{}

func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

func (c *NoopClient) ListRefs(ctx context.Context, req ListRefsRequest) ([]Reference, error) {
	_ = ctx
	if strings.TrimSpace(req.Endpoint) == "" {
		return nil, errors.New("git endpoint is required")
	}
	return []Reference{}, ErrGitClientNotConfigured
}

func (c *NoopClient) FetchRevision(ctx context.Context, req FetchRevisionRequest) (*RevisionSnapshot, error) {
	_ = ctx
	if strings.TrimSpace(req.Endpoint) == "" {
		return nil, errors.New("git endpoint is required")
	}
	if strings.TrimSpace(req.Revision) == "" {
		return nil, errors.New("git revision is required")
	}
	return nil, ErrGitClientNotConfigured
}
