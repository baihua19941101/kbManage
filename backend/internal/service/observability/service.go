package observability

import (
	"context"
	"errors"
	"strconv"
	"time"

	alertProvider "kbmanage/backend/internal/integration/observability/alerts"
	kubeadapter "kbmanage/backend/internal/kube/adapter"
)

var (
	ErrObservabilityScopeDenied = errors.New("observability scope access denied")
	ErrObservabilityUnavailable = errors.New("observability service is not fully configured")
)

type Service struct {
	scope              *ScopeService
	alerts             alertProvider.Provider
	datasourceService  *DataSourceService
	metricsService     *MetricsService
	logsService        *LogsService
	eventsService      *EventsService
	resourceCtxService *ResourceContextService
}

func NewService(scope *ScopeService) *Service {
	if scope == nil {
		scope = NewScopeService(nil, nil, nil)
	}
	alerts := alertProvider.NewMockProvider()
	logsSvc := NewLogsService(nil, scope)
	eventsSvc := NewEventsService(kubeadapter.NewKubeEventReader(nil), scope)
	metricsSvc := NewMetricsService(nil, scope)
	return &Service{
		scope:             scope,
		alerts:            alerts,
		datasourceService: NewDataSourceService(),
		metricsService:    metricsSvc,
		logsService:       logsSvc,
		eventsService:     eventsSvc,
		resourceCtxService: NewResourceContextService(
			logsSvc,
			eventsSvc,
			metricsSvc,
			alerts,
		),
	}
}

type OverviewRequest struct {
	ClusterIDs []uint64
	StartAt    string
	EndAt      string
}

type OverviewResult struct {
	Cards            []OverviewCard  `json:"cards"`
	HotAlerts        []AlertSummary  `json:"hotAlerts"`
	TopEvents        []EventTimeline `json:"topEvents"`
	MetricHighlights []MetricSummary `json:"metricHighlights"`
}

type OverviewCard struct {
	Title    string  `json:"title"`
	Value    float64 `json:"value"`
	Unit     string  `json:"unit,omitempty"`
	Trend    string  `json:"trend,omitempty"`
	Severity string  `json:"severity,omitempty"`
}

type AlertSummary struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
}

type EventTimeline = kubeadapter.EventItem

type MetricSummary struct {
	SubjectType  string   `json:"subjectType"`
	SubjectRef   string   `json:"subjectRef"`
	CPUUsage     float64  `json:"cpuUsage"`
	MemoryUsage  float64  `json:"memoryUsage"`
	Availability float64  `json:"availability"`
	AnomalyFlags []string `json:"anomalyFlags"`
	ObservedAt   string   `json:"observedAt"`
}

type DataSourceConfig struct {
	Name           string  `json:"name"`
	ProviderKind   string  `json:"providerKind"`
	BaseURL        string  `json:"baseUrl"`
	Status         string  `json:"status"`
	LastVerifiedAt *string `json:"lastVerifiedAt,omitempty"`
	LastError      string  `json:"lastError,omitempty"`
}

type ClusterObservabilityConfig struct {
	ClusterID     uint64           `json:"clusterId"`
	Metrics       DataSourceConfig `json:"metrics"`
	Logs          DataSourceConfig `json:"logs"`
	Alerts        DataSourceConfig `json:"alerts"`
	EventsEnabled bool             `json:"eventsEnabled"`
}

type ScopeFilterRequest struct {
	ClusterIDs   []uint64
	Namespace    string
	ResourceKind string
	ResourceName string
}

type LogQueryResult struct {
	QueryID       string             `json:"queryId"`
	Status        string             `json:"status"`
	Items         []LogEntryResponse `json:"items"`
	NextCursor    string             `json:"nextCursor,omitempty"`
	DataFreshness string             `json:"dataFreshness"`
}

type LogEntryResponse struct {
	Timestamp   string `json:"timestamp"`
	ClusterID   string `json:"clusterId,omitempty"`
	WorkspaceID string `json:"workspaceId,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	Workload    string `json:"workload,omitempty"`
	Pod         string `json:"pod,omitempty"`
	Container   string `json:"container,omitempty"`
	Message     string `json:"message"`
}

type EventListResult struct {
	Items []kubeadapter.EventItem `json:"items"`
}

type MetricSeriesResult struct {
	MetricKey     string       `json:"metricKey"`
	SubjectType   string       `json:"subjectType"`
	SubjectRef    string       `json:"subjectRef"`
	Window        MetricWindow `json:"window"`
	Points        []PointValue `json:"points"`
	DataFreshness string       `json:"dataFreshness"`
}

type MetricWindow struct {
	StartAt string `json:"startAt,omitempty"`
	EndAt   string `json:"endAt,omitempty"`
}

type PointValue struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

type ResourceContextResult struct {
	ResourceRef struct {
		ClusterID    string `json:"clusterId"`
		Namespace    string `json:"namespace"`
		ResourceKind string `json:"resourceKind"`
		ResourceName string `json:"resourceName"`
	} `json:"resourceRef"`
	LogSummary struct {
		QueryID     string `json:"queryId"`
		SampleCount int    `json:"sampleCount"`
	} `json:"logSummary"`
	EventSummary struct {
		Total        int `json:"total"`
		WarningCount int `json:"warningCount"`
	} `json:"eventSummary"`
	MetricSummary MetricSummary  `json:"metricSummary"`
	Alerts        []AlertSummary `json:"alerts"`
	DataFreshness string         `json:"dataFreshness"`
}

func (s *Service) Overview(ctx context.Context, req OverviewRequest) (OverviewResult, error) {
	if s == nil || s.scope == nil {
		return OverviewResult{}, ErrObservabilityUnavailable
	}

	clusterID := "default-cluster"
	if len(req.ClusterIDs) > 0 {
		clusterID = uint64ToString(req.ClusterIDs[0])
	}
	alerts, err := s.alerts.ListAlerts(ctx, alertProvider.AlertQuery{})
	if err != nil {
		return OverviewResult{}, err
	}
	events, err := s.eventsService.List(ctx, EventsQueryRequest{
		ClusterID:    clusterID,
		Namespace:    "default",
		ResourceKind: "Deployment",
		ResourceName: "mock-app",
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
	})
	if err != nil {
		return OverviewResult{}, err
	}

	now := time.Now().UTC()
	res := OverviewResult{
		Cards: []OverviewCard{
			{Title: "Healthy Clusters", Value: 1, Unit: "clusters", Trend: "stable"},
			{Title: "Active Alerts", Value: float64(len(alerts)), Unit: "alerts", Severity: "warning"},
		},
		MetricHighlights: []MetricSummary{
			{
				SubjectType:  "cluster",
				SubjectRef:   clusterID,
				CPUUsage:     0.38,
				MemoryUsage:  0.44,
				Availability: 99.9,
				AnomalyFlags: []string{},
				ObservedAt:   now.Format(time.RFC3339),
			},
		},
	}
	for _, a := range alerts {
		res.HotAlerts = append(res.HotAlerts, AlertSummary{
			ID: a.IncidentKey, Status: a.Status, Severity: a.Severity, Summary: a.Summary,
		})
	}
	res.TopEvents = append(res.TopEvents, events.Items...)
	return res, nil
}

func (s *Service) AuthorizeScope(
	ctx context.Context,
	userID uint64,
	filter ScopeFilter,
) (ScopeFilter, error) {
	if s == nil || s.scope == nil {
		return ScopeFilter{}, ErrObservabilityUnavailable
	}
	return s.scope.FilterByScope(ctx, userID, filter)
}

func (s *Service) GetClusterConfig(clusterID uint64) ClusterObservabilityConfig {
	if s == nil || s.datasourceService == nil {
		return ClusterObservabilityConfig{}
	}
	return s.datasourceService.GetClusterConfig(clusterID)
}

func (s *Service) PutClusterConfig(cfg ClusterObservabilityConfig) ClusterObservabilityConfig {
	if s == nil || s.datasourceService == nil {
		return ClusterObservabilityConfig{}
	}
	res, err := s.datasourceService.UpdateClusterConfig(cfg.ClusterID, UpdateObservabilityConfigRequest{
		Metrics: &DataSourceConfigInput{
			Name:         cfg.Metrics.Name,
			ProviderKind: cfg.Metrics.ProviderKind,
			BaseURL:      cfg.Metrics.BaseURL,
		},
		Logs: &DataSourceConfigInput{
			Name:         cfg.Logs.Name,
			ProviderKind: cfg.Logs.ProviderKind,
			BaseURL:      cfg.Logs.BaseURL,
		},
		Alerts: &DataSourceConfigInput{
			Name:         cfg.Alerts.Name,
			ProviderKind: cfg.Alerts.ProviderKind,
			BaseURL:      cfg.Alerts.BaseURL,
		},
		EventsEnabled: &cfg.EventsEnabled,
	})
	if err != nil {
		return cfg
	}
	return res
}

func uint64ToString(in uint64) string {
	return strconv.FormatUint(in, 10)
}

func (s *Service) UpdateClusterConfig(clusterID uint64, req UpdateObservabilityConfigRequest) (ClusterObservabilityConfig, error) {
	if s == nil || s.datasourceService == nil {
		return ClusterObservabilityConfig{}, ErrObservabilityUnavailable
	}
	return s.datasourceService.UpdateClusterConfig(clusterID, req)
}

func (s *Service) QueryLogs(ctx context.Context, req LogsQueryRequest) (LogQueryResult, error) {
	if s == nil || s.logsService == nil {
		return LogQueryResult{}, ErrObservabilityUnavailable
	}
	return s.logsService.Query(ctx, req)
}

func (s *Service) ListEvents(ctx context.Context, req EventsQueryRequest) (EventListResult, error) {
	if s == nil || s.eventsService == nil {
		return EventListResult{}, ErrObservabilityUnavailable
	}
	return s.eventsService.List(ctx, req)
}

func (s *Service) QueryMetricSeries(ctx context.Context, req MetricQueryRequest) (MetricSeriesResult, error) {
	if s == nil || s.metricsService == nil {
		return MetricSeriesResult{}, ErrObservabilityUnavailable
	}
	result, err := s.metricsService.QuerySeries(ctx, req)
	if err != nil {
		return MetricSeriesResult{}, err
	}

	return result, nil
}

func (s *Service) ResourceContext(ctx context.Context, req ResourceContextQuery) (ResourceContextResult, error) {
	if s == nil || s.resourceCtxService == nil {
		return ResourceContextResult{}, ErrObservabilityUnavailable
	}
	return s.resourceCtxService.Get(ctx, req)
}
