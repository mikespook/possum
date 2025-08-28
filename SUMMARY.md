# Possum - HTTP Middleware Toolkit - Summary

Possum is a Go package that provides a collection of HTTP middleware functions for building robust web applications and APIs. It offers essential functionality for authentication, CORS handling, logging, method filtering, WebSocket support, and standardized response formatting.

## Key Features

1. **Middleware Chain** - Compose multiple middleware handlers in a clean, predictable order
2. **Authentication** - JWT-based authentication for both HTTP and WebSocket connections
3. **CORS Handling** - Comprehensive Cross-Origin Resource Sharing support with configurable policies
4. **Logging** - Structured request/response logging using zerolog
5. **Method Filtering** - Allow or deny specific HTTP methods
6. **WebSocket Support** - WebSocket upgrade handler with built-in CORS and connection management
7. **Standardized Responses** - Consistent JSON response format with UUID tracking and error handling

## Architecture

The package is organized into several key components:

- `possum` - Main package with core middleware functions
- `possum/auth` - JWT token generation and parsing
- `possum/config` - Environment-based configuration
- `possum/log` - Structured logging with zerolog
- `possum/utils` - Utility functions for body tracing

## Usage Examples

### Simple HTTP Server with Middleware

```go
handler := possum.Chain(
    myHandler,
    possum.Log,           // Add request logging
    possum.Cors(nil),     // Add CORS with default config
)
```

### Protected Route with JWT Authentication

```go
protectedHandler := possum.Chain(
    myHandler,
    possum.Log,
    possum.HTTPAuth(secretKey, nil), // Protect with JWT auth
)
```

### WebSocket Endpoint

```go
http.HandleFunc("/ws", possum.WebSocketUpgrade(
    corsConfig, 
    myWebSocketHandler,
))
```

## Testing

The package includes comprehensive tests for all components with >80% coverage. Run tests with:

```bash
./run_tests.sh
```

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `github.com/google/uuid` - UUID generation
- `github.com/gorilla/websocket` - WebSocket implementation
- `github.com/rs/zerolog` - Structured logging

## Documentation

See the following files for detailed documentation:
- `README.md` - General overview and usage
- `API.md` - Detailed API documentation
- `STRUCTURE.md` - Project structure and organization