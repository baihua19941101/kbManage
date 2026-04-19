package contract_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"

	"gorm.io/gorm"
)

type sreContractCtx struct {
	App    *testutil.App
	Token  string
	DB     *gorm.DB
	UserID uint64
	Access testutil.ObservabilityAccessSeed
}

func newSREContractCtx(t *testing.T, roleKey string) *sreContractCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "sre-contract-" + strings.ReplaceAll(roleKey, "_", "-"),
		Password: "Contract@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "sre-contract", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	return &sreContractCtx{
		App:    app,
		Token:  token,
		DB:     app.DB,
		UserID: user.User.ID,
		Access: access,
	}
}

func performSREContractRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func mustDecodeSREObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustDecodeSREItems(t *testing.T, body []byte) []any {
	t.Helper()
	payload := mustDecodeSREObject(t, body)
	raw, ok := payload["items"]
	if !ok {
		t.Fatalf("missing items payload=%v", payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("items not array type=%T", raw)
	}
	return items
}

func seedSRECapacityEvidence(t *testing.T, ctx *sreContractCtx) {
	t.Helper()
	now := time.Now()
	baseline := &domain.CapacityBaseline{
		WorkspaceID:       ctx.Access.WorkspaceID,
		ProjectID:         &ctx.Access.ProjectID,
		Name:              "contract-capacity",
		ResourceDimension: "clusters",
		ForecastWindow:    "7d",
		ForecastResult:    "未来 7 天增长 10%",
		ConfidenceLevel:   "medium",
		Status:            domain.CapacityBaselineStatusActive,
		OwnerUserID:       ctx.UserID,
	}
	if err := ctx.DB.WithContext(context.Background()).Create(baseline).Error; err != nil {
		t.Fatalf("seed capacity baseline failed: %v", err)
	}
	runbook := &domain.RunbookArticle{
		WorkspaceID:         ctx.Access.WorkspaceID,
		ProjectID:           &ctx.Access.ProjectID,
		Title:               "contract-runbook",
		ScenarioType:        "capacity",
		RiskLevel:           "medium",
		ChecklistSummary:    "检查任务积压",
		RecoverySteps:       "扩容控制面",
		VerificationSummary: "观察恢复情况",
		Status:              domain.RunbookStatusActive,
		OwnerUserID:         ctx.UserID,
	}
	if err := ctx.DB.WithContext(context.Background()).Create(runbook).Error; err != nil {
		t.Fatalf("seed runbook failed: %v", err)
	}
	evidence := &domain.ScaleEvidence{
		WorkspaceID:        ctx.Access.WorkspaceID,
		ProjectID:          &ctx.Access.ProjectID,
		CapacityBaselineID: &baseline.ID,
		RunbookArticleID:   &runbook.ID,
		EvidenceType:       "loadtest",
		Scope:              "platform",
		SampleWindow:       "24h",
		Summary:            "压测显示 API 延迟升高",
		BottleneckSummary:  "api-server 并发",
		ForecastSummary:    "7 天后容量紧张",
		ConfidenceLevel:    "medium",
		Status:             domain.ScaleEvidenceStatusAnalyzed,
		CapturedAt:         now,
		OwnerUserID:        ctx.UserID,
	}
	if err := ctx.DB.WithContext(context.Background()).Create(evidence).Error; err != nil {
		t.Fatalf("seed scale evidence failed: %v", err)
	}
}
