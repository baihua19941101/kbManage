package client

import (
	"fmt"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
