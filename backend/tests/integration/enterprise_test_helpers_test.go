package integration_test

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

type enterpriseIntegrationCtx struct {
	App    *testutil.App
	Token  string
	DB     *gorm.DB
	UserID uint64
	Access testutil.ObservabilityAccessSeed
}

func newEnterpriseIntegrationCtx(t *testing.T, roleKey string) *enterpriseIntegrationCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "enterprise-integration-" + strings.ReplaceAll(roleKey, "_", "-"), Password: "Integration@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "enterprise-integration", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	return &enterpriseIntegrationCtx{App: app, Token: token, DB: app.DB, UserID: user.User.ID, Access: access}
}

func performEnterpriseIntegrationRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func mustDecodeEnterpriseIntegrationObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func seedEnterpriseIntegrationData(t *testing.T, ctx *enterpriseIntegrationCtx) uint64 {
	t.Helper()
	projectID := ctx.Access.ProjectID
	now := time.Now()
	records := []any{
		&domain.PermissionChangeTrail{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, SubjectType: "group", SubjectRef: "ops", ChangeType: "delegate", BeforeState: "viewer", AfterState: "admin", AuthorizationBasis: "change-req", ScopeType: "workspace", ScopeRef: "ws-1", EvidenceCompleteness: "complete", ChangedAt: now, ChangedBy: ctx.UserID},
		&domain.KeyOperationTrace{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, ActorType: "group", ActorRef: "ops", OperationType: "export", TargetType: "report", TargetRef: "r-1", ContextSummary: "weekly", RiskLevel: "medium", Outcome: "success", OccurredAt: now},
		&domain.GovernanceCoverageSnapshot{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, SnapshotAt: now, CoverageDomain: "identity", CoverageRate: 91, StatusBreakdown: "已达标:10", MissingReasonSummary: "0", ConfidenceLevel: "high", TrendSummary: "上升", Owner: "audit"},
		&domain.GovernanceActionItem{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, SourceType: "coverage", SourceRef: "c-1", Title: "修复缺口", Priority: "medium", Owner: "audit", Status: domain.GovernanceActionStatusOpen},
		&domain.DeliveryArtifact{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, ArtifactType: "ops-runbook", Title: "运维手册", VersionScope: "v1.0", EnvironmentScope: "prod", OwnerRole: "ops", ApplicabilityNote: "标准", Status: domain.DeliveryArtifactStatusActive},
	}
	for _, item := range records {
		if err := ctx.DB.WithContext(context.Background()).Create(item).Error; err != nil {
			t.Fatalf("seed record: %v", err)
		}
	}
	bundle := &domain.DeliveryReadinessBundle{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, Name: "标准交付包", TargetEnvironment: "prod", TargetAudience: "customer", ArtifactSummary: "安装+运维", ChecklistStatus: "ready", MissingItems: "", ReadinessConclusion: domain.DeliveryReadinessReady}
	if err := ctx.DB.WithContext(context.Background()).Create(bundle).Error; err != nil {
		t.Fatalf("seed bundle: %v", err)
	}
	check := &domain.DeliveryChecklistItem{BundleID: bundle.ID, CheckItem: "验收项", Category: "acceptance", Owner: "delivery", EvidenceRequirement: "截图", Status: "done"}
	if err := ctx.DB.WithContext(context.Background()).Create(check).Error; err != nil {
		t.Fatalf("seed checklist: %v", err)
	}
	return bundle.ID
}
