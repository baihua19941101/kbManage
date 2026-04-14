package contract_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	"kbmanage/backend/tests/testutil"
)

func TestGitOpsContract_AccessControlDeniedWithoutScopeBinding(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-access-denied-contract-user",
		Password: "GitOps@123",
	})
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	seed := seedGitOpsScopeFixture(t, app, user.User.ID, "workspace-owner")

	if err := app.DB.WithContext(context.Background()).Where("subject_type = ? AND subject_id = ?", "user", user.User.ID).Delete(&repository.ScopeRoleBinding{}).Error; err != nil {
		t.Fatalf("delete scope binding failed: %v", err)
	}

	cases := []struct {
		name   string
		method string
		target string
		body   string
	}{
		{name: "list sources", method: http.MethodGet, target: "/api/v1/gitops/sources?workspaceId=" + u64s(seed.WorkspaceID)},
		{name: "create source", method: http.MethodPost, target: "/api/v1/gitops/sources", body: `{"name":"demo","sourceType":"git","endpoint":"https://git.example.com/demo.git","workspaceId":` + u64s(seed.WorkspaceID) + `,"projectId":` + u64s(seed.ProjectID) + `}`},
		{name: "list target groups", method: http.MethodGet, target: "/api/v1/gitops/target-groups?workspaceId=" + u64s(seed.WorkspaceID)},
		{name: "list delivery units", method: http.MethodGet, target: "/api/v1/gitops/delivery-units?workspaceId=" + u64s(seed.WorkspaceID)},
		{name: "submit action", method: http.MethodPost, target: "/api/v1/gitops/delivery-units/" + u64s(seed.UnitID) + "/actions", body: `{"actionType":"sync","payload":{"reason":"contract-test"}}`},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := performObservabilityAuthedRequest(t, app.Router, token, tc.method, tc.target, tc.body)
			if resp.Code != http.StatusForbidden {
				t.Fatalf("expected status=403, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
		})
	}
}

func seedGitOpsScopeFixture(t *testing.T, app *testutil.App, userID uint64, roleKey string) gitOpsScopeSeed {
	t.Helper()

	seed := testutil.SeedObservabilityAccess(t, app.DB, userID, "gitops-access", roleKey)

	source := &domain.DeliverySource{
		Name:        "gitops-contract-source",
		SourceType:  domain.DeliverySourceTypeGit,
		Endpoint:    "https://git.example.com/contract.git",
		WorkspaceID: &seed.WorkspaceID,
		ProjectID:   &seed.ProjectID,
		Status:      domain.DeliverySourceStatusReady,
	}
	if err := app.DB.WithContext(context.Background()).Create(source).Error; err != nil {
		t.Fatalf("seed source failed: %v", err)
	}

	unit := &domain.ApplicationDeliveryUnit{
		Name:           "gitops-contract-unit",
		WorkspaceID:    seed.WorkspaceID,
		ProjectID:      &seed.ProjectID,
		SourceID:       source.ID,
		SourcePath:     "apps/demo",
		SyncMode:       domain.DeliverySyncModeManual,
		DeliveryStatus: domain.DeliveryUnitStatusReady,
	}
	if err := app.DB.WithContext(context.Background()).Create(unit).Error; err != nil {
		t.Fatalf("seed delivery unit failed: %v", err)
	}

	return gitOpsScopeSeed{
		WorkspaceID: seed.WorkspaceID,
		ProjectID:   seed.ProjectID,
		SourceID:    source.ID,
		UnitID:      unit.ID,
	}
}

type gitOpsScopeSeed struct {
	WorkspaceID uint64
	ProjectID   uint64
	SourceID    uint64
	UnitID      uint64
}

func u64s(v uint64) string {
	return strconv.FormatUint(v, 10)
}
