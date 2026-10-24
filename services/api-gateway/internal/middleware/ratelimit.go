package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type clientLimit struct {
	tokens     int
	lastRefill time.Time
}

type RateLimiter struct {
	mu         sync.Mutex
	clients    map[string]*clientLimit
	maxTokens  int
	refillRate time.Duration // time per token
}

func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:    make(map[string]*clientLimit),
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}

	// Periodic cleanup of stale clients
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			rl.mu.Lock()
			now := time.Now()
			for ip, cl := range rl.clients {
				if now.Sub(cl.lastRefill) > 10*time.Minute {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = strings.Split(xff, ",")[0]
		}

		rl.mu.Lock()
		cl, exists := rl.clients[ip]
		now := time.Now()

		if !exists {
			rl.clients[ip] = &clientLimit{
				tokens:     rl.maxTokens - 1,
				lastRefill: now,
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		elapsed := now.Sub(cl.lastRefill)
		tokensToAdd := int(elapsed / rl.refillRate)
		if tokensToAdd > 0 {
			cl.tokens += tokensToAdd
			if cl.tokens > rl.maxTokens {
				cl.tokens = rl.maxTokens
			}
			cl.lastRefill = now
		}

		if cl.tokens <= 0 {
			rl.mu.Unlock()
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		cl.tokens--
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
