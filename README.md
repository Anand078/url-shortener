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

## Project Structure

url-shortener/
├── cmd/
│   └── server/
│       └── main.go         # Entry point
├── internal/
│   ├── api/                # API handlers
│   │   ├── handlers.go
│   │   ├── handlers_test.go
│   │   └── routes.go
│   ├── storage/            # In-memory storage
│   │   ├── inmemory.go
│   │   └── inmemory_test.go
│   └── shortener/          # Core shortening logic
│       ├── shortener.go
│       └── shortener_test.go
├── Dockerfile              # For Docker containerization
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies lock file
└── README.md               # This file

## Project Structure

url-shortener/
├── cmd/
│   └── server/
│       └── main.go         # Entry point
├── internal/
│   ├── api/                # API handlers
│   │   ├── handlers.go
│   │   ├── handlers_test.go
│   │   └── routes.go
│   ├── storage/            # In-memory storage
│   │   ├── inmemory.go
│   │   └── inmemory_test.go
│   └── shortener/          # Core shortening logic
│       ├── shortener.go
│       └── shortener_test.go
├── Dockerfile              # For Docker containerization
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies lock file
└── README.md               # This file

## Requirements

- Go 1.20 or later
- Docker (optional, for containerization)

## Getting Started

### Running Locally

#### 1. Clone the repository:

git clone https://github.com/yourusername/url-shortener.git
cd url-shortener

#### 2. Install dependencies:
go mod download

#### 3. Run the application:
go run cmd/server/main.go

#### 4. The server will start on `http://localhost:8080`

### Using Docker

#### 1. Build the Docker image:
docker build -t url-shortener .

#### 2. Run the container:
docker run -p 8080:8080 url-shortener

#### 3. The service will be available at `http://localhost:8080`

## API Endpoints

### 1. Shorten URL

**Endpoint:** `POST /shorten`

**Request Body:**
```json
{
"url": "https://example.com/some/long/path"
}
Response:
json{
  "original_url": "https://example.com/some/long/path",
  "short_url": "http://localhost:8080/r/abc123de",
  "short_code": "abc123de"
}
```

### 2. Redirect to Original URL
Endpoint: GET /r/{code}
This endpoint redirects to the original URL associated with the provided short code.

 ### 3. Get Top Domains

**Endpoint:** `GET /metrics/top-domains`

**Response:**
```json
[
  {
    "domain": "example.com",
    "count": 5
  },
  {
    "domain": "google.com",
    "count": 3
  },
  {
    "domain": "github.com",
    "count": 2
  }
]
```

## Testing
Run the full test suite:

go test ./...


## Design Decisions

### URL Shortening Algorithm
The service uses a SHA-256 hash of the input URL and then takes the first 8 characters of the base64-encoded hash for the short code. This provides:

- Deterministic behavior (same input always produces same output)
- Fixed-length short codes

### Storage
The service uses in-memory storage with mutex locks for thread safety. The storage layer is designed with interfaces to allow for future extension to persistent storage options.

### API Design
The API is designed RESTfully with clear separation of concerns:

- shortener package handles the core URL shortening logic
- storage package manages data persistence
- api package handles HTTP requests and responses

## Future Improvements

- Persistent storage (database) to maintain URLs across restarts
- Rate limiting to prevent abuse
- Authentication and user management
- URL expiration functionality

