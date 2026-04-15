package contract_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

func TestSecurityPolicyAuditContract_QuerySecurityPolicyEvents(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-audit-contract",
		Password: "SecurityPolicyAudit@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-audit-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyBody := fmt.Sprintf(`{
		"name":"audit-policy",
		"workspaceId":%d,
		"projectId":%d,
		"scopeLevel":"project",
		"category":"admission",
		"ruleTemplate":{"denyPrivileged":true},
		"defaultEnforcementMode":"warn",
		"riskLevel":"high"
	}`, access.WorkspaceID, access.ProjectID)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/security-policies", strings.NewReader(policyBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create policy status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	actorID := seeded.User.ID
	if err := app.DB.Create(&domain.AuditEvent{
		RequestID:     "non-security-policy-event",
		ActorID:       &actorID,
		Action:        "workloadops.action.submit",
		ResourceType:  "workloadops",
		ResourceID:    "operation:100",
		Outcome:       domain.AuditOutcomeSuccess,
		AuditCategory: "workloadops",
		ActionScope:   "action",
		CreatedAt:     time.Now().UTC(),
	}).Error; err != nil {
		t.Fatalf("seed non-security-policy audit event failed: %v", err)
	}

	eventsReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/audit/security-policies/events?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&action=securitypolicy.policy.create",
		nil,
	)
	eventsReq.Header.Set("Authorization", "Bearer "+token)
	eventsResp := httptest.NewRecorder()
	app.Router.ServeHTTP(eventsResp, eventsReq)
	if eventsResp.Code != http.StatusOK {
		t.Fatalf("expected security policy audit query status=200, got=%d body=%s", eventsResp.Code, strings.TrimSpace(eventsResp.Body.String()))
	}
	payload := mustDecodeSecurityPolicyContractObject(t, eventsResp.Body.Bytes())
	items := mustReadSecurityPolicyContractArray(t, payload, "items")
	if len(items) == 0 {
		t.Fatalf("expected at least one security policy event payload=%v", payload)
	}
	for i := range items {
		item, ok := items[i].(map[string]any)
		if !ok {
			t.Fatalf("unexpected event item type=%T", items[i])
		}
		action := strings.TrimSpace(mustReadSecurityPolicyContractString(t, item, "Action"))
		if !strings.HasPrefix(action, "securitypolicy.") {
			t.Fatalf("expected only securitypolicy.* action, got action=%q item=%v", action, item)
		}
	}
	count := int(payload["count"].(float64))
	if count != len(items) {
		t.Fatalf("expected count equals items length, count=%d len=%d payload=%v", count, len(items), payload)
	}
}
