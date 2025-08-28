# Possum Package Structure

This document provides an overview of the Possum package structure and its key components.

## Project Structure

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

## Main Package Files

### Core Middleware
1. `chain.go` - Middleware composition
2. `cors.go` - CORS handling
3. `logger.go` - Request/response logging
4. `method.go` - HTTP method filtering
5. `auth.go` - JWT authentication
6. `websocket.go` - WebSocket support
7. `response.go` - Standardized response handling

### Test Files
Each component has corresponding test files following the `*_test.go` naming convention.

## Key Components

### 1. Middleware Chain (`chain.go`)
- `Chain()` - Composes multiple middleware into a single handler
- Applies middleware in reverse order (last added executes first)

### 2. Authentication (`auth.go`)
- `HTTPAuth()` - JWT authentication for HTTP requests
- `WebSocketAuth()` - JWT authentication for WebSocket connections
- Uses context to pass claims to handlers

### 3. CORS Handling (`cors.go`)
- `Cors()` - Cross-Origin Resource Sharing middleware
- Configurable through `CORSConfig` struct
- Handles preflight requests automatically

### 4. Logging (`logger.go`)
- `Log()` - Request/response logging middleware
- Uses zerolog for structured logging
- Captures request details and response status

### 5. Method Filtering (`method.go`)
- `AllowMethods()` - Restricts allowed HTTP methods
- `DenyMethods()` - Blocks specific HTTP methods

### 6. WebSocket Support (`websocket.go`)
- `WebSocketUpgrade()` - Upgrades HTTP connections to WebSocket
- Built-in ping/pong handling
- Automatic connection cleanup

### 7. Response Handling (`response.go`)
- `Response` - Standardized JSON response structure
- Predefined error responses
- UUID tracking for requests
- Stack trace support in debug mode

## Subpackages

### Authentication (`auth/`)
- `JWTClaims` - JWT token claims structure
- `GenerateJWT()` - Creates signed JWT tokens
- `ParseToken()` - Validates and parses JWT tokens

### Configuration (`config/`)
- `IsDebug()` - Checks if running in debug mode
- `IsDev()` - Checks if running in development mode
- Environment-based configuration using `POSSUM_ENV`

### Logging (`log/`)
- Wrapper around zerolog
- Configurable log levels and output destinations
- Global logger with context support

## Testing

Each component includes comprehensive tests:
- Unit tests for individual functions
- Integration tests for middleware chains
- Mock implementations for external dependencies

Run tests with:
```bash
./run_tests.sh
```

## Dependencies

1. `github.com/golang-jwt/jwt/v5` - JWT implementation
2. `github.com/google/uuid` - UUID generation
3. `github.com/gorilla/websocket` - WebSocket support
4. `github.com/rs/zerolog` - Structured logging

## Environment Variables

- `POSSUM_ENV` - Controls environment (development, production, test)
  - Affects debug logging and stack traces
  - Influences WebSocket CORS behavior

## Usage Patterns

### 1. Simple Middleware Chain
```go
handler := possum.Chain(
    myHandler,
    possum.Log,
    possum.Cors(nil),
)
```

### 2. Protected Routes
```go
protectedHandler := possum.Chain(
    myHandler,
    possum.Log,
    possum.HTTPAuth(secret, nil),
)
```

### 3. WebSocket Endpoint
```go
http.HandleFunc("/ws", possum.WebSocketUpgrade(
    corsConfig, 
    myWebSocketHandler,
))
```

### 4. Method Restriction
```go
handler := possum.Chain(
    myHandler,
    possum.AllowMethods("GET", "POST"),
)
```

## Error Handling

The package provides standardized error responses:
- `UnauthorizedResponse` - HTTP 401
- `MethodNotAllowedResponse` - HTTP 405
- `BadRequestResponse` - HTTP 400
- `NotFoundResponse` - HTTP 404
- `InternalServerErrorResponse` - HTTP 500

Handlers can also create custom responses using:
- `NewResponse()` - Creates a new response with UUID
- `Response.SetError()` - Sets error details
- `Response.SetData()` - Sets response data