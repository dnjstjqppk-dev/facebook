package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/handler"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nova-enterprise-secure-jwt-secret-key-32chars"
	}

	repo := repository.NewInMemoryProfileRepository()
	jwtMgr := auth.NewJWTManager(jwtSecret, 15*time.Minute, 14*24*time.Hour)
	profileSvc := service.NewProfileService(repo)
	profileHandler := handler.NewProfileHandler(profileSvc, jwtMgr)

	mux := http.NewServeMux()
	profileHandler.RegisterRoutes(mux)

	log.Printf("[Profiles Service] Listening on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
