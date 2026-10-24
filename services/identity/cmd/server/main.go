package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/handler"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nova-enterprise-secure-jwt-secret-key-32chars"
	}

	repo := repository.NewInMemoryUserRepository()
	jwtManager := service.NewJWTManager(jwtSecret, 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(repo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()
	authHandler.RegisterRoutes(mux)

	log.Printf("[Identity Service] Listening on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
