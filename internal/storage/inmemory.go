package storage

import (
	"net/url"
	"sort"
	"sync"
)

// Storage interface defines methods for URL storage
type Storage interface {
	// Save stores a URL with its short code
	Save(shortCode, originalURL string)

	// Find retrieves the original URL for a given short code
	Find(shortCode string) (string, bool)

	// FindShortCode checks if URL already has a short code
	FindShortCode(originalURL string) (string, bool)

	// IncrementDomainCount increases the count for a domain
	IncrementDomainCount(domain string)

	// GetTopDomains returns the top N domains by count
	GetTopDomains(n int) []DomainCount
}

// DomainCount represents a domain and its count
type DomainCount struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

// InMemoryStorage implements Storage using maps
type InMemoryStorage struct {
	shortToLong  map[string]string
	longToShort  map[string]string
	domainCounts map[string]int
	mutex        sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		shortToLong:  make(map[string]string),
		longToShort:  make(map[string]string),
		domainCounts: make(map[string]int),
		mutex:        sync.RWMutex{},
	}
}

// Save stores a URL with its short code
func (s *InMemoryStorage) Save(shortCode, originalURL string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.shortToLong[shortCode] = originalURL
	s.longToShort[originalURL] = shortCode
}

// Find retrieves the original URL for a given short code
func (s *InMemoryStorage) Find(shortCode string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	originalURL, exists := s.shortToLong[shortCode]
	return originalURL, exists
}

// FindShortCode checks if URL already has a short code
func (s *InMemoryStorage) FindShortCode(originalURL string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	shortCode, exists := s.longToShort[originalURL]
	return shortCode, exists
}

// IncrementDomainCount increases the count for a domain
func (s *InMemoryStorage) IncrementDomainCount(domain string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.domainCounts[domain]++
}

// GetTopDomains returns the top N domains by count
func (s *InMemoryStorage) GetTopDomains(n int) []DomainCount {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convert map to slice for sorting
	domains := make([]DomainCount, 0, len(s.domainCounts))
	for domain, count := range s.domainCounts {
		domains = append(domains, DomainCount{
			Domain: domain,
			Count:  count,
		})
	}

	// Sort by count descending
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Count > domains[j].Count
	})

	// Return up to n domains
	if len(domains) > n {
		domains = domains[:n]
	}

	return domains
}

// ExtractDomain extracts domain from a URL string
func ExtractDomain(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}
