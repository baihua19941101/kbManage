package adapter

import (
	"context"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "kbmanage/backend/internal/kube/client"
)

type EventItem struct {
	ClusterID    string
	Namespace    string
	InvolvedKind string
	InvolvedName string
	EventType    string
	Reason       string
	Message      string
	FirstSeenAt  string
	LastSeenAt   string
	Count        int
}

type EventReader interface {
	List(ctx context.Context, clusterID, namespace, kind, name string) ([]EventItem, error)
}

type MockEventReader struct{}

func NewMockEventReader() *MockEventReader {
	return &MockEventReader{}
}

func (r *MockEventReader) List(_ context.Context, clusterID, namespace, kind, name string) ([]EventItem, error) {
	now := time.Now().UTC()
	return []EventItem{
		{
			ClusterID:    clusterID,
			Namespace:    namespace,
			InvolvedKind: kind,
			InvolvedName: name,
			EventType:    "warning",
			Reason:       "BackOff",
			Message:      "mock event for observability timeline",
			FirstSeenAt:  now.Add(-10 * time.Minute).Format(time.RFC3339),
			LastSeenAt:   now.Format(time.RFC3339),
			Count:        2,
		},
	}, nil
}

type KubeEventReader struct {
	clientManager *kubeclient.Manager
	fallback      EventReader
}

func NewKubeEventReader(clientManager *kubeclient.Manager) *KubeEventReader {
	return &KubeEventReader{
		clientManager: clientManager,
		fallback:      NewMockEventReader(),
	}
}

func (r *KubeEventReader) List(ctx context.Context, clusterID, namespace, kind, name string) ([]EventItem, error) {
	if r == nil || r.clientManager == nil || clusterID == "" {
		return r.fallback.List(ctx, clusterID, namespace, kind, name)
	}

	clientset, ok := r.clientManager.Get(clusterID)
	if !ok {
		return r.fallback.List(ctx, clusterID, namespace, kind, name)
	}

	ns := namespace
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	list, err := clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return r.fallback.List(ctx, clusterID, namespace, kind, name)
	}

	out := make([]EventItem, 0, len(list.Items))
	for _, event := range list.Items {
		if kind != "" && !strings.EqualFold(event.InvolvedObject.Kind, kind) {
			continue
		}
		if name != "" && event.InvolvedObject.Name != name {
			continue
		}
		eventType := strings.ToLower(event.Type)
		if eventType == "" {
			eventType = "normal"
		}
		out = append(out, EventItem{
			ClusterID:    clusterID,
			Namespace:    event.Namespace,
			InvolvedKind: event.InvolvedObject.Kind,
			InvolvedName: event.InvolvedObject.Name,
			EventType:    eventType,
			Reason:       event.Reason,
			Message:      event.Message,
			FirstSeenAt:  event.FirstTimestamp.UTC().Format(time.RFC3339),
			LastSeenAt:   event.LastTimestamp.UTC().Format(time.RFC3339),
			Count:        int(event.Count),
		})
	}
	if len(out) == 0 {
		return r.fallback.List(ctx, clusterID, namespace, kind, name)
	}
	return out, nil
}
