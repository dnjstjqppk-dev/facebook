package router

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/services/api-gateway/internal/middleware"
)

type GatewayConfig struct {
	IdentityServiceURL string
	JWTSecret          string
}

func NewGatewayRouter(cfg GatewayConfig) http.Handler {
	mux := http.NewServeMux()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, 15*time.Minute, 14*24*time.Hour)
	authMw := middleware.NewAuthMiddleware(jwtMgr)
	rateLimiter := middleware.NewRateLimiter(100, 100*time.Millisecond)

	// Health Check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "api-gateway",
			"version": "1.0.0",
		})
	})

	// Reverse proxy to Identity Service
	identityURL, err := url.Parse(cfg.IdentityServiceURL)
	if err == nil {
		proxy := httputil.NewSingleHostReverseProxy(identityURL)
		mux.Handle("/api/v1/auth/", proxy)
	}

	// Protected feed routes placeholder
	protectedFeed := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Welcome to protected feed route",
			"user_id": claims.UserID,
			"role":    claims.Role,
		})
	})
	mux.Handle("GET /api/v1/feed", authMw.RequireAuth(protectedFeed))

	// Chain middlewares: CORS -> Rate Limiter -> Mux
	var handler http.Handler = mux
	handler = rateLimiter.Middleware(handler)
	handler = middleware.CORSMiddleware(handler)

	return handler
}
