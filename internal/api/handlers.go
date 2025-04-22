package api

import (
	"encoding/json"
	"net/http"

	"github.com/Anand078/url-shortener/internal/shortener"
	"github.com/Anand078/url-shortener/internal/storage"
	"github.com/gorilla/mux"
)

// Handler contains references to dependencies
type Handler struct {
	shortener *shortener.Shortener
	store     storage.Storage
}

// NewHandler creates a new API handler
func NewHandler(shortener *shortener.Shortener, store storage.Storage) *Handler {
	return &Handler{
		shortener: shortener,
		store:     store,
	}
}

// ShortenURLRequest represents the request body for URL shortening
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse represents the response for URL shortening
type ShortenURLResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	ShortCode   string `json:"short_code"`
}

// ShortenURLHandler handles URL shortening requests
func (h *Handler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenURLRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if URL is provided
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Generate short URL
	shortCode, err := h.shortener.Shorten(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Construct short URL
	shortURL := getBaseURL(r) + "/r/" + shortCode

	// Prepare response
	response := ShortenURLResponse{
		OriginalURL: req.URL,
		ShortURL:    shortURL,
		ShortCode:   shortCode,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RedirectHandler handles URL redirection
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get short code from URL path
	vars := mux.Vars(r)
	shortCode := vars["code"]

	// Look up original URL
	originalURL, exists := h.shortener.Expand(shortCode)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

// TopDomainsHandler returns top domains metrics
func (h *Handler) TopDomainsHandler(w http.ResponseWriter, r *http.Request) {
	// Get top 3 domains
	topDomains := h.store.GetTopDomains(3)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topDomains)
}

// getBaseURL constructs the base URL from the request
func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
