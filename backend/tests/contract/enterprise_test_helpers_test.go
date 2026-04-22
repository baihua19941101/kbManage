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

type enterpriseContractCtx struct {
	App    *testutil.App
	Token  string
	DB     *gorm.DB
	UserID uint64
	Access testutil.ObservabilityAccessSeed
}

func newEnterpriseContractCtx(t *testing.T, roleKey string) *enterpriseContractCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "enterprise-contract-" + strings.ReplaceAll(roleKey, "_", "-"), Password: "Contract@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "enterprise-contract", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	return &enterpriseContractCtx{App: app, Token: token, DB: app.DB, UserID: user.User.ID, Access: access}
}

func performEnterpriseContractRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func mustDecodeEnterpriseObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustDecodeEnterpriseItems(t *testing.T, body []byte) []any {
	t.Helper()
	payload := mustDecodeEnterpriseObject(t, body)
	raw, ok := payload["items"]
	if !ok {
		t.Fatalf("missing items: %v", payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("items type=%T", raw)
	}
	return items
}

func seedEnterpriseAuditData(t *testing.T, ctx *enterpriseContractCtx) uint64 {
	t.Helper()
	projectID := ctx.Access.ProjectID
	now := time.Now()
	trail := &domain.PermissionChangeTrail{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, SubjectType: "user", SubjectRef: "alice", ChangeType: "grant", BeforeState: "viewer", AfterState: "editor", AuthorizationBasis: "ticket-1", ScopeType: "workspace", ScopeRef: "ws-1", EvidenceCompleteness: "complete", ChangedAt: now, ChangedBy: ctx.UserID}
	if err := ctx.DB.WithContext(context.Background()).Create(trail).Error; err != nil {
		t.Fatalf("seed trail: %v", err)
	}
	op := &domain.KeyOperationTrace{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, ActorType: "user", ActorRef: "alice", OperationType: "role-bind", TargetType: "workspace", TargetRef: "ws-1", ContextSummary: "critical", RiskLevel: "high", Outcome: "success", OccurredAt: now, RelatedTrailID: &trail.ID}
	if err := ctx.DB.WithContext(context.Background()).Create(op).Error; err != nil {
		t.Fatalf("seed operation: %v", err)
	}
	coverage := &domain.GovernanceCoverageSnapshot{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, SnapshotAt: now, CoverageDomain: "rbac", CoverageRate: 87.5, StatusBreakdown: "已达标:7,未达标:1", MissingReasonSummary: "1 个项目未接入", ConfidenceLevel: "medium", TrendSummary: "稳定", Owner: "audit-team"}
	if err := ctx.DB.WithContext(context.Background()).Create(coverage).Error; err != nil {
		t.Fatalf("seed coverage: %v", err)
	}
	action := &domain.GovernanceActionItem{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, SourceType: "risk-event", SourceRef: "risk-1", Title: "复核高风险授权", Priority: "high", Owner: "audit-team", Status: domain.GovernanceActionStatusOpen}
	if err := ctx.DB.WithContext(context.Background()).Create(action).Error; err != nil {
		t.Fatalf("seed action: %v", err)
	}
	artifact := &domain.DeliveryArtifact{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, ArtifactType: "install-guide", Title: "安装手册", VersionScope: "v1.0", EnvironmentScope: "prod", OwnerRole: "delivery", ApplicabilityNote: "标准版", Status: domain.DeliveryArtifactStatusActive}
	if err := ctx.DB.WithContext(context.Background()).Create(artifact).Error; err != nil {
		t.Fatalf("seed artifact: %v", err)
	}
	bundle := &domain.DeliveryReadinessBundle{WorkspaceID: ctx.Access.WorkspaceID, ProjectID: &projectID, Name: "客户A交付包", TargetEnvironment: "prod", TargetAudience: "customer", ArtifactSummary: "含安装与升级材料", ChecklistStatus: "partial", MissingItems: "验收签字", ReadinessConclusion: domain.DeliveryReadinessConditionally}
	if err := ctx.DB.WithContext(context.Background()).Create(bundle).Error; err != nil {
		t.Fatalf("seed bundle: %v", err)
	}
	check := &domain.DeliveryChecklistItem{BundleID: bundle.ID, CheckItem: "确认安装文档", Category: "doc", Owner: "delivery", EvidenceRequirement: "pdf", Status: "done"}
	if err := ctx.DB.WithContext(context.Background()).Create(check).Error; err != nil {
		t.Fatalf("seed checklist: %v", err)
	}
	return bundle.ID
}
