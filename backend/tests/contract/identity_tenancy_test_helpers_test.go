package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/api/router"
	authSvc "kbmanage/backend/internal/service/auth"
	"kbmanage/backend/tests/testutil"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type identityTenancyContractCtx struct {
	Router *gin.Engine
	DB     *gorm.DB
	Token  string
	Access testutil.ObservabilityAccessSeed
	UserID uint64
}

func newIdentityTenancyContractCtx(t *testing.T, roleKey string) *identityTenancyContractCtx {
	t.Helper()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "identity-contract-" + strings.ReplaceAll(roleKey, "_", "-"),
		Password: "Contract@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "identity-contract", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	tokenSvc := authSvc.NewTokenService(app.Config.JWTSecret, app.Config.AccessTokenTTL, app.Config.RefreshTokenTTL)

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	v1 := engine.Group("/api/v1")
	authed := v1.Group("/")
	authed.Use(middleware.AuthRequired(tokenSvc))
	router.RegisterIdentityTenancyRoutes(authed, app.DB, nil)

	return &identityTenancyContractCtx{
		Router: engine,
		DB:     app.DB,
		Token:  token,
		Access: access,
		UserID: user.User.ID,
	}
}

func performIdentityTenancyContractRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func mustReadIdentityContractID(t *testing.T, body []byte, key string) uint64 {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v body=%s", err, strings.TrimSpace(string(body)))
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

func createIdentitySourceContract(t *testing.T, ctx *identityTenancyContractCtx, sourceType, loginMode string) uint64 {
	t.Helper()
	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/sources", `{
		"name":"`+sourceType+`-source",
		"sourceType":"`+sourceType+`",
		"loginMode":"`+loginMode+`",
		"scopeMode":"platform"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create identity source failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadIdentityContractID(t, resp.Body.Bytes(), "id")
}

func createOrganizationContract(t *testing.T, ctx *identityTenancyContractCtx, name string, parentUnitID uint64) uint64 {
	t.Helper()
	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/organizations", `{
		"unitType":"organization",
		"name":"`+name+`",
		"parentUnitId":`+strconv.FormatUint(parentUnitID, 10)+`
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create organization failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadIdentityContractID(t, resp.Body.Bytes(), "id")
}

func createRoleDefinitionContract(t *testing.T, ctx *identityTenancyContractCtx, name, roleLevel, inheritancePolicy string, delegable bool) uint64 {
	t.Helper()
	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/roles", `{
		"name":"`+name+`",
		"roleLevel":"`+roleLevel+`",
		"permissionSummary":"read,write",
		"inheritancePolicy":"`+inheritancePolicy+`",
		"delegable":`+strconv.FormatBool(delegable)+`
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create role failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadIdentityContractID(t, resp.Body.Bytes(), "id")
}

func createDelegationGrantContract(t *testing.T, ctx *identityTenancyContractCtx, delegateRef string) uint64 {
	t.Helper()
	validFrom := time.Now().UTC()
	validUntil := validFrom.Add(2 * time.Hour)
	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/delegations", `{
		"grantorRef":"`+strconv.FormatUint(ctx.UserID, 10)+`",
		"delegateRef":"`+delegateRef+`",
		"allowedRoleLevels":["project","workspace"],
		"validFrom":"`+validFrom.Format(time.RFC3339)+`",
		"validUntil":"`+validUntil.Format(time.RFC3339)+`",
		"reason":"contract delegation"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create delegation failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadIdentityContractID(t, resp.Body.Bytes(), "id")
}
