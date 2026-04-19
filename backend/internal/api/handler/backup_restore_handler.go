package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/service/backuprestore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BackupRestoreHandler struct {
	svc *backuprestore.Service
}

func NewBackupRestoreHandler(svc *backuprestore.Service) *BackupRestoreHandler {
	return &BackupRestoreHandler{svc: svc}
}

func (h *BackupRestoreHandler) ListPolicies(c *gin.Context) {
	items, err := h.svc.ListPolicies(c.Request.Context(), c.GetUint64(middleware.UserIDKey), backuprestore.PolicyListFilter{
		ScopeType: strings.TrimSpace(c.Query("scopeType")),
		Status:    strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BackupRestoreHandler) CreatePolicy(c *gin.Context) {
	var req backuprestore.CreateBackupPolicyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreatePolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *BackupRestoreHandler) RunPolicy(c *gin.Context) {
	policyID, ok := parseUint64Param(c, "policyId")
	if !ok {
		return
	}
	item, err := h.svc.RunPolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), policyID)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}

func (h *BackupRestoreHandler) ListRestorePoints(c *gin.Context) {
	policyID, _ := parseOptionalUint64(c.Query("policyId"))
	items, err := h.svc.ListRestorePoints(c.Request.Context(), c.GetUint64(middleware.UserIDKey), backuprestore.RestorePointListFilter{
		PolicyID: policyID,
		Result:   strings.TrimSpace(c.Query("result")),
	})
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BackupRestoreHandler) GetRestorePoint(c *gin.Context) {
	id, ok := parseUint64Param(c, "restorePointId")
	if !ok {
		return
	}
	item, err := h.svc.GetRestorePoint(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *BackupRestoreHandler) CreateRestoreJob(c *gin.Context) {
	var req backuprestore.CreateRestoreJobInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateRestoreJob(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}

func (h *BackupRestoreHandler) ListRestoreJobs(c *gin.Context) {
	items, err := h.svc.ListRestoreJobs(c.Request.Context(), c.GetUint64(middleware.UserIDKey), backuprestore.RestoreJobListFilter{
		JobType: strings.TrimSpace(c.Query("jobType")),
		Status:  strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BackupRestoreHandler) ValidateRestoreJob(c *gin.Context) {
	jobID, ok := parseUint64Param(c, "jobId")
	if !ok {
		return
	}
	item, err := h.svc.ValidateRestoreJobByID(c.Request.Context(), c.GetUint64(middleware.UserIDKey), jobID)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *BackupRestoreHandler) CreateMigrationPlan(c *gin.Context) {
	var req backuprestore.CreateMigrationPlanInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateMigrationPlan(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *BackupRestoreHandler) ListDrillPlans(c *gin.Context) {
	items, err := h.svc.ListDrillPlans(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BackupRestoreHandler) CreateDrillPlan(c *gin.Context) {
	var req backuprestore.CreateDRDrillPlanInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateDrillPlan(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *BackupRestoreHandler) RunDrillPlan(c *gin.Context) {
	planID, ok := parseUint64Param(c, "planId")
	if !ok {
		return
	}
	item, err := h.svc.RunDrillPlan(c.Request.Context(), c.GetUint64(middleware.UserIDKey), planID)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}

func (h *BackupRestoreHandler) GetDrillRecord(c *gin.Context) {
	recordID, ok := parseUint64Param(c, "recordId")
	if !ok {
		return
	}
	item, err := h.svc.GetDrillRecord(c.Request.Context(), c.GetUint64(middleware.UserIDKey), recordID)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *BackupRestoreHandler) GenerateDrillReport(c *gin.Context) {
	recordID, ok := parseUint64Param(c, "recordId")
	if !ok {
		return
	}
	item, err := h.svc.GenerateDrillReport(c.Request.Context(), c.GetUint64(middleware.UserIDKey), recordID)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *BackupRestoreHandler) ListAuditEvents(c *gin.Context) {
	items, err := h.svc.ListAuditEvents(
		c.Request.Context(),
		c.GetUint64(middleware.UserIDKey),
		strings.TrimSpace(c.Query("action")),
		strings.TrimSpace(c.Query("outcome")),
		strings.TrimSpace(c.Query("targetType")),
	)
	if err != nil {
		writeBackupRestoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func writeBackupRestoreError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, backuprestore.ErrBackupRestoreScopeDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, backuprestore.ErrBackupRestoreConflict), errors.Is(err, backuprestore.ErrBackupRestoreBlocked):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, backuprestore.ErrBackupRestoreInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func parseOptionalUint64(v string) (uint64, bool) {
	if strings.TrimSpace(v) == "" {
		return 0, false
	}
	out, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, false
	}
	return out, true
}
