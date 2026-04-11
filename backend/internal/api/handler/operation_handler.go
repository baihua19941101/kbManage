package handler

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	authSvc "kbmanage/backend/internal/service/auth"
	operationSvc "kbmanage/backend/internal/service/operation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationHandler struct {
	svc         *operationSvc.Service
	scopeAccess *authSvc.ScopeAccessService
}

func NewOperationHandler(svc *operationSvc.Service, scopeAccess *authSvc.ScopeAccessService) *OperationHandler {
	return &OperationHandler{svc: svc, scopeAccess: scopeAccess}
}

type createOperationRequest struct {
	IdempotencyKey string         `json:"idempotencyKey"`
	ClusterID      any            `json:"clusterId"`
	WorkspaceID    any            `json:"workspaceId"`
	ProjectID      any            `json:"projectId"`
	ResourceUID    string         `json:"resourceUid"`
	ResourceKind   string         `json:"resourceKind"`
	Namespace      string         `json:"namespace"`
	Name           string         `json:"name"`
	OperationType  string         `json:"operationType"`
	RiskLevel      string         `json:"riskLevel"`
	RiskConfirmed  bool           `json:"riskConfirmed"`
	Payload        map[string]any `json:"payload"`
}

func (h *OperationHandler) Create(c *gin.Context) {
	var req createOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clusterID, err := parseUint64(req.ClusterID, true, "clusterId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := parseUint64(req.WorkspaceID, false, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseUint64(req.ProjectID, false, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, duplicated, err := h.svc.Submit(c.Request.Context(), c.GetUint64(middleware.UserIDKey), operationSvc.SubmitOperationRequest{
		IdempotencyKey: strings.TrimSpace(req.IdempotencyKey),
		ClusterID:      clusterID,
		WorkspaceID:    workspaceID,
		ProjectID:      projectID,
		ResourceUID:    strings.TrimSpace(req.ResourceUID),
		ResourceKind:   strings.TrimSpace(req.ResourceKind),
		Namespace:      strings.TrimSpace(req.Namespace),
		Name:           strings.TrimSpace(req.Name),
		OperationType:  strings.TrimSpace(req.OperationType),
		RiskLevel:      strings.TrimSpace(req.RiskLevel),
		RiskConfirmed:  req.RiskConfirmed,
		Payload:        req.Payload,
	})
	if err != nil {
		writeOperationError(c, err)
		return
	}

	status := http.StatusAccepted
	if duplicated {
		status = http.StatusOK
	}
	c.JSON(status, toOperationResponse(item))
}

func (h *OperationHandler) GetByID(c *gin.Context) {
	operationID, err := strconv.ParseUint(c.Param("operationId"), 10, 64)
	if err != nil || operationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation id"})
		return
	}
	userID := c.GetUint64(middleware.UserIDKey)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}

	item, err := h.svc.GetByID(c.Request.Context(), operationID)
	if err != nil {
		writeOperationError(c, err)
		return
	}
	allowed, err := h.canViewOperation(c, userID, item)
	if err != nil {
		writeOperationError(c, err)
		return
	}
	if !allowed {
		writeOperationError(c, gorm.ErrRecordNotFound)
		return
	}
	c.JSON(http.StatusOK, toOperationResponse(item))
}

func (h *OperationHandler) canViewOperation(c *gin.Context, userID uint64, item *domain.OperationRequest) (bool, error) {
	if item == nil {
		return false, gorm.ErrRecordNotFound
	}
	if item.OperatorID == userID || h.scopeAccess == nil {
		return true, nil
	}

	clusterID, ok := authSvc.ParseClusterIDFromReference(item.TargetRef)
	if !ok {
		return false, nil
	}
	return h.scopeAccess.CanAccessClusterByPermission(c.Request.Context(), userID, clusterID, middleware.PermissionProjectRead)
}

func toOperationResponse(item *domain.OperationRequest) gin.H {
	return gin.H{
		"id":              item.ID,
		"requestId":       item.RequestID,
		"operatorId":      item.OperatorID,
		"operationType":   item.OperationType,
		"targetRef":       item.TargetRef,
		"status":          item.Status,
		"riskLevel":       item.RiskLevel,
		"progressMessage": item.ProgressMessage,
		"resultMessage":   item.ResultMessage,
		"failureReason":   item.FailureReason,
		"completedAt":     item.CompletedAt,
		"createdAt":       item.CreatedAt,
		"updatedAt":       item.UpdatedAt,
	}
}

func writeOperationError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	lower := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		status = http.StatusConflict
	case strings.Contains(lower, "in progress"):
		status = http.StatusConflict
	case strings.Contains(lower, "required"), strings.Contains(lower, "invalid"):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": err.Error()})
}

func parseUint64(v any, required bool, fieldName string) (uint64, error) {
	switch value := v.(type) {
	case nil:
		if required {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return 0, nil
	case uint64:
		if required && value == 0 {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return value, nil
	case uint:
		if required && value == 0 {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return uint64(value), nil
	case int:
		if value < 0 {
			return 0, fmt.Errorf("%s must be a positive integer", fieldName)
		}
		if required && value == 0 {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return uint64(value), nil
	case int64:
		if value < 0 {
			return 0, fmt.Errorf("%s must be a positive integer", fieldName)
		}
		if required && value == 0 {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return uint64(value), nil
	case float64:
		if value < 0 || value > math.MaxUint64 || math.Trunc(value) != value {
			return 0, fmt.Errorf("%s must be a positive integer", fieldName)
		}
		if required && value == 0 {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return uint64(value), nil
	case string:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			if required {
				return 0, fmt.Errorf("%s is required", fieldName)
			}
			return 0, nil
		}
		n, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a positive integer", fieldName)
		}
		if required && n == 0 {
			return 0, fmt.Errorf("%s is required", fieldName)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("%s has unsupported type", fieldName)
	}
}
