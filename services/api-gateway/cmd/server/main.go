package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dnjstjqppk-dev/facebook/services/api-gateway/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	identityURL := os.Getenv("IDENTITY_SERVICE_URL")
	if identityURL == "" {
		identityURL = "http://localhost:8081"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nova-enterprise-secure-jwt-secret-key-32chars"
	}

	cfg := router.GatewayConfig{
		IdentityServiceURL: identityURL,
		JWTSecret:          jwtSecret,
	}

	handler := router.NewGatewayRouter(cfg)

	log.Printf("[API Gateway] Listening on port :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Gateway server error: %v", err)
	}
}
