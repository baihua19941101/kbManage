package diff

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type DiffType string

const (
	DiffTypeAdded       DiffType = "added"
	DiffTypeModified    DiffType = "modified"
	DiffTypeRemoved     DiffType = "removed"
	DiffTypeUnavailable DiffType = "unavailable"
)

type Summary struct {
	HasChanges bool     `json:"hasChanges"`
	Digest     string   `json:"digest"`
	DiffType   DiffType `json:"diffType,omitempty"`
}

type Comparator interface {
	Compare(ctx context.Context, desired []byte, live []byte) (*Summary, error)
}

type NoopComparator struct{}

func NewNoopComparator() *NoopComparator {
	return &NoopComparator{}
}

func (c *NoopComparator) Compare(ctx context.Context, desired []byte, live []byte) (*Summary, error) {
	_ = ctx
	digest := sha256.Sum256(append(append([]byte{}, desired...), live...))
	if bytes.Equal(desired, live) {
		return &Summary{
			HasChanges: false,
			Digest:     hex.EncodeToString(digest[:]),
		}, nil
	}
	diffType := DiffTypeModified
	switch {
	case len(desired) > 0 && len(live) == 0:
		diffType = DiffTypeAdded
	case len(desired) == 0 && len(live) > 0:
		diffType = DiffTypeRemoved
	}
	return &Summary{
		HasChanges: true,
		Digest:     hex.EncodeToString(digest[:]),
		DiffType:   diffType,
	}, nil
}
