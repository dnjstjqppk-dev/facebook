package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dnjstjqppk-dev/facebook/services/media/internal/handler"
	"github.com/dnjstjqppk-dev/facebook/services/media/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/media/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	repo := repository.NewInMemoryMediaRepository()
	mediaSvc := service.NewMediaService(repo)
	mediaHandler := handler.NewMediaHandler(mediaSvc)

	mux := http.NewServeMux()
	mediaHandler.RegisterRoutes(mux)

	log.Printf("Media Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

