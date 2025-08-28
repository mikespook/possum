# Possum - Go HTTP Middleware Toolkit

Possum is a lightweight Go package that provides a collection of HTTP middleware functions for building robust web applications. It includes authentication, CORS handling, logging, method filtering, WebSocket support, and standardized response formatting.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Components](#core-components)
  - [Middleware Chain](#middleware-chain)
  - [Authentication](#authentication)
  - [CORS Handling](#cors-handling)
  - [Logging](#logging)
  - [Method Filtering](#method-filtering)
  - [WebSocket Support](#websocket-support)
  - [Response Handling](#response-handling)
- [Configuration](#configuration)
- [Testing](#testing)
- [Dependencies](#dependencies)

## Installation

```bash
go get github.com/mikespook/possum
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/mikespook/possum"
)

func main() {
    // Create a simple handler
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    }

    // Apply middleware using Chain
    http.HandleFunc("/", possum.Chain(
        handler,
        possum.Log,  // Add logging
        possum.Cors(nil), // Add CORS with default config
    ))

    http.ListenAndServe(":8080", nil)
}
```

## Core Components

### Middleware Chain

The `Chain` function allows composing multiple middleware handlers into a single handler, applying them in reverse order.

```go
func Chain(handler http.HandlerFunc, middlewares ...HandlerFunc) http.HandlerFunc
```

**Usage:**
```go
chainedHandler := possum.Chain(
    myHandler,
    possum.Log,
    possum.Cors(corsConfig),
    possum.HTTPAuth(secret, nil),
)
```

### Authentication

Possum provides JWT-based authentication for both HTTP and WebSocket connections.

#### HTTP Authentication

```go
func HTTPAuth(secret []byte, next http.HandlerFunc) http.HandlerFunc
```

**Usage:**
```go
// Protect a route with JWT authentication
protectedHandler := possum.HTTPAuth(secretKey, myHandler)
```

#### WebSocket Authentication

```go
func WebSocketAuth(secret []byte, next WebsocketHandlerFunc) WebsocketHandlerFunc
```

**Usage:**
```go
// Protect a WebSocket endpoint
wsHandler := possum.WebSocketAuth(secretKey, myWebSocketHandler)
```

### CORS Handling

Cross-Origin Resource Sharing (CORS) middleware for handling browser security restrictions.

```go
func Cors(config *CORSConfig, next http.HandlerFunc) http.HandlerFunc
```

**Configuration:**
```go
type CORSConfig struct {
    AllowOrigin      string   // Default: "*"
    AllowedOrigins   []string // Specific origins to allow
    AllowMethods     string   // Default: "*"
    AllowHeaders     string   // Default: "*"
    AllowCredentials bool     // Default: true
    ExposeHeaders    string   // Default: "*"
    MaxAge           int      // Default: 0
    ExemptMethods    []string // Default: ["OPTIONS"]
}
```

**Usage:**
```go
corsConfig := &possum.CORSConfig{
    AllowedOrigins: []string{"https://example.com"},
    AllowMethods:   "GET,POST,PUT,DELETE",
    AllowHeaders:   "Content-Type,Authorization",
}

handler := possum.Cors(corsConfig, myHandler)
```

### Logging

Request/response logging middleware using zerolog for structured logging.

```go
func Log(next http.HandlerFunc) http.HandlerFunc
```

**Usage:**
```go
loggedHandler := possum.Log(myHandler)
```

The logger captures:
- Request details (method, URL, headers, etc.)
- Response status codes
- Processing errors

### Method Filtering

Filter HTTP methods with allow and deny lists.

#### Allow Methods

```go
func AllowMethods(methods ...string) HandlerFunc
```

**Usage:**
```go
// Only allow GET and POST methods
handler := possum.Chain(
    myHandler,
    possum.AllowMethods("GET", "POST"),
)
```

#### Deny Methods

```go
func DenyMethods(methods ...string) HandlerFunc
```

**Usage:**
```go
// Block DELETE and PUT methods
handler := possum.Chain(
    myHandler,
    possum.DenyMethods("DELETE", "PUT"),
)
```

### WebSocket Support

WebSocket upgrade handler with built-in CORS support and connection management.

```go
func WebSocketUpgrade(corsConfig *CORSConfig, next WebsocketHandlerFunc) http.HandlerFunc
```

**Usage:**
```go
wsHandler := func(conn *websocket.Conn, r *http.Request) {
    // Handle WebSocket communication
    // Ping/pong and connection cleanup are handled automatically
}

http.HandleFunc("/ws", possum.WebSocketUpgrade(nil, wsHandler))
```

### Response Handling

Standardized JSON response structure with UUID tracking and error handling.

#### Response Structure

```go
type Response struct {
    UUID  uuid.UUID `json:"uuid"`
    Data  any       `json:"data,omitempty"`
    Error *Error    `json:"error,omitempty"`
}
```

#### Predefined Responses

Common HTTP error responses:
- `InternalServerErrorResponse`
- `UnauthorizedResponse`
- `MethodNotAllowedResponse`
- `BadRequestResponse`
- `NotFoundResponse`
- `NotImplementedResponse`
- `ConflictResponse`
- `ForbiddenResponse`

#### Creating Responses

```go
// Create a new response
resp := possum.NewResponse(r)

// Set data
resp.SetData(map[string]string{"message": "Hello"})

// Set error
resp.SetError(http.StatusBadRequest, "Invalid request")

// Write response
resp.Write(w)
```

## Configuration

### Environment-Based Configuration

The package uses environment variables for configuration:

```bash
# Set environment (development, production, test)
export POSSUM_ENV=production
```

### Logger Configuration

The logging system can be configured through the log package:

```go
logConfig := &log.Config{
    Level:    "info",      // trace, debug, info, warn, error, fatal, panic
    Filename: "app.log",   // File to write logs to, empty for stderr
}
log.Init(logConfig)
```

## Testing

Run tests with the provided script:

```bash
./run_tests.sh
```

Or manually:

```bash
go test -v ./...
```

## Dependencies

- [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) - JWT implementation
- [github.com/google/uuid](https://github.com/google/uuid) - UUID generation
- [github.com/gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [github.com/rs/zerolog](https://github.com/rs/zerolog) - Structured logging

## Reference Guide

### Main Package Functions

| Function | Description |
|---------|-------------|
| `Chain` | Composes multiple middleware into a single handler |
| `HTTPAuth` | JWT authentication for HTTP requests |
| `WebSocketAuth` | JWT authentication for WebSocket connections |
| `Cors` | CORS header management |
| `Log` | Request/response logging |
| `AllowMethods` | Restrict allowed HTTP methods |
| `DenyMethods` | Block specific HTTP methods |
| `WebSocketUpgrade` | WebSocket connection handler |

### Response Functions

| Function | Description |
|---------|-------------|
| `NewResponse` | Creates a new Response with UUID |
| `CloneResponse` | Deep copies a Response object |
| `WriteResponse` | Writes a Response with error handling |
| `SetError` | Sets error details in a Response |
| `SetData` | Sets data in a Response |
| `WriteHeader` | Sets HTTP status code |
| `Write` | Serializes Response to HTTP writer |

### CORS Configuration

| Field | Type | Description |
|------|------|-------------|
| `AllowOrigin` | string | Allowed origin pattern |
| `AllowedOrigins` | []string | List of specific allowed origins |
| `AllowMethods` | string | Allowed HTTP methods |
| `AllowHeaders` | string | Allowed headers |
| `AllowCredentials` | bool | Whether to allow credentials |
| `ExposeHeaders` | string | Headers to expose to client |
| `MaxAge` | int | How long to cache preflight response |
| `ExemptMethods` | []string | Methods that trigger preflight |

### Authentication Claims

JWT claims structure:

```go
type JWTClaims struct {
    UserID    uuid.UUID
    IssuedAt  time.Time
    ExpiresAt time.Time
    jwt.RegisteredClaims
}
```

### Constants

| Constant | Value | Description |
|---------|-------|-------------|
| `ClaimsKey` | `"claims"` | Context key for JWT claims |
| `UUIDKey` | `"uuid"` | Context key for request UUID |

## Examples

### Complete Example with Authentication

```go
package main

import (
    "net/http"
    "github.com/mikespook/possum"
    "github.com/google/uuid"
    "github.com/mikespook/possum/auth"
)

func main() {
    // Secret key for JWT
    secret := []byte("my-secret-key")
    
    // Protected endpoint
    protectedHandler := func(w http.ResponseWriter, r *http.Request) {
        // Get user claims from context
        claims, ok := r.Context().Value(possum.ClaimsKey).(*auth.JWTClaims)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Use claims
        resp := possum.NewResponse(r)
        resp.SetData(map[string]interface{}{
            "user_id": claims.UserID,
            "message": "Hello, authenticated user!",
        })
        resp.Write(w)
    }
    
    // Public endpoint to generate tokens
    loginHandler := func(w http.ResponseWriter, r *http.Request) {
        // In a real app, verify credentials here
        userID := uuid.New()
        _, token, err := auth.GenerateJWT(secret, userID, nil)
        if err != nil {
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }
        
        resp := possum.NewResponse(r)
        resp.SetData(map[string]string{
            "token": token,
        })
        resp.Write(w)
    }
    
    // Apply middleware
    http.HandleFunc("/login", possum.Chain(
        loginHandler,
        possum.Log,
        possum.Cors(nil),
    ))
    
    http.HandleFunc("/protected", possum.Chain(
        protectedHandler,
        possum.Log,
        possum.Cors(nil),
        possum.HTTPAuth(secret, nil),
    ))
    
    http.ListenAndServe(":8080", nil)
}
```