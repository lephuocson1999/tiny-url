# Tiny URL Shortener (MVP)

A minimal, idiomatic Go URL shortener service following Clean Architecture.

## Features
- Shorten long URLs to unique short codes
- Redirect from short code to original URL
- Get stats for a short code
- In-memory storage (MVP)

## Getting Started

### Prerequisites
- Go 1.21+

### Run the Server
```sh
go run ./cmd/server
```

Server listens on `:8080` by default.

### API Endpoints

#### Shorten URL
- **POST** `/api/v1/shorten`
- **Body:** `{ "url": "https://example.com" }`
- **Response:** `{ "short_url": "abc123", "expires_at": "..." }`

#### Redirect
- **GET** `/api/v1/{shortCode}`
- **Redirects** to original URL

#### Get Stats
- **GET** `/api/v1/stats/{shortCode}`
- **Response:**
  ```json
  {
    "original_url": "https://example.com",
    "short_url": "abc123",
    "access_count": 0,
    "created_at": "...",
    "last_accessed": null
  }
  ```

## Testing
```sh
go test ./...
```

---

This MVP is ready for extension with persistent storage, caching, and production features. 