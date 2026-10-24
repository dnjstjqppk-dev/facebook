package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestGateway_HealthCheck(t *testing.T) {
	gw := NewGateway(Config{
		Port:            8080,
		JWTSecret:       "test-secret-at-least-32-bytes-long!",
		AllowedOrigins:  []string{"*"},
		RateLimitPerMin: 100,
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	gw.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestGateway_UnauthorizedAccess(t *testing.T) {
	gw := NewGateway(Config{
		Port:            8080,
		JWTSecret:       "test-secret-at-least-32-bytes-long!",
		AllowedOrigins:  []string{"*"},
		RateLimitPerMin: 100,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rr := httptest.NewRecorder()

	gw.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing token, got %d", rr.Code)
	}
}

func TestGateway_AuthorizedAccess(t *testing.T) {
	secret := "test-secret-at-least-32-bytes-long!"
	gw := NewGateway(Config{
		Port:            8080,
		JWTSecret:       secret,
		AllowedOrigins:  []string{"*"},
		RateLimitPerMin: 100,
	})

	tokenPair, err := gw.authService.GenerateTokenPair(uuid.New(), "test@user.com", "ROLE_USER", true)
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	rr := httptest.NewRecorder()

	gw.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid token, got %d", rr.Code)
	}
}
