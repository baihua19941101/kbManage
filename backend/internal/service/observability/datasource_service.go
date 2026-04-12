package observability

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type DataSourceConfigInput struct {
	Name          string `json:"name"`
	ProviderKind  string `json:"providerKind"`
	BaseURL       string `json:"baseUrl"`
	AuthSecretRef string `json:"authSecretRef,omitempty"`
}

type UpdateObservabilityConfigRequest struct {
	Metrics       *DataSourceConfigInput `json:"metrics"`
	Logs          *DataSourceConfigInput `json:"logs"`
	Alerts        *DataSourceConfigInput `json:"alerts"`
	EventsEnabled *bool                  `json:"eventsEnabled"`
}

type DataSourceService struct {
	mu      sync.RWMutex
	configs map[uint64]ClusterObservabilityConfig
}

func NewDataSourceService() *DataSourceService {
	return &DataSourceService{
		configs: make(map[uint64]ClusterObservabilityConfig),
	}
}

func (s *DataSourceService) GetClusterConfig(clusterID uint64) ClusterObservabilityConfig {
	s.mu.RLock()
	cfg, ok := s.configs[clusterID]
	s.mu.RUnlock()
	if ok {
		return cfg
	}

	now := time.Now().UTC().Format(time.RFC3339)
	return ClusterObservabilityConfig{
		ClusterID: clusterID,
		Metrics: DataSourceConfig{
			Name:           "default-metrics",
			ProviderKind:   "prometheus-compatible",
			BaseURL:        "http://prometheus.local",
			Status:         "healthy",
			LastVerifiedAt: &now,
		},
		Logs: DataSourceConfig{
			Name:           "default-logs",
			ProviderKind:   "loki-compatible",
			BaseURL:        "http://loki.local",
			Status:         "healthy",
			LastVerifiedAt: &now,
		},
		Alerts: DataSourceConfig{
			Name:           "default-alerts",
			ProviderKind:   "alertmanager-compatible",
			BaseURL:        "http://alertmanager.local",
			Status:         "healthy",
			LastVerifiedAt: &now,
		},
		EventsEnabled: true,
	}
}

func (s *DataSourceService) UpdateClusterConfig(clusterID uint64, req UpdateObservabilityConfigRequest) (ClusterObservabilityConfig, error) {
	current := s.GetClusterConfig(clusterID)
	current.ClusterID = clusterID

	var err error
	if req.Metrics != nil {
		current.Metrics, err = buildDataSourceConfig(*req.Metrics)
		if err != nil {
			return ClusterObservabilityConfig{}, fmt.Errorf("metrics: %w", err)
		}
	}
	if req.Logs != nil {
		current.Logs, err = buildDataSourceConfig(*req.Logs)
		if err != nil {
			return ClusterObservabilityConfig{}, fmt.Errorf("logs: %w", err)
		}
	}
	if req.Alerts != nil {
		current.Alerts, err = buildDataSourceConfig(*req.Alerts)
		if err != nil {
			return ClusterObservabilityConfig{}, fmt.Errorf("alerts: %w", err)
		}
	}
	if req.EventsEnabled != nil {
		current.EventsEnabled = *req.EventsEnabled
	}

	s.mu.Lock()
	s.configs[clusterID] = current
	s.mu.Unlock()
	return current, nil
}

func buildDataSourceConfig(input DataSourceConfigInput) (DataSourceConfig, error) {
	if strings.TrimSpace(input.Name) == "" {
		return DataSourceConfig{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(input.ProviderKind) == "" {
		return DataSourceConfig{}, fmt.Errorf("providerKind is required")
	}
	if strings.TrimSpace(input.BaseURL) == "" {
		return DataSourceConfig{}, fmt.Errorf("baseUrl is required")
	}

	status, lastError := "healthy", ""
	if !strings.HasPrefix(strings.ToLower(input.BaseURL), "http://") && !strings.HasPrefix(strings.ToLower(input.BaseURL), "https://") {
		status = "unhealthy"
		lastError = "baseUrl must start with http:// or https://"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	return DataSourceConfig{
		Name:           input.Name,
		ProviderKind:   input.ProviderKind,
		BaseURL:        input.BaseURL,
		Status:         status,
		LastVerifiedAt: &now,
		LastError:      lastError,
	}, nil
}
