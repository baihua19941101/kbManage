package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	complianceSvc "kbmanage/backend/internal/service/compliance"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceHandler struct {
	svc         *complianceSvc.Service
	remediation *complianceSvc.RemediationService
	exceptions  *complianceSvc.ExceptionService
	rechecks    *complianceSvc.RecheckService
	overview    *complianceSvc.OverviewService
	trends      *complianceSvc.TrendService
	exports     *complianceSvc.ArchiveExportService
}

func NewComplianceHandler(svc *complianceSvc.Service, remediation *complianceSvc.RemediationService, exceptions *complianceSvc.ExceptionService, rechecks *complianceSvc.RecheckService, overview *complianceSvc.OverviewService, trends *complianceSvc.TrendService, exports *complianceSvc.ArchiveExportService) *ComplianceHandler {
	return &ComplianceHandler{svc: svc, remediation: remediation, exceptions: exceptions, rechecks: rechecks, overview: overview, trends: trends, exports: exports}
}

func (h *ComplianceHandler) ListBaselines(c *gin.Context) {
	items, err := h.svc.Baselines.List(c.Request.Context(), c.GetUint64(middleware.UserIDKey), complianceSvc.BaselineListFilter{StandardType: c.Query("standardType"), Status: c.Query("status")})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) GetBaseline(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "baselineId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Baselines.Get(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) CreateBaseline(c *gin.Context) {
	var req complianceSvc.CreateBaselineInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Baselines.Create(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
func (h *ComplianceHandler) UpdateBaseline(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "baselineId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req complianceSvc.UpdateBaselineInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Baselines.Update(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) ListScanProfiles(c *gin.Context) {
	workspaceID, _ := parseComplianceOptionalQueryUint64(c, "workspaceId")
	projectID, _ := parseComplianceOptionalQueryUint64(c, "projectId")
	items, err := h.svc.Profiles.List(c.Request.Context(), c.GetUint64(middleware.UserIDKey), complianceSvc.ScanProfileListFilter{WorkspaceID: derefUint64Value(workspaceID), ProjectID: derefUint64Value(projectID), ScopeType: c.Query("scopeType"), ScheduleMode: c.Query("scheduleMode"), Status: c.Query("status")})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) GetScanProfile(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "profileId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Profiles.Get(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) CreateScanProfile(c *gin.Context) {
	var req complianceSvc.CreateScanProfileInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Profiles.Create(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
func (h *ComplianceHandler) UpdateScanProfile(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "profileId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req complianceSvc.UpdateScanProfileInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Profiles.Update(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) ExecuteScanProfile(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "profileId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req complianceSvc.ExecuteScanInput
	_ = c.ShouldBindJSON(&req)
	item, err := h.svc.Scans.ExecuteNow(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}
func (h *ComplianceHandler) ListScans(c *gin.Context) {
	workspaceID, _ := parseComplianceOptionalQueryUint64(c, "workspaceId")
	projectID, _ := parseComplianceOptionalQueryUint64(c, "projectId")
	profileID, _ := parseComplianceOptionalQueryUint64(c, "profileId")
	items, err := h.svc.Scans.List(c.Request.Context(), c.GetUint64(middleware.UserIDKey), complianceSvc.ScanExecutionListFilter{WorkspaceID: derefUint64Value(workspaceID), ProjectID: derefUint64Value(projectID), ProfileID: derefUint64Value(profileID), Status: c.Query("status"), TriggerSource: c.Query("triggerSource")})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) GetScan(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "scanId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.Scans.Get(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) ListFindings(c *gin.Context) {
	workspaceID, _ := parseComplianceOptionalQueryUint64(c, "workspaceId")
	projectID, _ := parseComplianceOptionalQueryUint64(c, "projectId")
	scanExecutionID, _ := parseComplianceOptionalQueryUint64(c, "scanExecutionId")
	if scanExecutionID == nil {
		if raw := strings.TrimSpace(c.Param("scanId")); raw != "" {
			if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
				scanExecutionID = &parsed
			}
		}
	}
	items, err := h.svc.Findings.List(c.Request.Context(), c.GetUint64(middleware.UserIDKey), complianceSvc.FindingListFilter{WorkspaceID: derefUint64Value(workspaceID), ProjectID: derefUint64Value(projectID), ScanExecutionID: derefUint64Value(scanExecutionID), Result: c.Query("result"), RiskLevel: c.Query("riskLevel")})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ComplianceHandler) GetFinding(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "findingId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	finding, err := h.svc.Findings.Get(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	evidences, err := h.svc.Evidence.ListByFinding(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	remediations, _ := h.remediation.ListTasks(c.Request.Context(), complianceSvc.RemediationTaskFilter{})
	exceptions, _ := h.exceptions.ListExceptions(c.Request.Context(), complianceSvc.ExceptionFilter{})
	rechecks, _ := h.rechecks.ListTasks(c.Request.Context(), complianceSvc.RecheckFilter{})
	findingIDText := strconv.FormatUint(id, 10)
	filteredRemediations := make([]complianceSvc.RemediationTask, 0)
	for _, item := range remediations {
		if item.FindingID == findingIDText {
			filteredRemediations = append(filteredRemediations, item)
		}
	}
	filteredExceptions := make([]complianceSvc.ComplianceExceptionRequest, 0)
	for _, item := range exceptions {
		if item.FindingID == findingIDText {
			filteredExceptions = append(filteredExceptions, item)
		}
	}
	filteredRechecks := make([]complianceSvc.RecheckTask, 0)
	for _, item := range rechecks {
		if item.FindingID == findingIDText {
			filteredRechecks = append(filteredRechecks, item)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"id":                finding.ID,
		"scanExecutionId":   finding.ScanExecutionID,
		"controlId":         finding.ControlID,
		"controlTitle":      finding.ControlTitle,
		"result":            finding.Result,
		"riskLevel":         finding.RiskLevel,
		"clusterId":         finding.ClusterID,
		"namespace":         finding.Namespace,
		"resourceKind":      finding.ResourceKind,
		"resourceName":      finding.ResourceName,
		"remediationStatus": finding.RemediationStatus,
		"summary":           finding.Summary,
		"evidences":         evidences,
		"remediationTasks":  filteredRemediations,
		"exceptions":        filteredExceptions,
		"rechecks":          filteredRechecks,
	})
}

func (h *ComplianceHandler) GetRecheck(c *gin.Context) {
	recheckID := strings.TrimSpace(c.Param("recheckId"))
	item, err := h.rechecks.GetTask(c.Request.Context(), recheckID)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ComplianceHandler) GetArchiveExport(c *gin.Context) {
	exportID := strings.TrimSpace(c.Param("exportId"))
	item, err := h.exports.GetExport(c.Request.Context(), exportID)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ComplianceHandler) ListEvidence(c *gin.Context) {
	id, err := parseComplianceParamUint64(c, "findingId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.Evidence.ListByFinding(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) ListRemediationTasks(c *gin.Context) {
	items, err := h.remediation.ListTasks(c.Request.Context(), complianceSvc.RemediationTaskFilter{WorkspaceID: queryUint64Value(c, "workspaceId"), ProjectID: queryUint64Value(c, "projectId"), Owner: strings.TrimSpace(c.Query("owner")), Status: strings.TrimSpace(c.Query("status")), Priority: strings.TrimSpace(c.Query("priority"))})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) CreateRemediationTask(c *gin.Context) {
	findingID := strings.TrimSpace(c.Param("findingId"))
	var req complianceSvc.CreateRemediationTaskInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.remediation.CreateTask(c.Request.Context(), c.GetUint64(middleware.UserIDKey), findingID, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
func (h *ComplianceHandler) UpdateRemediationTask(c *gin.Context) {
	taskID := strings.TrimSpace(c.Param("taskId"))
	var req complianceSvc.UpdateRemediationTaskInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.remediation.UpdateTask(c.Request.Context(), c.GetUint64(middleware.UserIDKey), taskID, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) ListExceptions(c *gin.Context) {
	items, err := h.exceptions.ListExceptions(c.Request.Context(), complianceSvc.ExceptionFilter{WorkspaceID: queryUint64Value(c, "workspaceId"), ProjectID: queryUint64Value(c, "projectId"), Status: strings.TrimSpace(c.Query("status")), BaselineID: strings.TrimSpace(c.Query("baselineId"))})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) CreateException(c *gin.Context) {
	findingID := strings.TrimSpace(c.Param("findingId"))
	var req complianceSvc.CreateExceptionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.exceptions.CreateException(c.Request.Context(), c.GetUint64(middleware.UserIDKey), findingID, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
func (h *ComplianceHandler) ReviewException(c *gin.Context) {
	exceptionID := strings.TrimSpace(c.Param("exceptionId"))
	var req complianceSvc.ReviewExceptionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.exceptions.ReviewException(c.Request.Context(), c.GetUint64(middleware.UserIDKey), exceptionID, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) ListRechecks(c *gin.Context) {
	items, err := h.rechecks.ListTasks(c.Request.Context(), complianceSvc.RecheckFilter{WorkspaceID: queryUint64Value(c, "workspaceId"), ProjectID: queryUint64Value(c, "projectId"), Status: strings.TrimSpace(c.Query("status")), TriggerSource: strings.TrimSpace(c.Query("triggerSource"))})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) CreateRecheck(c *gin.Context) {
	findingID := strings.TrimSpace(c.Param("findingId"))
	var req complianceSvc.CreateRecheckInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.rechecks.CreateTask(c.Request.Context(), c.GetUint64(middleware.UserIDKey), findingID, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
func (h *ComplianceHandler) CompleteRecheck(c *gin.Context) {
	recheckID := strings.TrimSpace(c.Param("recheckId"))
	var req complianceSvc.CompleteRecheckInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.rechecks.CompleteTask(c.Request.Context(), c.GetUint64(middleware.UserIDKey), recheckID, req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) GetOverview(c *gin.Context) {
	item, err := h.overview.GetOverview(c.Request.Context(), complianceSvc.OverviewFilter{WorkspaceID: queryUint64Value(c, "workspaceId"), ProjectID: queryUint64Value(c, "projectId")})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *ComplianceHandler) ListTrends(c *gin.Context) {
	items, err := h.trends.GetTrends(c.Request.Context(), complianceSvc.TrendFilter{WorkspaceID: queryUint64Value(c, "workspaceId"), ProjectID: queryUint64Value(c, "projectId")})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ComplianceHandler) ListArchiveExports(c *gin.Context) {
	items, err := h.exports.ListExports(c.Request.Context(), complianceSvc.ArchiveExportFilter{WorkspaceID: queryUint64Value(c, "workspaceId"), ProjectID: queryUint64Value(c, "projectId"), ExportScope: strings.TrimSpace(c.Query("scope")), Status: strings.TrimSpace(c.Query("status"))})
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
func (h *ComplianceHandler) CreateArchiveExport(c *gin.Context) {
	var req complianceSvc.CreateArchiveExportInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.exports.CreateExport(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeComplianceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func writeComplianceError(c *gin.Context, err error) {
	switch {
	case err == nil:
		return
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	case strings.Contains(strings.ToLower(err.Error()), "denied"):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func parseComplianceParamUint64(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(c.Param(key)), 10, 64)
}
func parseComplianceOptionalQueryUint64(c *gin.Context, key string) (*uint64, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func queryUint64Value(c *gin.Context, key string) uint64 {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0
	}
	v, _ := strconv.ParseUint(raw, 10, 64)
	return v
}
func derefUint64Value(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}
