package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Anand078/url-shortener/internal/shortener"
	"github.com/Anand078/url-shortener/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// setupTestEnvironment creates test dependencies
func setupTestEnvironment() (*Handler, *mux.Router) {
	store := storage.NewInMemoryStorage()
	shortenerService := shortener.NewShortener(store)
	handler := NewHandler(shortenerService, store)
	router := SetupRoutes(handler)

	return handler, router
}

// TestShortenURLHandler tests the URL shortening endpoint
func TestShortenURLHandler(t *testing.T) {
	_, router := setupTestEnvironment()

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		validateBody   bool
	}{
		{
			name: "Valid URL",
			requestBody: map[string]string{
				"url": "https://example.com",
			},
			expectedStatus: http.StatusOK,
			validateBody:   true,
		},
		{
			name: "Invalid URL",
			requestBody: map[string]string{
				"url": "ftp://example.com",
			},
			expectedStatus: http.StatusBadRequest,
			validateBody:   false,
		},
		{
			name:           "Empty request",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			validateBody:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/shorten", bytes.NewReader(bodyBytes))

			// Execute request
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// Validate response body for successful requests
			if tt.validateBody {
				var response ShortenURLResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.NoError(t, err)
				assert.NotEmpty(t, response.ShortURL)
				assert.NotEmpty(t, response.ShortCode)
				assert.Equal(t, tt.requestBody["url"], response.OriginalURL)
			}
		})
	}
}

// TestRedirectHandler tests the redirection endpoint
func TestRedirectHandler(t *testing.T) {
	handler, router := setupTestEnvironment()

	// Create a shortened URL first
	originalURL := "https://example.com"
	shortCode, _ := handler.shortener.Shorten(originalURL)

	// Test successful redirection
	req, _ := http.NewRequest("GET", "/r/"+shortCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, originalURL, recorder.Header().Get("Location"))

	// Test non-existent short code
	req, _ = http.NewRequest("GET", "/r/nonexistent", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

// TestTopDomainsHandler tests the metrics endpoint
func TestTopDomainsHandler(t *testing.T) {
	handler, router := setupTestEnvironment()

	// Add some domain metrics
	handler.store.IncrementDomainCount("example.com")
	handler.store.IncrementDomainCount("example.com")
	handler.store.IncrementDomainCount("test.com")

	// Test the endpoint
	req, _ := http.NewRequest("GET", "/metrics/top-domains", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check response
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []storage.DomainCount
	err := json.Unmarshal(recorder.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
	assert.Equal(t, "example.com", response[0].Domain)
	assert.Equal(t, 2, response[0].Count)
	assert.Equal(t, "test.com", response[1].Domain)
	assert.Equal(t, 1, response[1].Count)
}
