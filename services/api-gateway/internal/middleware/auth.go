package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
)

type contextKey string

const (
	UserContextKey contextKey = "user_claims"
)

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"unauthorized: missing or invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := a.jwtManager.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, `{"error":"unauthorized: invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-User-Role", string(claims.Role))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
