package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/handler"
	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	repo := repository.NewInMemoryMessagingRepository()
	messagingSvc := service.NewMessagingService(repo)
	messagingHandler := handler.NewMessagingHandler(messagingSvc)

	mux := http.NewServeMux()
	messagingHandler.RegisterRoutes(mux)

	log.Printf("Messaging Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

