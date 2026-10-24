package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/handler"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nova-enterprise-secure-jwt-secret-key-32chars"
	}

	repo := repository.NewInMemorySocialGraphRepository()
	jwtMgr := auth.NewJWTManager(jwtSecret, 15*time.Minute, 14*24*time.Hour)
	socialSvc := service.NewSocialGraphService(repo)
	socialHandler := handler.NewSocialHandler(socialSvc, jwtMgr)

	mux := http.NewServeMux()
	socialHandler.RegisterRoutes(mux)

	log.Printf("[Social Graph Service] Listening on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
