# URL Shortener Service

A simple URL shortener service written in Go that provides an API to shorten URLs, redirect to the original URLs, and track domain metrics.

## Features

- URL shortening with custom short code generation
- Consistent short URLs (same input URL always produces the same short URL)
- Redirection from short URL to original URL
- In-memory storage for URL mappings
- Metrics tracking for top domains
- Comprehensive unit tests
- Docker support