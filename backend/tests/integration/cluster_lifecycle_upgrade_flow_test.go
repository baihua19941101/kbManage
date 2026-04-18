package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestClusterLifecycleIntegration_CreateUpgradeRetireFlow(t *testing.T) {
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "cl-int", Password: "Integration@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "cluster-lifecycle-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	driverResp := performClusterLifecycleIntegrationRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/drivers", `{
		"driverKey":"generic-driver",
		"version":"v1",
		"displayName":"Generic Driver",
		"providerType":"generic",
		"status":"active",
		"capabilityProfileVersion":"v1",
		"schemaVersion":"v1"
	}`)
	if driverResp.Code != http.StatusCreated {
		t.Fatalf("create driver failed status=%d body=%s", driverResp.Code, strings.TrimSpace(driverResp.Body.String()))
	}

	createResp := performClusterLifecycleIntegrationRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters", `{
		"name":"provisioned-int",
		"displayName":"Provisioned Int",
		"workspaceId":`+strconv.FormatUint(access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(access.ProjectID, 10)+`,
		"infrastructureType":"generic",
		"driverKey":"generic-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"nodePools":[{"name":"workers","role":"worker","desiredCount":3,"minCount":1,"maxCount":5,"version":"v1.30.1","zoneRefs":["zone-a"]}]
	}`)
	if createResp.Code != http.StatusAccepted {
		t.Fatalf("create cluster failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	planResp := performClusterLifecycleIntegrationRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/upgrade-plans", `{
		"targetVersion":"v1.31.0",
		"impactSummary":"minor upgrade"
	}`)
	if planResp.Code != http.StatusCreated {
		t.Fatalf("create upgrade plan failed status=%d body=%s", planResp.Code, strings.TrimSpace(planResp.Body.String()))
	}

	executeResp := performClusterLifecycleIntegrationRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/upgrade-plans/1/execute", "")
	if executeResp.Code != http.StatusAccepted {
		t.Fatalf("execute upgrade plan failed status=%d body=%s", executeResp.Code, strings.TrimSpace(executeResp.Body.String()))
	}

	scaleResp := performClusterLifecycleIntegrationRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/node-pools/1/scale", `{"desiredCount":4}`)
	if scaleResp.Code != http.StatusAccepted {
		t.Fatalf("scale node pool failed status=%d body=%s", scaleResp.Code, strings.TrimSpace(scaleResp.Body.String()))
	}

	retireResp := performClusterLifecycleIntegrationRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/retire", `{
		"reason":"retire integration cluster",
		"confirmationScope":"full",
		"conclusion":"done"
	}`)
	if retireResp.Code != http.StatusAccepted {
		t.Fatalf("retire cluster failed status=%d body=%s", retireResp.Code, strings.TrimSpace(retireResp.Body.String()))
	}
}
