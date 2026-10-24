package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity"
	"github.com/google/uuid"
)

type Config struct {
	Port            int
	JWTSecret       string
	AllowedOrigins  []string
	RateLimitPerMin int
}

type Gateway struct {
	config      Config
	authService *identity.AuthService
	rateLimiter *InMemoryRateLimiter
}

type InMemoryRateLimiter struct {
	mu      sync.Mutex
	limits  map[string]*ClientLimiter
	maxReqs int
}

type ClientLimiter struct {
	lastReset time.Time
	count     int
}

func NewGateway(cfg Config) *Gateway {
	return &Gateway{
		config:      cfg,
		authService: identity.NewAuthService(cfg.JWTSecret),
		rateLimiter: &InMemoryRateLimiter{
			limits:  make(map[string]*ClientLimiter),
			maxReqs: cfg.RateLimitPerMin,
		},
	}
}

// Handler returns the root HTTP handler with middleware chain
func (g *Gateway) Handler() http.Handler {
	mux := http.NewServeMux()

	// Public Health Check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "UP",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	// Public Auth proxy routes
	mux.HandleFunc("/api/v1/auth/login", g.handleLogin)
	mux.HandleFunc("/api/v1/auth/register", g.handleRegister)

	// Protected Feed & Profile routes
	mux.Handle("/api/v1/feed", g.authMiddleware(http.HandlerFunc(g.handleFeed)))
	mux.Handle("/api/v1/users/me", g.authMiddleware(http.HandlerFunc(g.handleCurrentUser)))

	return g.loggingMiddleware(g.corsMiddleware(g.rateLimitMiddleware(mux)))
}

// Middleware: CORS
func (g *Gateway) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := false
		for _, o := range g.config.AllowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware: Request Logging & Request-ID
func (g *Gateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", reqID)

		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s %s took %v", reqID, r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// Middleware: Rate Limiting
func (g *Gateway) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]

		g.rateLimiter.mu.Lock()
		limiter, exists := g.rateLimiter.limits[ip]
		now := time.Now()
		if !exists || now.Sub(limiter.lastReset) > time.Minute {
			g.rateLimiter.limits[ip] = &ClientLimiter{lastReset: now, count: 1}
			g.rateLimiter.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if limiter.count >= g.rateLimiter.maxReqs {
			g.rateLimiter.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded. Try again in a minute."})
			return
		}

		limiter.count++
		g.rateLimiter.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// Middleware: JWT Authentication
func (g *Gateway) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing Authorization header"})
			return
		}

		claims, err := g.authService.ValidateAccessToken(authHeader)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Inject user context headers for upstream handlers/microservices
		r.Header.Set("X-User-ID", claims.UserID.String())
		r.Header.Set("X-User-Email", claims.Email)
		r.Header.Set("X-User-Role", claims.Role)

		next.ServeHTTP(w, r)
	})
}

// Stubs for API endpoints
func (g *Gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Login endpoint ready"})
}

func (g *Gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Registration endpoint ready"})
}

func (g *Gateway) handleFeed(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"feed":    []interface{}{},
		"has_more": false,
	})
}

func (g *Gateway) handleCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	email := r.Header.Get("X-User-Email")
	role := r.Header.Get("X-User-Role")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}

func (g *Gateway) Start() error {
	addr := fmt.Sprintf(":%d", g.config.Port)
	log.Printf("Starting API Gateway on %s", addr)
	return http.ListenAndServe(addr, g.Handler())
}
