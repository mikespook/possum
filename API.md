# Possum API Documentation

This document provides detailed API documentation for the Possum package, a collection of HTTP middleware functions for building robust web applications in Go.

## Table of Contents

1. [Package Overview](#package-overview)
2. [Types](#types)
3. [Functions](#functions)
4. [Constants](#constants)
5. [Variables](#variables)
6. [Examples](#examples)

## Package Overview

Package possum provides a collection of HTTP middleware functions for building robust web applications. It includes authentication, CORS handling, logging, method filtering, WebSocket support, and standardized response formatting.

The package is designed to be modular and composable, allowing developers to chain middleware functions together to create customized request processing pipelines.

## Types

### CORSConfig

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

CORSConfig holds configuration for Cross-Origin Resource Sharing (CORS) headers.

**Fields:**
- `AllowOrigin`: Specifies the allowed origin for CORS requests. Default is "*".
- `AllowedOrigins`: List of specific origins that are allowed.
- `AllowMethods`: Specifies the allowed HTTP methods. Default is "*".
- `AllowHeaders`: Specifies the allowed headers. Default is "*".
- `AllowCredentials`: Whether to allow credentials. Default is true.
- `ExposeHeaders`: Headers that should be exposed to the client. Default is "*".
- `MaxAge`: How long the results of a preflight request can be cached.
- `ExemptMethods`: Methods that are exempt from CORS checks (typically OPTIONS).

### Error

```go
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Stack   []byte `json:"stack,omitempty"`
}
```

Error represents an error response with optional stack trace in debug mode.

### HandlerFunc

```go
type HandlerFunc func(http.HandlerFunc) http.HandlerFunc
```

HandlerFunc is a type for middleware functions that wrap http.HandlerFunc.

### Response

```go
type Response struct {
    UUID  uuid.UUID `json:"uuid"`
    Data  any       `json:"data,omitempty"`
    Error *Error    `json:"error,omitempty"`
    // contains filtered or unexported fields
}
```

Response represents a standardized JSON response structure with UUID tracking.

### WebsocketHandlerFunc

```go
type WebsocketHandlerFunc func(conn *websocket.Conn, r *http.Request)
```

WebsocketHandlerFunc is a type for WebSocket handlers.

## Functions

### AllowMethods

```go
func AllowMethods(methods ...string) HandlerFunc
```

AllowMethods creates a middleware that only allows specified HTTP methods to pass through. If the request method is not in the allowed list, it returns a 405 Method Not Allowed response.

**Parameters:**
- `methods`: List of HTTP methods to allow (e.g., "GET", "POST")

**Returns:**
- A HandlerFunc middleware

### Chain

```go
func Chain(handler http.HandlerFunc, middlewares ...HandlerFunc) http.HandlerFunc
```

Chain composes multiple middleware handlers into a single handler, applying them in reverse order. This allows for clean middleware composition where the first middleware in the list is the outermost in the chain.

**Parameters:**
- `handler`: The final http.HandlerFunc to be called
- `middlewares`: List of HandlerFunc middleware to apply

**Returns:**
- A composed http.HandlerFunc

### Cors

```go
func Cors(config *CORSConfig, next http.HandlerFunc) http.HandlerFunc
```

Cors is a middleware that handles Cross-Origin Resource Sharing (CORS) headers for HTTP requests. It manages preflight requests, origin validation, and setting appropriate CORS headers.

**Parameters:**
- `config`: CORS configuration. If nil, uses default configuration.
- `next`: The next http.HandlerFunc in the chain

**Returns:**
- An http.HandlerFunc with CORS support

### DenyMethods

```go
func DenyMethods(methods ...string) HandlerFunc
```

DenyMethods creates a middleware that blocks specified HTTP methods from passing through. If the request method is in the denied list, it returns a 405 Method Not Allowed response.

**Parameters:**
- `methods`: List of HTTP methods to deny

**Returns:**
- A HandlerFunc middleware

### HTTPAuth

```go
func HTTPAuth(secret []byte, next http.HandlerFunc) http.HandlerFunc
```

HTTPAuth is a middleware that wraps an http.HandlerFunc with JWT authentication logic. It validates the Bearer token in the Authorization header and adds the parsed claims to the request context.

**Parameters:**
- `secret`: Secret key used to sign/verify JWT tokens
- `next`: The next http.HandlerFunc in the chain

**Returns:**
- An http.HandlerFunc with JWT authentication

### Log

```go
func Log(next http.HandlerFunc) http.HandlerFunc
```

Log is a middleware that wraps an http.HandlerFunc with request/response logging functionality using zerolog.

**Parameters:**
- `next`: The next http.HandlerFunc in the chain

**Returns:**
- An http.HandlerFunc with logging

### NewResponse

```go
func NewResponse(r *http.Request) *Response
```

NewResponse creates a new Response object with a UUID from the request context. If no UUID exists in the context, a new one is generated.

**Parameters:**
- `r`: The http.Request

**Returns:**
- A pointer to a new Response

### WebSocketAuth

```go
func WebSocketAuth(secret []byte, next WebsocketHandlerFunc) WebsocketHandlerFunc
```

WebSocketAuth is a middleware that wraps a WebsocketHandlerFunc with JWT authentication logic. It validates the token from the "token" query parameter.

**Parameters:**
- `secret`: Secret key used to sign/verify JWT tokens
- `next`: The next WebsocketHandlerFunc in the chain

**Returns:**
- A WebsocketHandlerFunc with JWT authentication

### WebSocketUpgrade

```go
func WebSocketUpgrade(corsConfig *CORSConfig, next WebsocketHandlerFunc) http.HandlerFunc
```

WebSocketUpgrade is a middleware that handles WebSocket connections with CORS support. It upgrades HTTP connections to WebSocket connections and manages the connection lifecycle including ping/pong handling.

**Parameters:**
- `corsConfig`: CORS configuration for WebSocket connections
- `next`: The next WebsocketHandlerFunc in the chain

**Returns:**
- An http.HandlerFunc that handles WebSocket upgrades

## Response Methods

### CloneResponse

```go
func CloneResponse(r Response) Response
```

CloneResponse creates a deep copy of a Response object, including error details.

**Parameters:**
- `r`: Response to clone

**Returns:**
- A cloned Response

### WriteResponse

```go
func WriteResponse(w http.ResponseWriter, resp Response, err error)
```

WriteResponse writes a Response object to an HTTP response writer, handling errors appropriately. If err is not nil, it updates the response with error details.

**Parameters:**
- `w`: http.ResponseWriter to write to
- `resp`: Response to write
- `err`: Error to handle

### SetData

```go
func (resp *Response) SetData(data any)
```

SetData configures the Response object with response data.

**Parameters:**
- `data`: Data to set in the response

### SetError

```go
func (resp *Response) SetError(code int, message string)
```

SetError configures the Response object with error details, including stack trace in debug mode.

**Parameters:**
- `code`: HTTP status code
- `message`: Error message

### Write

```go
func (resp *Response) Write(w http.ResponseWriter)
```

Write serializes the Response object to the HTTP response writer with proper headers.

**Parameters:**
- `w`: http.ResponseWriter to write to

### WriteHeader

```go
func (resp *Response) WriteHeader(code int)
```

WriteHeader sets the HTTP status code for the Response object.

**Parameters:**
- `code`: HTTP status code

## CORS Configuration Methods

### Init

```go
func (config *CORSConfig) Init()
```

Init initializes the CORS configuration by caching exempt methods for faster lookup.

### SkipMethod

```go
func (config *CORSConfig) SkipMethod(method string) bool
```

SkipMethod checks if a method should be exempt from CORS processing.

**Parameters:**
- `method`: HTTP method to check

**Returns:**
- true if the method should be exempt, false otherwise

## Constants

```go
const (
    // UUIDKey is the key used to store the UUID in the context
    UUIDKey ContextKey = "uuid"
    
    // ClaimsKey is the key used to store JWT claims in the context
    ClaimsKey ContextKey = "claims"
)
```

## Variables

### Error Responses

```go
var (
    // Predefined error responses for common HTTP status codes.
    InternalServerErrorResponse = Response{...}
    UnauthorizedResponse = Response{...}
    MethodNotAllowedResponse = Response{...}
    BadRequestResponse = Response{...}
    NotFoundResponse = Response{...}
    NotImplementedResponse = Response{...}
    ConflictResponse = Response{...}
    ForbiddenResponse = Response{...}
)
```

Predefined Response objects for common HTTP error statuses.

### Errors

```go
var (
    ErrUnauthorized = errors.New("Unauthorized")
    ErrHijackerNotImplement = fmt.Errorf("underlying ResponseWriter does not implement Hijacker")
)
```

Common error values used throughout the package.

## Examples

### Chaining Middleware

```go
handler := possum.Chain(
    myHandler,
    possum.Log,
    possum.Cors(nil),
    possum.HTTPAuth(secret, nil),
)
```

### Creating Custom Responses

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    resp := possum.NewResponse(r)
    resp.SetData(map[string]string{"message": "Hello, World!"})
    resp.WriteHeader(http.StatusCreated)
    resp.Write(w)
}
```

### Handling Errors

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

### WebSocket Handler

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