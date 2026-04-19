package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"

	"gorm.io/gorm"
)

type sreIntegrationCtx struct {
	App    *testutil.App
	Token  string
	DB     *gorm.DB
	UserID uint64
	Access testutil.ObservabilityAccessSeed
}

func newSREIntegrationCtx(t *testing.T, roleKey string) *sreIntegrationCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "sre-integration-" + strings.ReplaceAll(roleKey, "_", "-"),
		Password: "Integration@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "sre-integration", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	return &sreIntegrationCtx{
		App:    app,
		Token:  token,
		DB:     app.DB,
		UserID: user.User.ID,
		Access: access,
	}
}

func performSREIntegrationRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func seedSREIntegrationEvidence(t *testing.T, ctx *sreIntegrationCtx) {
	t.Helper()
	now := time.Now()
	projectID := ctx.Access.ProjectID
	baseline := &domain.CapacityBaseline{
		WorkspaceID:       ctx.Access.WorkspaceID,
		ProjectID:         &projectID,
		Name:              "integration-capacity",
		ResourceDimension: "clusters",
		ForecastWindow:    "7d",
		ForecastResult:    "未来 7 天增长 20%",
		ConfidenceLevel:   "medium",
		Status:            domain.CapacityBaselineStatusWarning,
		OwnerUserID:       ctx.UserID,
	}
	if err := ctx.DB.WithContext(context.Background()).Create(baseline).Error; err != nil {
		t.Fatalf("seed capacity baseline failed: %v", err)
	}
	runbook := &domain.RunbookArticle{
		WorkspaceID:         ctx.Access.WorkspaceID,
		ProjectID:           &projectID,
		Title:               "integration-runbook",
		ScenarioType:        "capacity",
		RiskLevel:           "medium",
		ChecklistSummary:    "检查容量",
		RecoverySteps:       "扩容",
		VerificationSummary: "观察状态",
		Status:              domain.RunbookStatusActive,
		OwnerUserID:         ctx.UserID,
	}
	if err := ctx.DB.WithContext(context.Background()).Create(runbook).Error; err != nil {
		t.Fatalf("seed runbook failed: %v", err)
	}
	evidence := &domain.ScaleEvidence{
		WorkspaceID:        ctx.Access.WorkspaceID,
		ProjectID:          &projectID,
		CapacityBaselineID: &baseline.ID,
		RunbookArticleID:   &runbook.ID,
		EvidenceType:       "loadtest",
		Scope:              "platform",
		SampleWindow:       "24h",
		Summary:            "压测达到阈值",
		BottleneckSummary:  "api-server",
		ForecastSummary:    "未来 7 天容量不足",
		ConfidenceLevel:    "medium",
		Status:             domain.ScaleEvidenceStatusAnalyzed,
		CapturedAt:         now,
		OwnerUserID:        ctx.UserID,
	}
	if err := ctx.DB.WithContext(context.Background()).Create(evidence).Error; err != nil {
		t.Fatalf("seed evidence failed: %v", err)
	}
}
