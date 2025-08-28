# Possum Package - LLM Coder Guide

Possum is a lightweight, modular Go web framework that provides essential middleware functionality for building HTTP APIs and web applications. It offers a collection of independent packages that can be used together or separately.

## Table of Contents
- [Overview](#overview)
- [Core Concepts](#core-concepts)
- [Package Structure](#package-structure)
- [Key Components](#key-components)
  - [Auth](#auth)
  - [Chain](#chain)
  - [CORS](#cors)
  - [Logger](#logger)
  - [Method](#method)
  - [Response](#response)
  - [WebSocket](#websocket)
- [Usage Examples](#usage-examples)
- [Testing](#testing)
- [Dependencies](#dependencies)
- [Development Guidelines](#development-guidelines)

## Overview

Possum is designed as a collection of modular middleware components for Go web applications. Each component is self-contained and can be used independently or in combination with others. The framework emphasizes simplicity, flexibility, and ease of use.

Possum provides a toolkit of HTTP middleware functions for building robust web applications and APIs, including authentication, CORS handling, logging, method filtering, WebSocket support, and standardized response formatting.

## Core Concepts

1. **Middleware Pattern**: Each component follows the standard Go HTTP middleware pattern, wrapping handlers to add functionality.
2. **Modularity**: Components are independent and can be mixed and matched based on project needs.
3. **Chainable**: The `chain` package allows combining multiple middleware into a single handler.
4. **Standards Compliant**: All components follow HTTP standards and best practices.
5. **Context Integration**: Uses request context to pass data between middleware and handlers.
6. **Environment Awareness**: Behavior can be controlled through environment variables like `POSSUM_ENV`.

## Package Structure

The possum package contains several independent modules:

```
possum/
├── auth/                 # Authentication utilities
│   ├── jwt.go           # JWT token generation and parsing
│   └── jwt_test.go      # Tests for JWT functionality
├── config/              # Configuration utilities
│   ├── config.go        # Environment-based configuration
│   └── config_test.go   # Tests for configuration
├── log/                 # Logging utilities
│   ├── config.go        # Logger configuration
│   ├── logger.go        # Logger implementation
│   └── logger_test.go   # Tests for logger
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── *.go                 # Main package files
└── *_test.go            # Tests for main package files
```

### Main Package Files

1. `auth.go` - JWT-based authentication middleware
2. `chain.go` - Utility for chaining multiple middleware handlers
3. `cors.go` - Cross-Origin Resource Sharing (CORS) middleware
4. `logger.go` - Request/response logging middleware
5. `method.go` - HTTP method override functionality
6. `response.go` - Standardized JSON response formatting
7. `websocket.go` - WebSocket connection handling utilities

Each module has corresponding test files (e.g., `auth_test.go`).

## Key Components

### Auth

The `auth` package provides JWT-based authentication middleware for both HTTP and WebSocket connections.

**Key Features:**
- JWT token validation for HTTP requests
- JWT token validation for WebSocket connections
- Token generation with custom claims
- Configurable signing methods (HMAC, RSA, ECDSA)
- Role-based access control support
- Context integration for claims passing

**Main Functions:**
- `HTTPAuth(secret []byte, next http.HandlerFunc) http.HandlerFunc`: Middleware that validates JWT tokens for HTTP requests
- `WebSocketAuth(secret []byte, next WebsocketHandlerFunc) WebsocketHandlerFunc`: Middleware that validates JWT tokens for WebSocket connections
- `GenerateJWT(secret []byte, userID uuid.UUID, customClaims jwt.Claims) (*jwt.Token, string, error)`: Generates a new JWT token
- Context integration using `ClaimsKey` to store and retrieve claims

**Context Integration:**
- Stores claims in request context for access in downstream handlers
- Uses `context.WithValue` to pass claims through the request lifecycle
- Access claims using `r.Context().Value(possum.ClaimsKey)`

**JWT Claims Structure:**
```go
type JWTClaims struct {
    UserID    uuid.UUID
    IssuedAt  time.Time
    ExpiresAt time.Time
    jwt.RegisteredClaims
}
```

### Chain

The `chain` package provides utilities for combining multiple middleware into a single handler. In the main package, the `Chain` function is used to compose middleware.

**Key Features:**
- Simple middleware chaining
- Preserves order of execution (applies middleware in reverse order)
- Compatible with standard `http.HandlerFunc`

**Main Functions:**
- `Chain(handler http.HandlerFunc, middlewares ...HandlerFunc) http.HandlerFunc`: Composes multiple middleware handlers into a single handler, applying them in reverse order

**Usage Pattern:**
```go
handler := possum.Chain(
    myHandler,
    possum.Log,           // Add request logging
    possum.Cors(nil),     // Add CORS with default config
    possum.HTTPAuth(secretKey, nil), // Add JWT auth
)
```

**Order of Execution:**
Middleware is applied in reverse order, meaning the first middleware in the list is the outermost in the chain and executes first.

### CORS

The `cors` package implements Cross-Origin Resource Sharing middleware with comprehensive configuration options.

**Key Features:**
- Configurable CORS headers
- Support for preflight requests
- Origin validation
- Method and header controls
- Environment-aware behavior (different defaults for development/production)

**Main Functions:**
- `Cors(config *CORSConfig, next http.HandlerFunc) http.HandlerFunc`: Middleware that handles Cross-Origin Resource Sharing headers for HTTP requests

**Configuration Options:**
```go
type CORSConfig struct {
    AllowOrigin      string   `mapstructure:"allow_origin,omitempty"`
    AllowedOrigins   []string `mapstructure:"allowed_origins,omitempty"`
    AllowMethods     string   `mapstructure:"allow_methods,omitempty"`
    AllowHeaders     string   `mapstructure:"allow_headers,omitempty"`
    AllowCredentials bool     `mapstructure:"allow_credentials,omitempty"`
    ExposeHeaders    string   `mapstructure:"expose_headers,omitempty"`
    MaxAge           int      `mapstructure:"max_age,omitempty"`
    ExemptMethods    []string `mapstructure:"exempt_methods,omitempty"`
}
```

**Configuration Details:**
- `AllowOrigin`: Specifies the allowed origin for CORS requests. Default is "*".
- `AllowedOrigins`: List of specific origins that are allowed.
- `AllowMethods`: Specifies the allowed HTTP methods. Default is "*".
- `AllowHeaders`: Specifies the allowed headers. Default is "*".
- `AllowCredentials`: Whether to allow credentials. Default is true.
- `ExposeHeaders`: Headers that should be exposed to the client. Default is "*".
- `MaxAge`: How long the results of a preflight request can be cached.
- `ExemptMethods`: Methods that are exempt from CORS checks (typically OPTIONS).

**Default Behavior:**
- In production mode, more restrictive defaults are applied
- OPTIONS requests are automatically handled
- Preflight responses are cached based on MaxAge setting

### Logger

The `logger` package provides structured logging for HTTP requests and responses using zerolog.

**Key Features:**
- Request/response logging with timing
- Structured JSON output using zerolog
- Configurable log levels
- Request ID generation and tracking
- Environment-aware logging (more verbose in development)

**Main Functions:**
- `Log(next http.HandlerFunc) http.HandlerFunc`: Middleware that wraps handlers with request/response logging functionality

**Context Integration:**
- Adds request ID to context for traceability using `UUIDKey`
- Provides logger instance in context for use in handlers through the `log` subpackage
- Request UUIDs are generated automatically and can be accessed via context

**Logging Details:**
- Captures request method, URL, headers, and timing
- Logs response status codes and processing time
- Includes error details when applicable
- In debug mode, includes stack traces for errors

### Method

The `method` package provides HTTP method filtering functionality to allow or deny specific HTTP methods.

**Key Features:**
- Allow specific HTTP methods to pass through
- Deny specific HTTP methods with automatic 405 responses
- Flexible method restriction for different endpoints

**Main Functions:**
- `AllowMethods(methods ...string) HandlerFunc`: Creates a middleware that only allows specified HTTP methods
- `DenyMethods(methods ...string) HandlerFunc`: Creates a middleware that blocks specified HTTP methods

**Usage Examples:**
```go
// Only allow GET and POST methods
handler := possum.Chain(
    myHandler,
    possum.AllowMethods("GET", "POST"),
)

// Block DELETE and PUT methods
handler := possum.Chain(
    myHandler,
    possum.DenyMethods("DELETE", "PUT"),
)
```

**Response Behavior:**
- When a method is not allowed, automatically returns a 405 Method Not Allowed response
- Includes the `Allow` header listing permitted methods
- Uses the predefined `MethodNotAllowedResponse` for consistent error formatting

### Response

The `response` package provides standardized JSON response formatting with UUID tracking and error handling.

**Key Features:**
- Consistent JSON response structure with UUID tracking
- Standard HTTP status codes
- Predefined error responses for common scenarios
- Error response formatting with optional stack traces
- Data serialization

**Response Structure:**
```go
type Response struct {
    UUID  uuid.UUID `json:"uuid"`
    Data  any       `json:"data,omitempty"`
    Error *Error    `json:"error,omitempty"`
}
```

**Main Functions:**
- `NewResponse(r *http.Request) *Response`: Creates a new Response with UUID from request context
- `CloneResponse(r Response) Response`: Creates a deep copy of a Response
- `WriteResponse(w http.ResponseWriter, resp Response, err error)`: Writes a Response with error handling
- `SetData(data any)`: Configures the Response with data
- `SetError(code int, message string)`: Configures the Response with error details
- `WriteHeader(code int)`: Sets the HTTP status code
- `Write(w http.ResponseWriter)`: Serializes the Response to the HTTP response writer

**Predefined Error Responses:**
- `InternalServerErrorResponse` (HTTP 500)
- `UnauthorizedResponse` (HTTP 401)
- `MethodNotAllowedResponse` (HTTP 405)
- `BadRequestResponse` (HTTP 400)
- `NotFoundResponse` (HTTP 404)
- `NotImplementedResponse` (HTTP 501)
- `ConflictResponse` (HTTP 409)
- `ForbiddenResponse` (HTTP 403)

**Error Structure:**
```go
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Stack   []byte `json:"stack,omitempty"` // Included in debug mode
}
```

**Usage Example:**
```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    resp := possum.NewResponse(r)
    
    // Some operation that might fail
    data, err := getData()
    if err != nil {
        resp.SetError(http.StatusInternalServerError, "Failed to retrieve data")
        resp.Write(w)
        return
    }
    
    resp.SetData(data)
    resp.Write(w)
}
```

### WebSocket

The `websocket` package provides utilities for handling WebSocket connections with built-in authentication and CORS support.

**Key Features:**
- WebSocket connection upgrade with CORS handling
- JWT authentication for WebSocket connections
- Automatic ping/pong handling for connection health
- Graceful connection cleanup
- Built-in error handling

**Main Functions:**
- `WebSocketUpgrade(corsConfig *CORSConfig, next WebsocketHandlerFunc) http.HandlerFunc`: Middleware that handles WebSocket connections with CORS support
- `WebSocketAuth(secret []byte, next WebsocketHandlerFunc) WebsocketHandlerFunc`: Middleware that wraps WebSocket handlers with JWT authentication

**WebSocket Handler Type:**
```go
type WebsocketHandlerFunc func(conn *websocket.Conn, r *http.Request)
```

**Connection Management:**
- Automatic handling of ping/pong messages to maintain connection health
- Proper connection cleanup when handlers complete
- Hijacker interface implementation for WebSocket upgrade

**Usage Example:**
```go
wsHandler := func(conn *websocket.Conn, r *http.Request) {
    // Read message from client
    _, message, err := conn.ReadMessage()
    if err != nil {
        log.Printf("Error reading message: %v", err)
        return
    }
    
    // Echo message back
    err = conn.WriteMessage(websocket.TextMessage, message)
    if err != nil {
        log.Printf("Error writing message: %v", err)
        return
    }
}

http.HandleFunc("/ws", possum.WebSocketUpgrade(nil, wsHandler))
```

**Dependencies:**
- Uses `github.com/gorilla/websocket` for underlying WebSocket implementation

## Usage Examples

### Basic Server with Middleware

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
        possum.Log,           // Add request logging
        possum.Cors(nil),     // Add CORS with default config
    ))

    http.ListenAndServe(":8080", nil)
}
```

### Protected Route with JWT Authentication

```go
// Protect a route with JWT authentication
protectedHandler := possum.Chain(
    myHandler,
    possum.Log,
    possum.HTTPAuth(secretKey, nil), // Protect with JWT auth
)

http.HandleFunc("/protected", protectedHandler)
```

### WebSocket Endpoint

```go
wsHandler := func(conn *websocket.Conn, r *http.Request) {
    // Handle WebSocket communication
    // Ping/pong and connection cleanup are handled automatically
}

http.HandleFunc("/ws", possum.WebSocketUpgrade(
    nil, 
    wsHandler,
))
```

### Method Restriction

```go
// Only allow GET and POST methods
handler := possum.Chain(
    myHandler,
    possum.AllowMethods("GET", "POST"),
)

http.HandleFunc("/api/data", handler)
```

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

## Testing

Each component includes comprehensive tests that demonstrate usage patterns with >80% coverage:

1. **Unit Tests**: Each module has dedicated test files (e.g., `auth_test.go`)
2. **Integration Patterns**: Tests show how components work together
3. **Example Usage**: Test functions serve as documentation for real-world usage
4. **Mock Implementations**: Tests include mock implementations for external dependencies

**Running Tests:**
```bash
./run_tests.sh
# or
go test -v ./...
```

**Test Structure:**
- Test functions follow the pattern `Test<Component>_<Functionality>`
- Use table-driven tests where appropriate
- Include both success and error cases
- Demonstrate proper usage patterns
- Include tests for edge cases and error conditions

**Test Features:**
- Tests validate middleware behavior in isolation
- Integration tests verify chained middleware functionality
- Mock HTTP servers for testing WebSocket functionality
- Test coverage reports for quality assurance

## Dependencies

The possum package relies on several well-maintained third-party libraries:

- `github.com/golang-jwt/jwt/v5`: JWT implementation for authentication
- `github.com/google/uuid`: UUID generation for request IDs
- `github.com/gorilla/websocket`: WebSocket protocol implementation
- `github.com/rs/zerolog`: High-performance logging library

### Subpackages Dependencies

The possum package also includes subpackages with their own dependencies:

- `config/`: Environment-based configuration utilities
- `log/`: Enhanced logging wrapper around zerolog with configuration options

These dependencies are carefully selected for their stability, performance, and minimal footprint.

**Go Version Requirement:** Go 1.23.5 or higher

## Development Guidelines

1. **Middleware Pattern**: All middleware should follow the signature `func(http.HandlerFunc) http.HandlerFunc`
2. **Context Usage**: Use request context to pass data between middleware and handlers, with standardized context keys
3. **Error Handling**: Handle errors gracefully and return appropriate HTTP status codes using predefined response structures
4. **Testing**: Each component should have comprehensive tests with >80% coverage
5. **Documentation**: All exported functions should have clear documentation comments
6. **Environment Awareness**: Respect environment variables like `POSSUM_ENV` for behavior changes
7. **Security**: Validate inputs and sanitize outputs to prevent common web vulnerabilities
8. **Performance**: Minimize allocations and optimize for high-throughput scenarios
9. **Compatibility**: Maintain backward compatibility when possible, follow semantic versioning

## Constants and Keys

The package provides standardized context keys for data passing:

- `UUIDKey`: Key for storing request UUIDs in context
- `ClaimsKey`: Key for storing JWT claims in context

## Installation

```bash
go get github.com/mikespook/possum
```