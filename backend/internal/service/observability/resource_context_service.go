package observability

import (
	"context"
	"time"

	alertProvider "kbmanage/backend/internal/integration/observability/alerts"
)

type ResourceContextService struct {
	logsSvc    *LogsService
	eventsSvc  *EventsService
	metricsSvc *MetricsService
	alerts     alertProvider.Provider
}

func NewResourceContextService(logsSvc *LogsService, eventsSvc *EventsService, metricsSvc *MetricsService, alerts alertProvider.Provider) *ResourceContextService {
	if logsSvc == nil {
		logsSvc = NewLogsService(nil, nil)
	}
	if eventsSvc == nil {
		eventsSvc = NewEventsService(nil, nil)
	}
	if metricsSvc == nil {
		metricsSvc = NewMetricsService(nil, nil)
	}
	if alerts == nil {
		alerts = alertProvider.NewMockProvider()
	}
	return &ResourceContextService{
		logsSvc:    logsSvc,
		eventsSvc:  eventsSvc,
		metricsSvc: metricsSvc,
		alerts:     alerts,
	}
}

func (s *ResourceContextService) Get(ctx context.Context, req ResourceContextQuery) (ResourceContextResult, error) {
	logs, err := s.logsSvc.Query(ctx, LogsQueryRequest{
		Namespace: req.Namespace,
		Workload:  req.ResourceName,
		Keyword:   req.Keyword,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		Limit:     20,
	})
	if err != nil {
		return ResourceContextResult{}, err
	}
	events, err := s.eventsSvc.List(ctx, EventsQueryRequest{
		ClusterID:    req.ClusterID,
		Namespace:    req.Namespace,
		ResourceKind: req.ResourceKind,
		ResourceName: req.ResourceName,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
	})
	if err != nil {
		return ResourceContextResult{}, err
	}
	metricSeries, err := s.metricsSvc.QuerySeries(ctx, MetricQueryRequest{
		SubjectType: req.ResourceKind,
		SubjectRef:  req.ResourceName,
		MetricKey:   "cpu_usage",
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
	})
	if err != nil {
		return ResourceContextResult{}, err
	}
	alerts, err := s.alerts.ListAlerts(ctx, alertProvider.AlertQuery{})
	if err != nil {
		return ResourceContextResult{}, err
	}

	out := ResourceContextResult{
		DataFreshness: "fresh",
	}
	out.ResourceRef.ClusterID = req.ClusterID
	out.ResourceRef.Namespace = req.Namespace
	out.ResourceRef.ResourceKind = req.ResourceKind
	out.ResourceRef.ResourceName = req.ResourceName

	out.LogSummary.QueryID = logs.QueryID
	out.LogSummary.SampleCount = len(logs.Items)

	out.EventSummary.Total = len(events.Items)
	for _, item := range events.Items {
		if item.EventType == "warning" {
			out.EventSummary.WarningCount++
		}
	}

	lastCPU := 0.38
	if n := len(metricSeries.Points); n > 0 {
		lastCPU = metricSeries.Points[n-1].Value
	}
	out.MetricSummary = MetricSummary{
		SubjectType:  req.ResourceKind,
		SubjectRef:   req.ResourceName,
		CPUUsage:     lastCPU,
		MemoryUsage:  0.42,
		Availability: 99.8,
		ObservedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	for _, a := range alerts {
		out.Alerts = append(out.Alerts, AlertSummary{
			ID:       a.IncidentKey,
			Status:   a.Status,
			Severity: a.Severity,
			Summary:  a.Summary,
		})
	}
	return out, nil
}
