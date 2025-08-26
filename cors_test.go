package possum

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCORSConfigInit tests the Init method of CORSConfig for initializing CORS settings.
func TestCORSConfigInit(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		config         CORSConfig
		exemptMethods  []string
		methodToCheck  string
		expectedResult bool
	}{
		{
			name: "Empty exempt methods",
			config: CORSConfig{
				ExemptMethods: []string{},
			},
			methodToCheck:  "OPTIONS",
			expectedResult: false,
		},
		{
			name: "Single exempt method",
			config: CORSConfig{
				ExemptMethods: []string{"OPTIONS"},
			},
			methodToCheck:  "OPTIONS",
			expectedResult: true,
		},
		{
			name: "Multiple exempt methods",
			config: CORSConfig{
				ExemptMethods: []string{"OPTIONS", "HEAD"},
			},
			methodToCheck:  "HEAD",
			expectedResult: true,
		},
		{
			name: "Method not exempt",
			config: CORSConfig{
				ExemptMethods: []string{"OPTIONS", "HEAD"},
			},
			methodToCheck:  "GET",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the config
			tc.config.Init()

			// Check if method is exempt
			result := tc.config.SkipMethod(tc.methodToCheck)

			// Verify result
			if result != tc.expectedResult {
				t.Errorf("Expected SkipMethod(%s) to be %v, got %v", 
					tc.methodToCheck, tc.expectedResult, result)
			}
		})
	}
}

// TestIsOriginAllowed tests the isOriginAllowed function for validating request origins against CORS rules.
func TestIsOriginAllowed(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		config         *CORSConfig
		origin         string
		expectedResult bool
	}{
		{
			name: "Empty origin",
			config: &CORSConfig{
				AllowOrigin: "*",
			},
			origin:         "",
			expectedResult: false,
		},
		{
			name: "Wildcard allow origin without credentials",
			config: &CORSConfig{
				AllowOrigin:      "*",
				AllowCredentials: false,
			},
			origin:         "http://example.com",
			expectedResult: true,
		},
		{
			name: "Wildcard allow origin with credentials",
			config: &CORSConfig{
				AllowOrigin:      "*",
				AllowCredentials: true,
			},
			origin:         "http://example.com",
			expectedResult: false,
		},
		{
			name: "Specific origin in allowed list",
			config: &CORSConfig{
				AllowedOrigins: []string{"http://example.com", "http://test.com"},
			},
			origin:         "http://example.com",
			expectedResult: true,
		},
		{
			name: "Origin not in allowed list",
			config: &CORSConfig{
				AllowedOrigins: []string{"http://example.com", "http://test.com"},
			},
			origin:         "http://other.com",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test the function
			result := isOriginAllowed(tc.config, tc.origin)

			// Verify result
			if result != tc.expectedResult {
				t.Errorf("Expected isOriginAllowed to be %v, got %v", 
					tc.expectedResult, result)
			}
		})
	}
}

// TestCors tests the Cors middleware function for handling Cross-Origin Resource Sharing (CORS) headers.
func TestCors(t *testing.T) {
	// Create a mock handler
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Test cases
	tests := []struct {
		name               string
		config             *CORSConfig
		method             string
		origin             string
		expectedStatus     int
		expectedHeaders    map[string]string
		expectedHeadersSet bool
	}{
		{
			name:           "Default config with OPTIONS method",
			config:         nil, // Use default config
			method:         "OPTIONS",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "*",
				"Access-Control-Allow-Headers":     "*",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":    "*",
				"Vary":                             "Origin",
			},
			expectedHeadersSet: true,
		},
		{
			name: "Custom config with GET method",
			config: &CORSConfig{
				AllowOrigin:      "http://allowed.com",
				AllowedOrigins:   []string{"http://allowed.com"},
				AllowMethods:     "GET, POST",
				AllowHeaders:     "Content-Type, Authorization",
				AllowCredentials: true,
				ExposeHeaders:    "X-Custom-Header",
				MaxAge:           3600,
				ExemptMethods:    []string{"OPTIONS"},
			},
			method:         "GET",
			origin:         "http://allowed.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://allowed.com",
				"Access-Control-Allow-Methods":     "GET, POST",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":    "X-Custom-Header",
				"Access-Control-Max-Age":           "3600",
				"Vary":                             "Origin",
			},
			expectedHeadersSet: true,
		},
		{
			name: "Origin not allowed",
			config: &CORSConfig{
				AllowedOrigins: []string{"http://allowed.com"},
				AllowMethods:   "GET, POST",
				ExemptMethods:  []string{"OPTIONS"},
			},
			method:             "GET",
			origin:             "http://notallowed.com",
			expectedStatus:     http.StatusOK,
			expectedHeadersSet: false,
		},
		{
			name: "WebSocket upgrade request",
			config: &CORSConfig{
				AllowOrigin:    "*",
				AllowHeaders:   "Content-Type",
				ExemptMethods:  []string{"OPTIONS"},
			},
			method:         "GET",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "http://example.com",
				"Access-Control-Allow-Headers": "Content-Type, Sec-WebSocket-Key, Sec-WebSocket-Protocol, Sec-WebSocket-Version",
				"Vary":                         "Origin",
			},
			expectedHeadersSet: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize config if provided
			if tc.config != nil {
				tc.config.Init()
			}

			// Create a request
			req := httptest.NewRequest(tc.method, "/", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			
			// For WebSocket test
			if tc.name == "WebSocket upgrade request" {
				req.Header.Set("Upgrade", "websocket")
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create the CORS middleware
			handler := Cors(tc.config, mockHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check headers
			if tc.expectedHeadersSet {
				for key, expectedValue := range tc.expectedHeaders {
					actualValue := rr.Header().Get(key)
					if actualValue != expectedValue {
						t.Errorf("Expected header %s to be %s, got %s", 
							key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}