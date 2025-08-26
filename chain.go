package possum

import "net/http"

type HandlerFunc func(http.HandlerFunc) http.HandlerFunc

// Chain composes multiple middleware handlers into a single handler, applying them in reverse order.
func Chain(handler http.HandlerFunc, middlewares ...HandlerFunc) http.HandlerFunc {
	for i := range middlewares {
		handler = middlewares[len(middlewares)-1-i](handler)
	}
	return handler
}
