package observability

import (
	"context"
	"errors"
	"fmt"
	"time"

	logProvider "kbmanage/backend/internal/integration/observability/logs"
)

type LogsService struct {
	provider logProvider.Provider
	scope    *ScopeService
}

func NewLogsService(provider logProvider.Provider, scope *ScopeService) *LogsService {
	if provider == nil {
		provider = logProvider.NewLokiProvider("", nil)
	}
	return &LogsService{provider: provider, scope: scope}
}

func (s *LogsService) Query(ctx context.Context, req LogsQueryRequest) (LogQueryResult, error) {
	if s == nil || s.provider == nil {
		return LogQueryResult{}, ErrObservabilityUnavailable
	}
	if s.scope != nil {
		access, ok := AccessContextFromContext(ctx)
		if ok {
			filtered, err := s.scope.FilterByScope(ctx, access.UserID, ScopeFilter{ClusterIDs: req.ClusterIDs})
			if err != nil {
				return LogQueryResult{}, ErrObservabilityScopeDenied
			}
			req.ClusterIDs = filtered.ClusterIDs
			if access.ClusterConstrained && len(req.ClusterIDs) == 0 {
				return LogQueryResult{}, ErrObservabilityScopeDenied
			}
		}
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}
	items, err := s.provider.Query(ctx, logProvider.QueryRequest{
		ClusterIDs: req.ClusterIDs,
		Namespace:  req.Namespace,
		Workload:   req.Workload,
		Pod:        req.Pod,
		Container:  req.Container,
		Keyword:    req.Keyword,
		StartAt:    req.StartAt,
		EndAt:      req.EndAt,
		Limit:      req.Limit,
	})
	if err != nil {
		return LogQueryResult{}, normalizeObservabilityProviderError(err)
	}

	out := LogQueryResult{
		QueryID:       fmt.Sprintf("lq-%d", time.Now().UnixNano()),
		Status:        "ready",
		DataFreshness: "fresh",
	}
	if len(req.ClusterIDs) > 0 {
		outCluster := fmt.Sprintf("%d", req.ClusterIDs[0])
		for _, item := range items.Items {
			out.Items = append(out.Items, LogEntryResponse{
				Timestamp: item.Timestamp,
				ClusterID: outCluster,
				Namespace: req.Namespace,
				Workload:  req.Workload,
				Pod:       req.Pod,
				Container: req.Container,
				Message:   item.Message,
			})
		}
		return out, nil
	}

	for _, item := range items.Items {
		out.Items = append(out.Items, LogEntryResponse{
			Timestamp: item.Timestamp,
			Namespace: req.Namespace,
			Workload:  req.Workload,
			Pod:       req.Pod,
			Container: req.Container,
			Message:   item.Message,
		})
	}
	return out, nil
}

func normalizeObservabilityProviderError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrObservabilityScopeDenied) || errors.Is(err, ErrObservabilityUnavailable) {
		return err
	}
	return errors.New("observability upstream temporarily unavailable")
}
