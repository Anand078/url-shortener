package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/shorten", handler.ShortenURLHandler).Methods("POST")
	router.HandleFunc("/r/{code}", handler.RedirectHandler).Methods("GET")
	router.HandleFunc("/metrics/top-domains", handler.TopDomainsHandler).Methods("GET")

	// Add middleware for logging
	router.Use(loggingMiddleware)

	return router
}

// loggingMiddleware logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		// fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
