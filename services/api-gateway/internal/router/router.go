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
	FeedServiceURL        string
	MediaServiceURL       string
	MessagingServiceURL   string
	GroupsServiceURL      string
	PagesServiceURL       string
	EventsServiceURL      string
	MarketplaceServiceURL string
	PaymentsServiceURL    string
	AdsServiceURL         string
	ModerationServiceURL  string
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

	// Reverse proxy to Feed Service
	if cfg.FeedServiceURL != "" {
		if u, err := url.Parse(cfg.FeedServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/posts/", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/feed", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/memories", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/saved", authMw.RequireAuth(proxy))
		}
	} else {
		// Protected feed routes fallback
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
	}

	// Reverse proxy to Media Service
	if cfg.MediaServiceURL != "" {
		if u, err := url.Parse(cfg.MediaServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/media/", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/stories", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/stories/", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/reels", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/live/", authMw.RequireAuth(proxy))
		}
	}

	// Reverse proxy to Messaging Service
	if cfg.MessagingServiceURL != "" {
		if u, err := url.Parse(cfg.MessagingServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/conversations/", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/messages", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/messages/", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/presence", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/presence/", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/calls", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/calls/", authMw.RequireAuth(proxy))
		}
	}

	// Reverse proxy to Community Services (Groups, Pages, Events)
	if cfg.GroupsServiceURL != "" {
		if u, err := url.Parse(cfg.GroupsServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/groups", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/groups/", authMw.RequireAuth(proxy))
		}
	}
	if cfg.PagesServiceURL != "" {
		if u, err := url.Parse(cfg.PagesServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/pages", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/pages/", authMw.RequireAuth(proxy))
		}
	}
	if cfg.EventsServiceURL != "" {
		if u, err := url.Parse(cfg.EventsServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/events", authMw.RequireAuth(proxy))
			mux.Handle("/api/v1/events/", authMw.RequireAuth(proxy))
		}
	}

	// Reverse proxy to Commerce Services (Marketplace, Payments)
	if cfg.MarketplaceServiceURL != "" {
		if u, err := url.Parse(cfg.MarketplaceServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/marketplace/", authMw.RequireAuth(proxy))
		}
	}
	if cfg.PaymentsServiceURL != "" {
		if u, err := url.Parse(cfg.PaymentsServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/payments/", authMw.RequireAuth(proxy))
		}
	}

	// Reverse proxy to Ads and Moderation
	if cfg.AdsServiceURL != "" {
		if u, err := url.Parse(cfg.AdsServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/ads/", authMw.RequireAuth(proxy))
		}
	}
	if cfg.ModerationServiceURL != "" {
		if u, err := url.Parse(cfg.ModerationServiceURL); err == nil {
			proxy := httputil.NewSingleHostReverseProxy(u)
			mux.Handle("/api/v1/moderation/", authMw.RequireAuth(proxy))
		}
	}

	// Chain middlewares: CORS -> Rate Limiter -> Mux
	var handler http.Handler = mux
	handler = rateLimiter.Middleware(handler)
	handler = middleware.CORSMiddleware(handler)

	return handler
}
