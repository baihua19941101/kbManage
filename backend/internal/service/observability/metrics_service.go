package observability

import (
	"context"

	metricProvider "kbmanage/backend/internal/integration/observability/metrics"
)

type MetricsService struct {
	provider metricProvider.Provider
	scope    *ScopeService
}

func NewMetricsService(provider metricProvider.Provider, scope *ScopeService) *MetricsService {
	if provider == nil {
		provider = metricProvider.NewPrometheusProvider("", nil)
	}
	return &MetricsService{provider: provider, scope: scope}
}

func (s *MetricsService) QuerySeries(ctx context.Context, req MetricQueryRequest) (MetricSeriesResult, error) {
	if s == nil || s.provider == nil {
		return MetricSeriesResult{}, ErrObservabilityUnavailable
	}
	if s.scope != nil {
		if access, ok := AccessContextFromContext(ctx); ok {
			filtered, err := s.scope.FilterByScope(ctx, access.UserID, ScopeFilter{ClusterIDs: req.ClusterIDs})
			if err != nil {
				return MetricSeriesResult{}, ErrObservabilityScopeDenied
			}
			req.ClusterIDs = filtered.ClusterIDs
			if access.ClusterConstrained && len(req.ClusterIDs) == 0 {
				return MetricSeriesResult{}, ErrObservabilityScopeDenied
			}
		}
	}

	series, err := s.provider.QuerySeries(ctx, metricProvider.SeriesQuery{
		ClusterIDs: req.ClusterIDs,
		Subject:    req.SubjectRef,
		MetricKey:  req.MetricKey,
		StartAt:    req.StartAt,
		EndAt:      req.EndAt,
		Step:       req.Step,
	})
	if err != nil {
		return MetricSeriesResult{}, normalizeObservabilityProviderError(err)
	}
	out := MetricSeriesResult{
		MetricKey:     req.MetricKey,
		SubjectType:   req.SubjectType,
		SubjectRef:    req.SubjectRef,
		Window:        MetricWindow{StartAt: req.StartAt, EndAt: req.EndAt},
		DataFreshness: "fresh",
	}
	for _, point := range series.Points {
		out.Points = append(out.Points, PointValue{
			Timestamp: point.Timestamp,
			Value:     point.Value,
		})
	}
	return out, nil
}
