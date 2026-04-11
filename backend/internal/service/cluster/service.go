package cluster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/kube/adapter"
	kubeclient "kbmanage/backend/internal/kube/client"
	"kbmanage/backend/internal/repository"
)

type ConnectivityResult struct {
	ClusterID uint64 `json:"clusterId"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

type RegisterClusterRequest struct {
	Name              string `json:"name"`
	DisplayName       string `json:"displayName"`
	Description       string `json:"description"`
	Environment       string `json:"environment"`
	APIServer         string `json:"apiServer"`
	AuthType          string `json:"authType"`
	KubeConfig        string `json:"kubeConfig"`
	CredentialType    string `json:"credentialType"`
	CredentialPayload string `json:"credentialPayload"`
}

type HealthSummary struct {
	ClusterID uint64 `json:"clusterId"`
	Healthy   int64  `json:"healthy"`
	Degraded  int64  `json:"degraded"`
	Unknown   int64  `json:"unknown"`
	Total     int64  `json:"total"`
}

type SyncTriggerResult struct {
	ClusterID uint64 `json:"clusterId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type SyncDispatcher interface {
	Enqueue(clusterID uint64) error
}

type Service struct {
	clusterRepo     *repository.ClusterRepository
	credentialRepo  *repository.ClusterCredentialRepository
	resourceRepo    *repository.ResourceInventoryRepository
	resourceIndexer adapter.ResourceIndexer
	syncDispatcher  SyncDispatcher
	cipher          CredentialCipher
	clientManager   *kubeclient.Manager
}

func NewService(
	clusterRepo *repository.ClusterRepository,
	credentialRepo *repository.ClusterCredentialRepository,
	resourceRepo *repository.ResourceInventoryRepository,
	cipher CredentialCipher,
	resourceIndexer adapter.ResourceIndexer,
	syncDispatcher SyncDispatcher,
	clientManager *kubeclient.Manager,
) *Service {
	if resourceIndexer == nil {
		resourceIndexer = adapter.NoopResourceIndexer{}
	}
	if clientManager == nil {
		clientManager = kubeclient.NewManager()
	}
	return &Service{
		clusterRepo:     clusterRepo,
		credentialRepo:  credentialRepo,
		resourceRepo:    resourceRepo,
		resourceIndexer: resourceIndexer,
		syncDispatcher:  syncDispatcher,
		cipher:          cipher,
		clientManager:   clientManager,
	}
}

func (s *Service) RegisterCluster(ctx context.Context, req RegisterClusterRequest) (*domain.Cluster, error) {
	cluster := &domain.Cluster{
		Name:      strings.TrimSpace(req.Name),
		APIServer: strings.TrimSpace(req.APIServer),
		Status:    domain.ClusterStatusUnknown,
	}
	if err := s.clusterRepo.Create(ctx, cluster); err != nil {
		return nil, err
	}

	authType := strings.TrimSpace(req.AuthType)
	if authType == "" {
		authType = strings.TrimSpace(req.CredentialType)
	}
	if authType == "" {
		authType = "kubeconfig"
	}
	credentialPayload := req.KubeConfig
	if strings.TrimSpace(credentialPayload) == "" {
		credentialPayload = req.CredentialPayload
	}
	if strings.TrimSpace(credentialPayload) == "" {
		return nil, errors.New("credential payload is required")
	}
	ciphertext, err := s.cipher.Encrypt(credentialPayload)
	if err != nil {
		return nil, err
	}

	cred := &repository.ClusterCredential{
		ClusterID:            cluster.ID,
		AuthType:             authType,
		KubeConfigCiphertext: ciphertext,
	}
	if err := s.credentialRepo.UpsertByClusterID(ctx, cred); err != nil {
		return nil, err
	}
	if err := s.enqueueInitialSync(cluster.ID); err != nil {
		return nil, err
	}
	return cluster, nil
}

func (s *Service) enqueueInitialSync(clusterID uint64) error {
	if s.syncDispatcher != nil {
		return s.syncDispatcher.Enqueue(clusterID)
	}
	if s.resourceIndexer == nil {
		return nil
	}
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("resource sync panic recovered for cluster=%d: %v", clusterID, recovered)
			}
		}()
		if err := s.resourceIndexer.SyncCluster(context.Background(), clusterID); err != nil {
			log.Printf("resource sync failed for cluster=%d: %v", clusterID, err)
		}
	}()
	return nil
}

func (s *Service) TriggerResourceSync(ctx context.Context, clusterID uint64) (*SyncTriggerResult, error) {
	if s == nil {
		return nil, errors.New("cluster service is not initialized")
	}
	if clusterID == 0 {
		return nil, errors.New("clusterID is required")
	}
	if _, err := s.clusterRepo.GetByID(ctx, clusterID); err != nil {
		return nil, fmt.Errorf("failed to load cluster: %w", err)
	}

	if s.syncDispatcher != nil {
		if err := s.syncDispatcher.Enqueue(clusterID); err != nil {
			return nil, fmt.Errorf("failed to enqueue resource sync: %w", err)
		}
		return &SyncTriggerResult{
			ClusterID: clusterID,
			Status:    "accepted",
			Message:   "cluster resource sync task enqueued",
		}, nil
	}

	if s.resourceIndexer == nil {
		return nil, errors.New("resource sync is not configured")
	}

	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("resource sync panic recovered for cluster=%d: %v", clusterID, recovered)
			}
		}()
		if err := s.resourceIndexer.SyncCluster(context.Background(), clusterID); err != nil {
			log.Printf("resource sync failed for cluster=%d: %v", clusterID, err)
		}
	}()
	return &SyncTriggerResult{
		ClusterID: clusterID,
		Status:    "accepted",
		Message:   "cluster resource sync started in background",
	}, nil
}

func (s *Service) ValidateConnectivity(ctx context.Context, clusterID uint64) (*ConnectivityResult, error) {
	cluster, err := s.clusterRepo.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	cred, err := s.credentialRepo.GetByClusterID(ctx, clusterID)
	if err != nil {
		_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusDegraded)
		return &ConnectivityResult{
			ClusterID: clusterID,
			Success:   false,
			Message:   "cluster credential not found",
		}, nil
	}

	kubeConfig, err := s.cipher.Decrypt(cred.KubeConfigCiphertext)
	if err != nil {
		_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusDegraded)
		return &ConnectivityResult{
			ClusterID: clusterID,
			Success:   false,
			Message:   "failed to decrypt cluster credential",
		}, nil
	}

	cfg, err := s.clientManager.BuildRESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusDegraded)
		return &ConnectivityResult{
			ClusterID: clusterID,
			Success:   false,
			Message:   "invalid kubeconfig: " + strings.TrimSpace(err.Error()),
		}, nil
	}

	if strings.TrimSpace(cfg.Host) == "" {
		cfg.Host = strings.TrimSpace(cluster.APIServer)
	}
	if strings.TrimSpace(cfg.Host) == "" {
		_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusDegraded)
		return &ConnectivityResult{
			ClusterID: clusterID,
			Success:   false,
			Message:   "api server endpoint is empty",
		}, nil
	}

	if err := s.clientManager.Register(strconv.FormatUint(clusterID, 10), cfg); err != nil {
		_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusDegraded)
		return &ConnectivityResult{
			ClusterID: clusterID,
			Success:   false,
			Message:   "failed to initialize cluster client: " + strings.TrimSpace(err.Error()),
		}, nil
	}

	if err := s.clientManager.CheckConnectivity(cfg); err != nil {
		_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusDegraded)
		return &ConnectivityResult{
			ClusterID: clusterID,
			Success:   false,
			Message:   "connectivity probe failed: " + strings.TrimSpace(err.Error()),
		}, nil
	}

	_ = s.clusterRepo.UpdateStatus(ctx, clusterID, domain.ClusterStatusHealthy)
	return &ConnectivityResult{
		ClusterID: clusterID,
		Success:   true,
		Message:   fmt.Sprintf("connectivity probe succeeded for %s", cfg.Host),
	}, nil
}

func (s *Service) GetHealthSummary(ctx context.Context, clusterID uint64) (*HealthSummary, error) {
	_, err := s.clusterRepo.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	counts, err := s.resourceRepo.CountByHealth(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	summary := &HealthSummary{
		ClusterID: clusterID,
		Healthy:   counts["healthy"],
		Degraded:  counts["degraded"],
		Unknown:   counts["unknown"],
	}
	summary.Total = summary.Healthy + summary.Degraded + summary.Unknown
	return summary, nil
}

func (s *Service) ListResources(ctx context.Context, filter repository.ResourceListFilter) ([]repository.ResourceInventory, error) {
	return s.resourceRepo.List(ctx, filter)
}

func (s *Service) GetResourceDetail(ctx context.Context, filter repository.ResourceDetailFilter) (*repository.ResourceInventory, error) {
	return s.resourceRepo.GetDetail(ctx, filter)
}

func (s *Service) ListClusters(ctx context.Context) ([]domain.Cluster, error) {
	return s.clusterRepo.List(ctx)
}

func (s *Service) ListClustersByIDs(ctx context.Context, clusterIDs []uint64) ([]domain.Cluster, error) {
	if len(clusterIDs) == 0 {
		return []domain.Cluster{}, nil
	}
	return s.clusterRepo.ListByIDs(ctx, clusterIDs)
}
