# Possum - Go HTTP Middleware Toolkit

Possum is a lightweight, modular Go package that provides essential HTTP middleware functions for building robust web applications and APIs. It offers a collection of independent components that can be used together or separately.

## Key Features

- **Authentication**: JWT-based authentication for HTTP and WebSocket connections
- **CORS Handling**: Comprehensive Cross-Origin Resource Sharing support
- **Logging**: Structured request/response logging with zerolog
- **Method Filtering**: Allow or deny specific HTTP methods
- **WebSocket Support**: WebSocket upgrade handler with connection management
- **Response Formatting**: Standardized JSON responses with UUID tracking
- **Middleware Chaining**: Compose multiple middleware handlers in a clean, predictable order

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
        possum.Log,      // Add logging
        possum.Cors(nil), // Add CORS with default config
    ))

    http.ListenAndServe(":8080", nil)
}
```

## Core Components

1. **Middleware Chain** - Compose multiple middleware handlers into a single handler
2. **Authentication** - JWT-based authentication for both HTTP and WebSocket connections
3. **CORS Handling** - Cross-Origin Resource Sharing middleware with configurable policies
4. **Logging** - Structured request/response logging using zerolog
5. **Method Filtering** - Allow or deny specific HTTP methods
6. **WebSocket Support** - WebSocket upgrade handler with built-in connection management
7. **Response Handling** - Consistent JSON response format with UUID tracking

## Documentation

For comprehensive documentation, please refer to:
- [LLM_CODER_GUIDE.md](LLM_CODER_GUIDE.md) - Complete guide for developers and AI models
- [API.md](API.md) - Detailed API documentation
- [STRUCTURE.md](STRUCTURE.md) - Project structure and organization

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

## License

MIT