package operation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

var (
	errOperatorIDRequired        = errors.New("operator id is required")
	errClusterIDRequired         = errors.New("clusterId is required")
	errResourceKindRequired      = errors.New("resourceKind is required")
	errOperationTypeRequired     = errors.New("operationType is required")
	errRiskConfirmationRequired  = errors.New("risk confirmation is required for high-risk operation")
	errIdempotentRequestInFlight = errors.New("idempotent request is in progress")
)

type SubmitOperationRequest struct {
	IdempotencyKey string         `json:"idempotencyKey"`
	ClusterID      uint64         `json:"clusterId"`
	WorkspaceID    uint64         `json:"workspaceId,omitempty"`
	ProjectID      uint64         `json:"projectId,omitempty"`
	ResourceUID    string         `json:"resourceUid,omitempty"`
	ResourceKind   string         `json:"resourceKind"`
	Namespace      string         `json:"namespace,omitempty"`
	Name           string         `json:"name,omitempty"`
	OperationType  string         `json:"operationType"`
	RiskLevel      string         `json:"riskLevel,omitempty"`
	RiskConfirmed  bool           `json:"riskConfirmed"`
	Payload        map[string]any `json:"payload,omitempty"`
}

type Service struct {
	repo        *repository.OperationRepository
	idempotency *IdempotencyService
	queue       QueueService
}

func NewService(repo *repository.OperationRepository, idempotency *IdempotencyService, queue QueueService) *Service {
	if idempotency == nil {
		idempotency = NewIdempotencyService(nil)
	}
	if queue == nil {
		queue = NewQueueService(nil)
	}
	return &Service{repo: repo, idempotency: idempotency, queue: queue}
}

func (s *Service) Submit(ctx context.Context, operatorID uint64, req SubmitOperationRequest) (*domain.OperationRequest, bool, error) {
	if operatorID == 0 {
		return nil, false, errOperatorIDRequired
	}
	if req.ClusterID == 0 {
		return nil, false, errClusterIDRequired
	}
	resourceKind := strings.TrimSpace(req.ResourceKind)
	if resourceKind == "" {
		return nil, false, errResourceKindRequired
	}
	operationType := strings.TrimSpace(req.OperationType)
	if operationType == "" {
		return nil, false, errOperationTypeRequired
	}

	riskLevel := normalizeRiskLevel(req.RiskLevel, operationType)
	if riskLevel == domain.RiskLevelHigh && !req.RiskConfirmed {
		return nil, false, errRiskConfirmationRequired
	}

	requestID := s.idempotency.BuildRequestID(operatorID, req.IdempotencyKey)
	if strings.TrimSpace(req.IdempotencyKey) != "" {
		acquired, err := s.idempotency.Acquire(ctx, requestID)
		if err != nil {
			return nil, false, err
		}
		if acquired {
			defer func() {
				_ = s.idempotency.Release(ctx, requestID)
			}()
		}

		existing, err := s.repo.GetByRequestID(ctx, requestID)
		switch {
		case err == nil:
			return existing, true, nil
		case !errors.Is(err, gorm.ErrRecordNotFound):
			return nil, false, err
		}

		if !acquired {
			return nil, false, errIdempotentRequestInFlight
		}
	}

	item := &domain.OperationRequest{
		RequestID:     requestID,
		OperatorID:    operatorID,
		OperationType: operationType,
		TargetRef:     buildTargetRef(req),
		Status:        domain.OperationStatusPending,
		RiskLevel:     riskLevel,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, false, err
	}
	if err := s.queue.Enqueue(ctx, item.ID); err != nil {
		return nil, false, err
	}

	return item, false, nil
}

func (s *Service) GetByID(ctx context.Context, operationID uint64) (*domain.OperationRequest, error) {
	if operationID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return s.repo.GetByID(ctx, operationID)
}

func normalizeRiskLevel(rawRiskLevel, operationType string) domain.RiskLevel {
	normalized := strings.ToLower(strings.TrimSpace(rawRiskLevel))
	switch normalized {
	case "high", "critical":
		return domain.RiskLevelHigh
	case "medium":
		return domain.RiskLevelMedium
	case "low":
		return domain.RiskLevelLow
	}

	switch strings.ToLower(strings.TrimSpace(operationType)) {
	case "delete", "drain", "cordon":
		return domain.RiskLevelHigh
	case "scale", "restart", "update":
		return domain.RiskLevelMedium
	default:
		return domain.RiskLevelLow
	}
}

func buildTargetRef(req SubmitOperationRequest) string {
	parts := []string{
		fmt.Sprintf("cluster:%d", req.ClusterID),
	}
	if ns := strings.TrimSpace(req.Namespace); ns != "" {
		parts = append(parts, "ns:"+ns)
	}
	if kind := strings.TrimSpace(req.ResourceKind); kind != "" {
		parts = append(parts, "kind:"+kind)
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		parts = append(parts, "name:"+name)
	}
	if uid := strings.TrimSpace(req.ResourceUID); uid != "" {
		parts = append(parts, "uid:"+uid)
	}
	return strings.Join(parts, "/")
}
