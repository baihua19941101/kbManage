package observability

import (
	"context"
	"strconv"
	"strings"

	kubeadapter "kbmanage/backend/internal/kube/adapter"
)

type EventsService struct {
	reader kubeadapter.EventReader
	scope  *ScopeService
}

func NewEventsService(reader kubeadapter.EventReader, scope *ScopeService) *EventsService {
	if reader == nil {
		reader = kubeadapter.NewMockEventReader()
	}
	return &EventsService{reader: reader, scope: scope}
}

func (s *EventsService) List(ctx context.Context, req EventsQueryRequest) (EventListResult, error) {
	if s == nil || s.reader == nil {
		return EventListResult{}, ErrObservabilityUnavailable
	}
	if s.scope != nil {
		clusterIDs := parseEventClusterIDs(req.ClusterID)
		if access, ok := AccessContextFromContext(ctx); ok {
			filtered, err := s.scope.FilterByScope(ctx, access.UserID, ScopeFilter{ClusterIDs: clusterIDs})
			if err != nil {
				return EventListResult{}, ErrObservabilityScopeDenied
			}
			if access.ClusterConstrained && len(filtered.ClusterIDs) == 0 {
				return EventListResult{}, ErrObservabilityScopeDenied
			}
			if len(filtered.ClusterIDs) > 0 {
				req.ClusterID = strconv.FormatUint(filtered.ClusterIDs[0], 10)
			}
		}
	}

	items, err := s.reader.List(ctx, req.ClusterID, req.Namespace, req.ResourceKind, req.ResourceName)
	if err != nil {
		return EventListResult{}, normalizeObservabilityProviderError(err)
	}

	if req.EventType == "" {
		return EventListResult{Items: items}, nil
	}
	filtered := make([]kubeadapter.EventItem, 0, len(items))
	for _, item := range items {
		if strings.EqualFold(item.EventType, req.EventType) {
			filtered = append(filtered, item)
		}
	}
	return EventListResult{Items: filtered}, nil
}

func parseEventClusterIDs(raw string) []uint64 {
	out := make([]uint64, 0)
	for _, segment := range strings.Split(strings.TrimSpace(raw), ",") {
		item := strings.TrimSpace(segment)
		if item == "" {
			continue
		}
		id, err := strconv.ParseUint(item, 10, 64)
		if err != nil || id == 0 {
			continue
		}
		out = append(out, id)
	}
	return out
}
