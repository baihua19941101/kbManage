package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

type clusterLifecycleIntegrationCtx struct {
	app    *testutil.App
	token  string
	access testutil.ObservabilityAccessSeed
	userID uint64
}

func newClusterLifecycleIntegrationCtx(t *testing.T, roleKey string) clusterLifecycleIntegrationCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "cl-int-" + roleKey, Password: "Integration@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "cluster-lifecycle-int", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	return clusterLifecycleIntegrationCtx{app: app, token: token, access: access, userID: user.User.ID}
}

func performClusterLifecycleIntegrationRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func createDriverForIntegration(t *testing.T, ctx clusterLifecycleIntegrationCtx) {
	t.Helper()
	resp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/drivers", `{
		"driverKey":"generic-driver",
		"version":"v1",
		"displayName":"Generic Driver",
		"providerType":"generic",
		"status":"active",
		"capabilityProfileVersion":"v1",
		"schemaVersion":"v1",
		"capabilities":[{"capabilityDomain":"network","supportLevel":"native","compatibilityStatus":"compatible"}]
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create driver failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}

func createClusterForIntegration(t *testing.T, ctx clusterLifecycleIntegrationCtx) {
	t.Helper()
	createDriverForIntegration(t, ctx)
	resp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters", `{
		"name":"cluster-int",
		"displayName":"Cluster Int",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.access.ProjectID, 10)+`,
		"infrastructureType":"generic",
		"driverKey":"generic-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"nodePools":[{"name":"workers","role":"worker","desiredCount":3,"minCount":1,"maxCount":5,"version":"v1.30.1","zoneRefs":["zone-a"]}]
	}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("create cluster failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}

func seedRunningLifecycleOperation(t *testing.T, ctx clusterLifecycleIntegrationCtx, clusterID uint64) {
	t.Helper()
	if err := ctx.app.DB.WithContext(context.Background()).Create(&domain.LifecycleOperation{
		ClusterID:     &clusterID,
		OperationType: domain.LifecycleOperationUpgrade,
		Status:        domain.LifecycleOperationRunning,
		RiskLevel:     domain.LifecycleRiskHigh,
		RequestedBy:   ctx.userID,
	}).Error; err != nil {
		t.Fatalf("seed running lifecycle operation failed: %v", err)
	}
}
