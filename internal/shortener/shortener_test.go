package shortener

import (
	"testing"

	"github.com/Anand078/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
)

// TestShorten tests the URL shortening functionality
func TestShorten(t *testing.T) {
	store := storage.NewInMemoryStorage()
	shortener := NewShortener(store)

	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "Valid URL",
			url:         "https://example.com/test",
			expectError: false,
		},
		{
			name:        "Invalid URL scheme",
			url:         "ftp://example.com",
			expectError: true,
		},
		{
			name:        "Missing scheme",
			url:         "example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortCode, err := shortener.Shorten(tt.url)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, shortCode)

				// Test consistency - same URL should get same code
				shortCode2, err := shortener.Shorten(tt.url)
				assert.NoError(t, err)
				assert.Equal(t, shortCode, shortCode2)
			}
		})
	}
}

// TestExpand tests the URL expansion functionality
func TestExpand(t *testing.T) {
	store := storage.NewInMemoryStorage()
	shortener := NewShortener(store)

	// Create a short URL
	originalURL := "https://example.com/test"
	shortCode, err := shortener.Shorten(originalURL)
	assert.NoError(t, err)

	// Test expansion
	expandedURL, exists := shortener.Expand(shortCode)
	assert.True(t, exists)
	assert.Equal(t, originalURL, expandedURL)

	// Test non-existent code
	_, exists = shortener.Expand("nonexistent")
	assert.False(t, exists)
}

// TestValidateURL tests URL validation
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "Valid HTTP URL",
			url:         "http://example.com",
			expectError: false,
		},
		{
			name:        "Valid HTTPS URL",
			url:         "https://example.com",
			expectError: false,
		},
		{
			name:        "Invalid scheme",
			url:         "ftp://example.com",
			expectError: true,
		},
		{
			name:        "Missing scheme",
			url:         "example.com",
			expectError: true,
		},
		{
			name:        "Empty URL",
			url:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
