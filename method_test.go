package possum

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAllowMethods tests the AllowMethods middleware function for restricting HTTP methods.
func TestAllowMethods(t *testing.T) {
	// Create a mock handler
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	// Test cases
	tests := []struct {
		name           string
		allowedMethods []string
		requestMethod  string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Allowed method - GET",
			allowedMethods: []string{"GET", "POST"},
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Allowed method - POST",
			allowedMethods: []string{"GET", "POST"},
			requestMethod:  "POST",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Not allowed method",
			allowedMethods: []string{"GET", "POST"},
			requestMethod:  "DELETE",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "", // Method not allowed response is handled by the middleware
		},
		{
			name:           "Empty allowed methods",
			allowedMethods: []string{},
			requestMethod:  "GET",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "", // Method not allowed response is handled by the middleware
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified method
			req := httptest.NewRequest(tc.requestMethod, "/", nil)
			
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create the middleware
			middleware := AllowMethods(tc.allowedMethods...)
			handler := middleware(mockHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check response body if expected
			if tc.expectedBody != "" && rr.Body.String() != tc.expectedBody {
				t.Errorf("Expected body %s, got %s", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

// TestDenyMethods tests the DenyMethods middleware function for blocking specific HTTP methods.
func TestDenyMethods(t *testing.T) {
	// Create a mock handler
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	// Test cases
	tests := []struct {
		name           string
		deniedMethods  []string
		requestMethod  string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Denied method",
			deniedMethods:  []string{"DELETE", "PUT"},
			requestMethod:  "DELETE",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "", // Method not allowed response is handled by the middleware
		},
		{
			name:           "Not denied method",
			deniedMethods:  []string{"DELETE", "PUT"},
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Empty denied methods",
			deniedMethods:  []string{},
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Multiple denied methods",
			deniedMethods:  []string{"DELETE", "PUT", "PATCH"},
			requestMethod:  "PUT",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "", // Method not allowed response is handled by the middleware
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified method
			req := httptest.NewRequest(tc.requestMethod, "/", nil)
			
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create the middleware
			middleware := DenyMethods(tc.deniedMethods...)
			handler := middleware(mockHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check response body if expected
			if tc.expectedBody != "" && rr.Body.String() != tc.expectedBody {
				t.Errorf("Expected body %s, got %s", tc.expectedBody, rr.Body.String())
			}
		})
	}
}