package client

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Manager struct {
	mu      sync.RWMutex
	clients map[string]kubernetes.Interface
}

func NewManager() *Manager {
	return &Manager{clients: make(map[string]kubernetes.Interface)}
}

func (m *Manager) Register(clusterID string, cfg *rest.Config) error {
	if clusterID == "" {
		return fmt.Errorf("clusterID is required")
	}
	if cfg == nil {
		return fmt.Errorf("rest config is required")
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[clusterID] = clientset
	return nil
}

func (m *Manager) Get(clusterID string) (kubernetes.Interface, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[clusterID]
	return c, ok
}

func (m *Manager) BuildRESTConfigFromKubeConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig == "" {
		return nil, errors.New("kubeconfig is required")
	}
	cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (m *Manager) CheckConnectivity(cfg *rest.Config) error {
	if cfg == nil {
		return errors.New("rest config is required")
	}

	probeConfig := rest.CopyConfig(cfg)
	if probeConfig.Timeout <= 0 {
		probeConfig.Timeout = 5 * time.Second
	}

	clientset, err := kubernetes.NewForConfig(probeConfig)
	if err != nil {
		return err
	}

	_, err = clientset.Discovery().ServerVersion()
	return err
}
