package possum

import (
	"net/http"
)

// AllowMethods creates a middleware that only allows specified HTTP methods to pass through.
func AllowMethods(methods ...string) HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, method := range methods {
				if r.Method == method {
					next(w, r)
					return
				}
			}
			MethodNotAllowedResponse.Write(w)
		}
	}
}

// DenyMethods creates a middleware that blocks specified HTTP methods from passing through.
func DenyMethods(methods ...string) HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, method := range methods {
				if r.Method == method {
					MethodNotAllowedResponse.Write(w)
					return
				}
			}
			next(w, r)
		}
	}
}
