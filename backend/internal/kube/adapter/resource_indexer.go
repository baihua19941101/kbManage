package adapter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	kubeclient "kbmanage/backend/internal/kube/client"
	"kbmanage/backend/internal/repository"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ResourceIndexer indexes Kubernetes resources for a cluster.
type ResourceIndexer interface {
	SyncCluster(ctx context.Context, clusterID uint64) error
}

type credentialDecryptor interface {
	Decrypt(ciphertext string) (string, error)
}

type KubeResourceIndexer struct {
	clusterRepo    *repository.ClusterRepository
	credentialRepo *repository.ClusterCredentialRepository
	resourceRepo   *repository.ResourceInventoryRepository
	cipher         credentialDecryptor
	clientManager  *kubeclient.Manager
}

type NoopResourceIndexer struct{}

func NewResourceIndexer(
	clusterRepo *repository.ClusterRepository,
	credentialRepo *repository.ClusterCredentialRepository,
	resourceRepo *repository.ResourceInventoryRepository,
	cipher credentialDecryptor,
	clientManager *kubeclient.Manager,
) ResourceIndexer {
	if clusterRepo == nil || credentialRepo == nil || resourceRepo == nil || cipher == nil {
		return NoopResourceIndexer{}
	}
	if clientManager == nil {
		clientManager = kubeclient.NewManager()
	}
	return &KubeResourceIndexer{
		clusterRepo:    clusterRepo,
		credentialRepo: credentialRepo,
		resourceRepo:   resourceRepo,
		cipher:         cipher,
		clientManager:  clientManager,
	}
}

func (NoopResourceIndexer) SyncCluster(_ context.Context, _ uint64) error {
	return errors.New("resource indexer is not configured")
}

func (i *KubeResourceIndexer) SyncCluster(ctx context.Context, clusterID uint64) error {
	if clusterID == 0 {
		return errors.New("clusterID is required")
	}

	if i.clusterRepo != nil {
		_ = i.clusterRepo.MarkSyncRunning(ctx, clusterID)
	}
	var syncErr error
	defer func() {
		if i.clusterRepo == nil {
			return
		}
		if syncErr != nil {
			_ = i.clusterRepo.MarkSyncFailed(ctx, clusterID, syncErr.Error())
			return
		}
		_ = i.clusterRepo.MarkSyncSuccess(ctx, clusterID)
	}()

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
	}

	clientset, err := i.resolveClient(ctx, clusterID)
	if err != nil {
		syncErr = err
		return syncErr
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		syncErr = fmt.Errorf("list namespaces: %w", err)
		return syncErr
	}

	items := make([]repository.ResourceInventory, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		name := strings.TrimSpace(ns.Name)
		if name == "" {
			continue
		}

		health := "unknown"
		switch strings.ToLower(string(ns.Status.Phase)) {
		case "active":
			health = "healthy"
		case "terminating":
			health = "degraded"
		}

		items = append(items, repository.ResourceInventory{
			ClusterID: clusterID,
			Namespace: name,
			Kind:      "Namespace",
			Name:      name,
			Health:    health,
		})
	}

	if err := i.resourceRepo.ReplaceClusterSnapshot(ctx, clusterID, items); err != nil {
		syncErr = fmt.Errorf("replace resource snapshot: %w", err)
		return syncErr
	}
	return nil
}

func (i *KubeResourceIndexer) resolveClient(ctx context.Context, clusterID uint64) (kubernetes.Interface, error) {
	clusterKey := strconv.FormatUint(clusterID, 10)
	if clientset, ok := i.clientManager.Get(clusterKey); ok {
		return clientset, nil
	}

	cluster, err := i.clusterRepo.GetByID(ctx, clusterID)
	if err != nil {
		return nil, fmt.Errorf("load cluster: %w", err)
	}
	cred, err := i.credentialRepo.GetByClusterID(ctx, clusterID)
	if err != nil {
		return nil, fmt.Errorf("load cluster credential: %w", err)
	}

	kubeConfig, err := i.cipher.Decrypt(cred.KubeConfigCiphertext)
	if err != nil {
		return nil, fmt.Errorf("decrypt kubeconfig: %w", err)
	}

	cfg, err := i.clientManager.BuildRESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("build rest config: %w", err)
	}
	if strings.TrimSpace(cfg.Host) == "" {
		cfg.Host = strings.TrimSpace(cluster.APIServer)
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return nil, errors.New("api server endpoint is empty")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 8 * time.Second
	}

	if err := i.clientManager.Register(clusterKey, cfg); err != nil {
		return nil, fmt.Errorf("register cluster client: %w", err)
	}
	clientset, ok := i.clientManager.Get(clusterKey)
	if !ok {
		return nil, errors.New("cluster client is unavailable after registration")
	}
	return clientset, nil
}
