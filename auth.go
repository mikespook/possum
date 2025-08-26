package possum

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/mikespook/possum/auth"
)

type ContextKey string

const (
	ClaimsKey = ContextKey("claims")
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

// HTTPAuth is a middleware that wraps an http.HandlerFunc with JWT authentication logic.
func HTTPAuth(secret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			UnauthorizedResponse.Write(w)
			return
		}
		// Check if the header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			UnauthorizedResponse.Write(w)
			return
		}
		token := parts[1]

		// Validate JWT token
		claims, err := auth.ParseToken(secret, token)
		if err != nil {
			UnauthorizedResponse.Write(w)
			return
		}
		// Call the next handler
		next(w, r.WithContext(context.WithValue(r.Context(), ClaimsKey, claims)))
	}
}

// WebSocketAuth is a middleware that wraps a WebsocketHandlerFunc with JWT authentication logic.
func WebSocketAuth(secret []byte, next WebsocketHandlerFunc) WebsocketHandlerFunc {
	return func(conn *websocket.Conn, r *http.Request) {
		token := r.URL.Query().Get("token")
		claims, err := auth.ParseToken(secret, token)
		if err != nil {
			conn.WriteMessage(websocket.CloseMessage, []byte("Invalid token"))
			return
		}
		// Call the next handler
		next(conn, r.WithContext(context.WithValue(r.Context(), ClaimsKey, claims)))
		// Call the next handler
	}
}
