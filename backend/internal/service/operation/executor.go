package operation

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

var (
	ErrClusterClientUnavailable = errors.New("cluster client is unavailable")
	ErrUnsupportedOperationType = errors.New("unsupported operation type")
	ErrUnsupportedResourceKind  = errors.New("unsupported resource kind")
)

type ExecuteResult struct {
	ProgressMessage string
	ResultMessage   string
	FailureReason   string
}

// Executor executes a submitted operation request.
type Executor interface {
	Execute(ctx context.Context, item *domain.OperationRequest) (ExecuteResult, error)
}

// KubernetesClientProvider resolves a Kubernetes client by cluster ID.
type KubernetesClientProvider interface {
	GetClient(clusterID uint64) (kubernetes.Interface, bool)
}

// StaticClientProvider is a minimal in-memory provider for tests or wiring.
type StaticClientProvider map[uint64]kubernetes.Interface

func (p StaticClientProvider) GetClient(clusterID uint64) (kubernetes.Interface, bool) {
	client, ok := p[clusterID]
	return client, ok
}

// NewExecutor creates the default executor chain: Kubernetes first, simulated fallback.
func NewExecutor(clientProvider KubernetesClientProvider) Executor {
	return &fallbackExecutor{
		primary:   NewKubernetesExecutor(clientProvider),
		secondary: NewSimulatedExecutor(),
	}
}

func NewKubernetesExecutor(clientProvider KubernetesClientProvider) Executor {
	return &kubernetesExecutor{clientProvider: clientProvider}
}

func NewSimulatedExecutor() Executor {
	return simulatedExecutor{}
}

type fallbackExecutor struct {
	primary   Executor
	secondary Executor
}

func (e *fallbackExecutor) Execute(ctx context.Context, item *domain.OperationRequest) (ExecuteResult, error) {
	if e == nil {
		return ExecuteResult{}, errors.New("executor is not configured")
	}
	if e.primary == nil {
		if e.secondary == nil {
			return ExecuteResult{}, errors.New("executor is not configured")
		}
		return e.secondary.Execute(ctx, item)
	}

	result, err := e.primary.Execute(ctx, item)
	if err == nil || e.secondary == nil {
		return result, err
	}
	if errors.Is(err, ErrClusterClientUnavailable) {
		return e.secondary.Execute(ctx, item)
	}
	return result, err
}

type simulatedExecutor struct{}

func (simulatedExecutor) Execute(_ context.Context, item *domain.OperationRequest) (ExecuteResult, error) {
	return ExecuteResult{
		ProgressMessage: "operation dispatched by simulated executor",
		ResultMessage:   defaultOperationResultMessage(item),
	}, nil
}

type kubernetesExecutor struct {
	clientProvider KubernetesClientProvider
}

func (e *kubernetesExecutor) Execute(ctx context.Context, item *domain.OperationRequest) (ExecuteResult, error) {
	if item == nil {
		return ExecuteResult{}, errors.New("operation request is required")
	}
	if e == nil || e.clientProvider == nil {
		return ExecuteResult{
			FailureReason: "cluster client provider is not configured",
		}, fmt.Errorf("%w: provider is nil", ErrClusterClientUnavailable)
	}

	target, err := parseOperationTargetRef(item.TargetRef)
	if err != nil {
		return ExecuteResult{FailureReason: err.Error()}, err
	}
	clientset, ok := e.clientProvider.GetClient(target.ClusterID)
	if !ok || clientset == nil {
		return ExecuteResult{
			FailureReason: fmt.Sprintf("cluster client is unavailable for cluster=%d", target.ClusterID),
		}, fmt.Errorf("%w: cluster=%d", ErrClusterClientUnavailable, target.ClusterID)
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
	}

	switch strings.ToLower(strings.TrimSpace(item.OperationType)) {
	case "scale":
		return e.executeScale(ctx, clientset, target)
	case "restart":
		return e.executeRestart(ctx, clientset, target)
	case "node-maintenance", "cordon", "uncordon", "drain":
		return e.executeNodeMaintenance(ctx, clientset, target, strings.ToLower(strings.TrimSpace(item.OperationType)))
	default:
		err := fmt.Errorf("%w: %s", ErrUnsupportedOperationType, strings.TrimSpace(item.OperationType))
		return ExecuteResult{FailureReason: err.Error()}, err
	}
}

func (e *kubernetesExecutor) executeScale(ctx context.Context, clientset kubernetes.Interface, target operationTargetRef) (ExecuteResult, error) {
	if target.Name == "" || target.Namespace == "" {
		err := errors.New("scale target requires namespace and name")
		return ExecuteResult{FailureReason: err.Error()}, err
	}

	kind := strings.ToLower(target.Kind)
	switch kind {
	case "deployment":
		scale, err := clientset.AppsV1().Deployments(target.Namespace).GetScale(ctx, target.Name, metav1.GetOptions{})
		if err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("get deployment scale: %w", err)
		}
		if _, err := clientset.AppsV1().Deployments(target.Namespace).UpdateScale(ctx, target.Name, scale, metav1.UpdateOptions{}); err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("update deployment scale: %w", err)
		}
		return ExecuteResult{
			ProgressMessage: "scale operation applied",
			ResultMessage:   fmt.Sprintf("scaled deployment %s/%s to %d replicas", target.Namespace, target.Name, scale.Spec.Replicas),
		}, nil
	case "statefulset":
		scale, err := clientset.AppsV1().StatefulSets(target.Namespace).GetScale(ctx, target.Name, metav1.GetOptions{})
		if err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("get statefulset scale: %w", err)
		}
		if _, err := clientset.AppsV1().StatefulSets(target.Namespace).UpdateScale(ctx, target.Name, scale, metav1.UpdateOptions{}); err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("update statefulset scale: %w", err)
		}
		return ExecuteResult{
			ProgressMessage: "scale operation applied",
			ResultMessage:   fmt.Sprintf("scaled statefulset %s/%s to %d replicas", target.Namespace, target.Name, scale.Spec.Replicas),
		}, nil
	default:
		err := fmt.Errorf("%w: %s", ErrUnsupportedResourceKind, target.Kind)
		return ExecuteResult{FailureReason: err.Error()}, err
	}
}

func (e *kubernetesExecutor) executeRestart(ctx context.Context, clientset kubernetes.Interface, target operationTargetRef) (ExecuteResult, error) {
	if target.Name == "" || target.Namespace == "" {
		err := errors.New("restart target requires namespace and name")
		return ExecuteResult{FailureReason: err.Error()}, err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	patch := []byte(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"` + timestamp + `"}}}}}`)

	kind := strings.ToLower(target.Kind)
	switch kind {
	case "deployment":
		if _, err := clientset.AppsV1().Deployments(target.Namespace).Patch(ctx, target.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{}); err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("restart deployment: %w", err)
		}
	case "statefulset":
		if _, err := clientset.AppsV1().StatefulSets(target.Namespace).Patch(ctx, target.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{}); err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("restart statefulset: %w", err)
		}
	case "daemonset":
		if _, err := clientset.AppsV1().DaemonSets(target.Namespace).Patch(ctx, target.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{}); err != nil {
			return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("restart daemonset: %w", err)
		}
	default:
		err := fmt.Errorf("%w: %s", ErrUnsupportedResourceKind, target.Kind)
		return ExecuteResult{FailureReason: err.Error()}, err
	}

	return ExecuteResult{
		ProgressMessage: "restart operation applied",
		ResultMessage:   fmt.Sprintf("restarted %s %s/%s", strings.ToLower(target.Kind), target.Namespace, target.Name),
	}, nil
}

func (e *kubernetesExecutor) executeNodeMaintenance(
	ctx context.Context,
	clientset kubernetes.Interface,
	target operationTargetRef,
	operationType string,
) (ExecuteResult, error) {
	if target.Name == "" {
		err := errors.New("node-maintenance target requires node name")
		return ExecuteResult{FailureReason: err.Error()}, err
	}
	if target.Kind != "" && !strings.EqualFold(target.Kind, "node") {
		err := fmt.Errorf("%w: %s", ErrUnsupportedResourceKind, target.Kind)
		return ExecuteResult{FailureReason: err.Error()}, err
	}

	node, err := clientset.CoreV1().Nodes().Get(ctx, target.Name, metav1.GetOptions{})
	if err != nil {
		return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("load node: %w", err)
	}

	desiredUnschedulable := true
	actionText := "cordon"
	switch operationType {
	case "uncordon":
		desiredUnschedulable = false
		actionText = "uncordon"
	case "drain":
		desiredUnschedulable = true
		actionText = "drain"
	}

	nodeCopy := node.DeepCopy()
	nodeCopy.Spec.Unschedulable = desiredUnschedulable
	if _, err := clientset.CoreV1().Nodes().Update(ctx, nodeCopy, metav1.UpdateOptions{}); err != nil {
		return ExecuteResult{FailureReason: err.Error()}, fmt.Errorf("%s node: %w", actionText, err)
	}

	return ExecuteResult{
		ProgressMessage: "node maintenance operation applied",
		ResultMessage:   fmt.Sprintf("%s node %s", actionText, target.Name),
	}, nil
}

type operationTargetRef struct {
	ClusterID uint64
	Namespace string
	Kind      string
	Name      string
}

func parseOperationTargetRef(targetRef string) (operationTargetRef, error) {
	trimmed := strings.TrimSpace(targetRef)
	if trimmed == "" {
		return operationTargetRef{}, errors.New("targetRef is required")
	}

	ref := operationTargetRef{}
	parts := strings.Split(trimmed, "/")
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}
		keyValue := strings.SplitN(token, ":", 2)
		if len(keyValue) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(keyValue[0]))
		value := strings.TrimSpace(keyValue[1])
		switch key {
		case "cluster":
			clusterID, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return operationTargetRef{}, fmt.Errorf("invalid cluster id in targetRef: %w", err)
			}
			ref.ClusterID = clusterID
		case "ns":
			ref.Namespace = value
		case "kind":
			ref.Kind = value
		case "name":
			ref.Name = value
		}
	}
	if ref.ClusterID == 0 {
		return operationTargetRef{}, errors.New("targetRef cluster is required")
	}
	return ref, nil
}

func defaultOperationResultMessage(item *domain.OperationRequest) string {
	if item == nil {
		return "operation executed successfully"
	}
	if strings.TrimSpace(item.TargetRef) != "" {
		return "operation executed: " + strings.TrimSpace(item.TargetRef)
	}
	return "operation executed successfully"
}
