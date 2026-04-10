package cluster

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type ConnectivityResult struct {
	ClusterID uint64 `json:"clusterId"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

type RegisterClusterRequest struct {
	Name       string `json:"name"`
	APIServer  string `json:"apiServer"`
	AuthType   string `json:"authType"`
	KubeConfig string `json:"kubeConfig"`
}

type HealthSummary struct {
	ClusterID uint64 `json:"clusterId"`
	Healthy   int64  `json:"healthy"`
	Degraded  int64  `json:"degraded"`
	Unknown   int64  `json:"unknown"`
	Total     int64  `json:"total"`
}

type Service struct {
	clusterRepo    *repository.ClusterRepository
	credentialRepo *repository.ClusterCredentialRepository
	resourceRepo   *repository.ResourceInventoryRepository
	cipher         CredentialCipher
}

func NewService(
	clusterRepo *repository.ClusterRepository,
	credentialRepo *repository.ClusterCredentialRepository,
	resourceRepo *repository.ResourceInventoryRepository,
	cipher CredentialCipher,
) *Service {
	return &Service{
		clusterRepo:    clusterRepo,
		credentialRepo: credentialRepo,
		resourceRepo:   resourceRepo,
		cipher:         cipher,
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
		authType = "kubeconfig"
	}
	ciphertext, err := s.cipher.Encrypt(req.KubeConfig)
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
	return cluster, nil
}

func (s *Service) ValidateConnectivity(ctx context.Context, clusterID uint64) (*ConnectivityResult, error) {
	_, err := s.clusterRepo.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	return &ConnectivityResult{
		ClusterID: clusterID,
		Success:   true,
		Message:   "connectivity check stub: accepted",
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
