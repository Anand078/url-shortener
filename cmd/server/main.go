package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Anand078/url-shortener/internal/api"
	"github.com/Anand078/url-shortener/internal/shortener"
	"github.com/Anand078/url-shortener/internal/storage"
)

func main() {
	// Initialize storage
	store := storage.NewInMemoryStorage()

	// Initialize URL shortener service
	shortenerService := shortener.NewShortener(store)

	// Initialize API handlers
	handler := api.NewHandler(shortenerService, store)

	// Set up routes
	router := api.SetupRoutes(handler)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
