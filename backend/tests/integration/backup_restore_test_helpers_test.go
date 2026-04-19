package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/api/router"
	authSvc "kbmanage/backend/internal/service/auth"
	"kbmanage/backend/tests/testutil"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type backupRestoreIntegrationCtx struct {
	Router *gin.Engine
	DB     *gorm.DB
	Token  string
	Access testutil.ObservabilityAccessSeed
}

func newBackupRestoreIntegrationCtx(t *testing.T, roleKey string) *backupRestoreIntegrationCtx {
	t.Helper()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "backup-integration-" + strings.ReplaceAll(roleKey, "_", "-"),
		Password: "Integration@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "backup-integration", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	tokenSvc := authSvc.NewTokenService(app.Config.JWTSecret, app.Config.AccessTokenTTL, app.Config.RefreshTokenTTL)

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	v1 := engine.Group("/api/v1")
	authed := v1.Group("/")
	authed.Use(middleware.AuthRequired(tokenSvc))
	router.RegisterAuditRoutes(authed, app.DB)
	router.RegisterBackupRestoreRoutes(authed, app.DB, nil)

	return &backupRestoreIntegrationCtx{
		Router: engine,
		DB:     app.DB,
		Token:  token,
		Access: access,
	}
}

func performBackupRestoreIntegrationRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	return resp
}

func createBackupRestoreIntegrationPolicy(t *testing.T, ctx *backupRestoreIntegrationCtx, name string) uint64 {
	t.Helper()
	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies", `{
		"name":"`+name+`",
		"scopeType":"namespace",
		"scopeRef":"orders-prod",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"executionMode":"manual",
		"retentionRule":"14d",
		"consistencyLevel":"application-consistent",
		"status":"active"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create policy failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadBackupRestoreIntegrationID(t, resp.Body.Bytes(), "id")
}

func runBackupRestoreIntegrationPolicy(t *testing.T, ctx *backupRestoreIntegrationCtx, policyID uint64) uint64 {
	t.Helper()
	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies/"+strconv.FormatUint(policyID, 10)+"/run", "")
	if resp.Code != http.StatusAccepted {
		t.Fatalf("run policy failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadBackupRestoreIntegrationID(t, resp.Body.Bytes(), "id")
}

func createBackupRestoreIntegrationDrillPlan(t *testing.T, ctx *backupRestoreIntegrationCtx, name string) uint64 {
	t.Helper()
	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans", `{
		"name":"`+name+`",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"scopeSelection":{"namespaces":["orders-prod"]},
		"rpoTargetMinutes":15,
		"rtoTargetMinutes":30,
		"roleAssignments":["sre","biz-owner"],
		"cutoverProcedure":["freeze writes","switch traffic"],
		"validationChecklist":["verify api","verify jobs"]
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create drill plan failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadBackupRestoreIntegrationID(t, resp.Body.Bytes(), "id")
}

func mustReadBackupRestoreIntegrationID(t *testing.T, body []byte, key string) uint64 {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	number, ok := raw.(float64)
	if !ok || number <= 0 {
		t.Fatalf("field %q must be positive number, got=%T value=%v", key, raw, raw)
	}
	return uint64(number)
}
