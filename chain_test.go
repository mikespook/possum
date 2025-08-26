package possum

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestChain tests the Chain function for composing multiple middleware handlers in the correct order.
func TestChain(t *testing.T) {
	// Create a counter to track middleware execution
	var executionOrder []string

	// Create a base handler
	baseHandler := func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "base")
		w.WriteHeader(http.StatusOK)
	}

	// Create middleware functions
	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		}
	}

	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		}
	}

	middleware3 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware3-before")
			next(w, r)
			executionOrder = append(executionOrder, "middleware3-after")
		}
	}

	// Test cases
	tests := []struct {
		name             string
		middlewares      []HandlerFunc
		expectedOrder    []string
		expectedStatus   int
	}{
		{
			name:             "No middleware",
			middlewares:      []HandlerFunc{},
			expectedOrder:    []string{"base"},
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Single middleware",
			middlewares:      []HandlerFunc{middleware1},
			expectedOrder:    []string{"middleware1-before", "base", "middleware1-after"},
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Multiple middlewares",
			middlewares:      []HandlerFunc{middleware1, middleware2, middleware3},
			expectedOrder:    []string{
				"middleware1-before", 
				"middleware2-before", 
				"middleware3-before", 
				"base", 
				"middleware3-after", 
				"middleware2-after", 
				"middleware1-after",
			},
			expectedStatus:   http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset execution order
			executionOrder = []string{}

			// Create a request
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			// Chain the middlewares
			handler := Chain(baseHandler, tc.middlewares...)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check execution order
			if len(executionOrder) != len(tc.expectedOrder) {
				t.Errorf("Expected %d executions, got %d", len(tc.expectedOrder), len(executionOrder))
			}

			for i, step := range tc.expectedOrder {
				if i >= len(executionOrder) || executionOrder[i] != step {
					t.Errorf("Expected step %d to be %s, got %s", i, step, 
						func() string {
							if i >= len(executionOrder) {
								return "none"
							}
							return executionOrder[i]
						}())
				}
			}
		})
	}
}