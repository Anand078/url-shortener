package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInMemoryStorage tests the in-memory storage functionality
func TestInMemoryStorage(t *testing.T) {
	storage := NewInMemoryStorage()

	// Test Save and Find
	shortCode := "abc123"
	originalURL := "https://example.com"

	storage.Save(shortCode, originalURL)

	foundURL, exists := storage.Find(shortCode)
	assert.True(t, exists)
	assert.Equal(t, originalURL, foundURL)

	// Test FindShortCode
	foundCode, exists := storage.FindShortCode(originalURL)
	assert.True(t, exists)
	assert.Equal(t, shortCode, foundCode)

	// Test non-existent short code
	_, exists = storage.Find("nonexistent")
	assert.False(t, exists)

	// Test non-existent URL
	_, exists = storage.FindShortCode("https://nonexistent.com")
	assert.False(t, exists)
}

// TestDomainMetrics tests the domain metrics functionality
func TestDomainMetrics(t *testing.T) {
	storage := NewInMemoryStorage()

	// Add domain counts
	storage.IncrementDomainCount("example.com")
	storage.IncrementDomainCount("example.com")
	storage.IncrementDomainCount("example.com")
	storage.IncrementDomainCount("test.com")
	storage.IncrementDomainCount("test.com")
	storage.IncrementDomainCount("another.com")

	// Get top domains
	topDomains := storage.GetTopDomains(2)

	// Check results
	assert.Equal(t, 2, len(topDomains))
	assert.Equal(t, "example.com", topDomains[0].Domain)
	assert.Equal(t, 3, topDomains[0].Count)
	assert.Equal(t, "test.com", topDomains[1].Domain)
	assert.Equal(t, 2, topDomains[1].Count)

	// Test with limit higher than available domains
	allDomains := storage.GetTopDomains(10)
	assert.Equal(t, 3, len(allDomains))
}

// TestExtractDomain tests the domain extraction function
func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedDomain string
		expectError    bool
	}{
		{
			name:           "Simple URL",
			url:            "https://example.com",
			expectedDomain: "example.com",
			expectError:    false,
		},
		{
			name:           "URL with path",
			url:            "https://example.com/test/path",
			expectedDomain: "example.com",
			expectError:    false,
		},
		{
			name:           "URL with subdomain",
			url:            "https://sub.example.com",
			expectedDomain: "sub.example.com",
			expectError:    false,
		},
		{
			name:           "URL with port",
			url:            "https://example.com:8080",
			expectedDomain: "example.com",
			expectError:    false,
		},
		{
			name:           "Invalid URL",
			url:            "not-a-url",
			expectedDomain: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := ExtractDomain(tt.url)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDomain, domain)
			}
		})
	}
}
