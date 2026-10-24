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
	IdentityServiceURL    string
	SocialGraphServiceURL string
	ProfilesServiceURL    string
	JWTSecret             string
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
	if cfg.IdentityServiceURL != "" {
		if u, err := url.Parse(cfg.IdentityServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/auth/", proxy)
		}
	}

	// Reverse proxy to Social Graph Service
	if cfg.SocialGraphServiceURL != "" {
		if u, err := url.Parse(cfg.SocialGraphServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/friends/", proxy)
			mux.Handle("/api/v1/users/block", proxy)
		}
	}

	// Reverse proxy to Profiles Service
	if cfg.ProfilesServiceURL != "" {
		if u, err := url.Parse(cfg.ProfilesServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/profiles/", proxy)
		}
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
