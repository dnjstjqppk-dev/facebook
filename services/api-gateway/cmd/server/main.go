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

	socialGraphURL := os.Getenv("SOCIAL_GRAPH_SERVICE_URL")
	if socialGraphURL == "" {
		socialGraphURL = "http://localhost:8082"
	}

	profilesURL := os.Getenv("PROFILES_SERVICE_URL")
	if profilesURL == "" {
		profilesURL = "http://localhost:8083"
	}

	feedURL := os.Getenv("FEED_SERVICE_URL")
	if feedURL == "" {
		feedURL = "http://localhost:8084"
	}

	mediaURL := os.Getenv("MEDIA_SERVICE_URL")
	if mediaURL == "" {
		mediaURL = "http://localhost:8085"
	}

	messagingURL := os.Getenv("MESSAGING_SERVICE_URL")
	if messagingURL == "" {
		messagingURL = "http://localhost:8086"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nova-enterprise-secure-jwt-secret-key-32chars"
	}

	cfg := router.GatewayConfig{
		IdentityServiceURL:    identityURL,
		SocialGraphServiceURL: socialGraphURL,
		ProfilesServiceURL:    profilesURL,
		FeedServiceURL:        feedURL,
		MediaServiceURL:       mediaURL,
		MessagingServiceURL:   messagingURL,
		JWTSecret:             jwtSecret,
	}

	handler := router.NewGatewayRouter(cfg)

	log.Printf("[API Gateway] Listening on port :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Gateway server error: %v", err)
	}
}
