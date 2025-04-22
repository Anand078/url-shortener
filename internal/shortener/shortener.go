package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/Anand078/url-shortener/internal/storage"
)

// Shortener handles URL shortening operations
type Shortener struct {
	store storage.Storage
}

// NewShortener creates a new URL shortener
func NewShortener(store storage.Storage) *Shortener {
	return &Shortener{
		store: store,
	}
}

// Shorten creates a shortened version of the URL
func (s *Shortener) Shorten(originalURL string) (string, error) {
	// Validate URL
	if err := validateURL(originalURL); err != nil {
		return "", err
	}

	// Check if URL already has a short code
	if shortCode, exists := s.store.FindShortCode(originalURL); exists {
		return shortCode, nil
	}

	// Generate new short code
	shortCode := generateShortCode(originalURL)

	// Extract and store domain for metrics
	domain, err := storage.ExtractDomain(originalURL)
	if err == nil && domain != "" {
		s.store.IncrementDomainCount(domain)
	}

	// Save mapping
	s.store.Save(shortCode, originalURL)

	return shortCode, nil
}

// Expand returns the original URL for a given short code
func (s *Shortener) Expand(shortCode string) (string, bool) {
	return s.store.Find(shortCode)
}

// validateURL checks if the URL is valid
func validateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must have http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// generateShortCode creates a unique short code for a URL
func generateShortCode(url string) string {
	// Create SHA256 hash of the URL
	hash := sha256.Sum256([]byte(url))

	// Encode to base64 for shorter representation
	encoded := base64.URLEncoding.EncodeToString(hash[:])

	// Take first 8 characters for the short code
	// This provides enough uniqueness while keeping the URL short
	shortCode := encoded[:8]

	// Make it URL-safe by removing any special characters
	shortCode = strings.ReplaceAll(shortCode, "_", "a")
	shortCode = strings.ReplaceAll(shortCode, "-", "b")

	return shortCode
}
