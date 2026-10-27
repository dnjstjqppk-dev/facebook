package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/handler"
	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	repo := repository.NewInMemoryFeedRepository()
	feedSvc := service.NewFeedService(repo, nil)
	feedHandler := handler.NewFeedHandler(feedSvc)

	mux := http.NewServeMux()
	feedHandler.RegisterRoutes(mux)

	log.Printf("Feed Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

