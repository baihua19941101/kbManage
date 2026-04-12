package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	obsSvc "kbmanage/backend/internal/service/observability"
)

type ObservabilityHandler struct {
	service *obsSvc.Service
}

func NewObservabilityHandler(service *obsSvc.Service) *ObservabilityHandler {
	if service == nil {
		service = obsSvc.NewService(nil)
	}
	return &ObservabilityHandler{service: service}
}

func (h *ObservabilityHandler) Overview(c *gin.Context) {
	clusterIDs, err := parseUint64List(c.QueryArray("clusterIds"), c.Query("clusterIds"), c.Query("clusterId"))
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	res, err := h.service.Overview(c.Request.Context(), obsSvc.OverviewRequest{
		ClusterIDs: clusterIDs,
		StartAt:    c.Query("startAt"),
		EndAt:      c.Query("endAt"),
	})
	if err != nil {
		writeObservabilityServiceError(c, "overview_failed", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ObservabilityHandler) QueryLogs(c *gin.Context) {
	clusterIDs, err := parseUint64List(c.QueryArray("clusterIds"), c.Query("clusterIds"), c.Query("clusterId"))
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "limit must be integer")
		return
	}

	namespace := firstNonEmpty(c.QueryArray("namespaces"), c.Query("namespaces"), c.Query("namespace"))
	res, err := h.service.QueryLogs(c.Request.Context(), obsSvc.LogsQueryRequest{
		ClusterIDs: clusterIDs,
		Namespace:  namespace,
		Workload:   c.Query("workload"),
		Pod:        c.Query("pod"),
		Container:  c.Query("container"),
		Keyword:    c.Query("keyword"),
		StartAt:    c.Query("startAt"),
		EndAt:      c.Query("endAt"),
		Limit:      limit,
	})
	if err != nil {
		writeObservabilityServiceError(c, "query_logs_failed", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ObservabilityHandler) ListEvents(c *gin.Context) {
	clusterID := firstNonEmpty(c.QueryArray("clusterIds"), c.Query("clusterIds"), c.Query("clusterId"))
	namespace := firstNonEmpty(c.QueryArray("namespaces"), c.Query("namespaces"), c.Query("namespace"))

	res, err := h.service.ListEvents(c.Request.Context(), obsSvc.EventsQueryRequest{
		ClusterID:    firstToken(clusterID),
		Namespace:    firstToken(namespace),
		ResourceKind: c.Query("resourceKind"),
		ResourceName: c.Query("resourceName"),
		EventType:    c.Query("eventType"),
		StartAt:      c.Query("startAt"),
		EndAt:        c.Query("endAt"),
	})
	if err != nil {
		writeObservabilityServiceError(c, "list_events_failed", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ObservabilityHandler) QueryMetricSeries(c *gin.Context) {
	clusterIDs, err := parseUint64List(c.QueryArray("clusterIds"), c.Query("clusterIds"), c.Query("clusterId"))
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	subjectType := c.Query("subjectType")
	subjectRef := c.Query("subjectRef")
	metricKey := c.Query("metricKey")
	if subjectType == "" || subjectRef == "" || metricKey == "" {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "subjectType, subjectRef and metricKey are required")
		return
	}

	res, err := h.service.QueryMetricSeries(c.Request.Context(), obsSvc.MetricQueryRequest{
		ClusterIDs:  clusterIDs,
		SubjectType: subjectType,
		SubjectRef:  subjectRef,
		MetricKey:   metricKey,
		StartAt:     c.Query("startAt"),
		EndAt:       c.Query("endAt"),
		Step:        c.Query("step"),
	})
	if err != nil {
		writeObservabilityServiceError(c, "query_metrics_failed", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ObservabilityHandler) ResourceContext(c *gin.Context) {
	clusterID := c.Query("clusterId")
	namespace := c.Query("namespace")
	resourceKind := c.Query("resourceKind")
	resourceName := c.Query("resourceName")
	if clusterID == "" || namespace == "" || resourceKind == "" || resourceName == "" {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "clusterId, namespace, resourceKind, resourceName are required")
		return
	}

	res, err := h.service.ResourceContext(c.Request.Context(), obsSvc.ResourceContextQuery{
		ClusterID:    clusterID,
		Namespace:    namespace,
		ResourceKind: resourceKind,
		ResourceName: resourceName,
		Keyword:      c.Query("keyword"),
		StartAt:      c.Query("startAt"),
		EndAt:        c.Query("endAt"),
	})
	if err != nil {
		writeObservabilityServiceError(c, "resource_context_failed", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func writeObservabilityError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

func writeObservabilityServiceError(c *gin.Context, code string, err error) {
	if errors.Is(err, obsSvc.ErrObservabilityScopeDenied) || errors.Is(err, obsSvc.ErrInvalidObservabilityUser) {
		writeObservabilityError(c, http.StatusForbidden, "forbidden", "observability scope access denied")
		return
	}
	writeObservabilityError(c, http.StatusInternalServerError, code, err.Error())
}

func parseUint64List(arr []string, csv string, single string) ([]uint64, error) {
	rawItems := make([]string, 0, len(arr)+2)
	rawItems = append(rawItems, arr...)
	if csv != "" {
		rawItems = append(rawItems, strings.Split(csv, ",")...)
	}
	if single != "" {
		rawItems = append(rawItems, single)
	}

	out := make([]uint64, 0, len(rawItems))
	for _, item := range rawItems {
		token := strings.TrimSpace(item)
		if token == "" {
			continue
		}
		v, err := strconv.ParseUint(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid clusterIds value: %q", token)
		}
		out = append(out, v)
	}
	return out, nil
}

func firstNonEmpty(values []string, fallbacks ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	for _, v := range fallbacks {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func firstToken(in string) string {
	parts := strings.Split(in, ",")
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token != "" {
			return token
		}
	}
	return ""
}
