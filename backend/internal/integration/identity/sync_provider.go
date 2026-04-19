package identity

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type SyncedPrincipal struct {
	ExternalRef   string `json:"externalRef"`
	PrincipalType string `json:"principalType"`
}

type SyncRequest struct {
	SourceID   uint64
	SourceType string
}

type SyncResult struct {
	State          string            `json:"state"`
	CompletedAt    *time.Time        `json:"completedAt,omitempty"`
	Principals     []SyncedPrincipal `json:"principals,omitempty"`
	Message        string            `json:"message,omitempty"`
	NormalizedCode string            `json:"normalizedCode,omitempty"`
}

type SyncProvider interface {
	SyncDirectory(context.Context, SyncRequest) (SyncResult, error)
}

type StaticSyncProvider struct{}

func NewStaticSyncProvider() SyncProvider {
	return &StaticSyncProvider{}
}

func (p *StaticSyncProvider) SyncDirectory(_ context.Context, req SyncRequest) (SyncResult, error) {
	now := time.Now()
	return SyncResult{
		State:       "succeeded",
		CompletedAt: &now,
		Principals: []SyncedPrincipal{
			{
				ExternalRef:   fmt.Sprintf("%s-source-%d-admin", strings.TrimSpace(req.SourceType), req.SourceID),
				PrincipalType: "user",
			},
		},
		Message:        "directory sync completed",
		NormalizedCode: "identity_sync_ok",
	}, nil
}
