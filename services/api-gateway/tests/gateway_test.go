package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/api-gateway/internal/router"
)

func TestGatewayHealthz(t *testing.T) {
	cfg := router.GatewayConfig{
		IdentityServiceURL: "http://localhost:8081",
		JWTSecret:          "gateway-test-secret-32-chars-long!",
	}
	h := router.NewGatewayRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got: %d", rec.Code)
	}

	// Verify CORS header injected
	originHeader := rec.Header().Get("Access-Control-Allow-Origin")
	if originHeader != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got: %s", originHeader)
	}
}

func TestGatewayProtectedFeedWithoutToken(t *testing.T) {
	cfg := router.GatewayConfig{
		IdentityServiceURL: "http://localhost:8081",
		JWTSecret:          "gateway-test-secret-32-chars-long!",
	}
	h := router.NewGatewayRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for unauthenticated feed access, got: %d", rec.Code)
	}
}

func TestGatewayProtectedFeedWithValidToken(t *testing.T) {
	secret := "gateway-test-secret-32-chars-long!"
	cfg := router.GatewayConfig{
		IdentityServiceURL: "http://localhost:8081",
		JWTSecret:          secret,
	}
	h := router.NewGatewayRouter(cfg)

	jwtMgr := auth.NewJWTManager(secret, time.Hour, 24*time.Hour)
	testUser := &domain.User{
		ID:    "test-user-uuid-1234",
		Email: "test@example.com",
		Role:  domain.RoleUserStandard,
	}
	tokens, err := jwtMgr.GenerateTokens(testUser)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK for valid token, got: %d (body: %s)", rec.Code, rec.Body.String())
	}
}
