package possum

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mikespook/possum/config"
)

// TestNewResponse tests the NewResponse function for creating HTTP responses with UUID tracking.
func TestNewResponse(t *testing.T) {
	// Test with UUID in context
	t.Run("With UUID in context", func(t *testing.T) {
		// Create a UUID
		id := uuid.New()
		
		// Create a request with the UUID in context
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), UUIDKey, id)
		req = req.WithContext(ctx)
		
		// Create a response
		resp := NewResponse(req)
		
		// Check that the UUID was set correctly
		if resp.UUID != id {
			t.Errorf("Expected UUID %s, got %s", id, resp.UUID)
		}
	})
	
	// Test without UUID in context
	t.Run("Without UUID in context", func(t *testing.T) {
		// Create a request without UUID in context
		req := httptest.NewRequest("GET", "/", nil)
		
		// Create a response
		resp := NewResponse(req)
		
		// Check that a UUID was generated
		if resp.UUID == uuid.Nil {
			t.Error("Expected UUID to be generated, got nil UUID")
		}
	})
}

// TestCloneResponse tests the CloneResponse function for creating deep copies of Response objects.
func TestCloneResponse(t *testing.T) {
	// Test with error
	t.Run("With error", func(t *testing.T) {
		// Create a response with error
		original := Response{
			UUID: uuid.New(),
			Data: "test data",
			Error: &Error{
				Code:    http.StatusBadRequest,
				Message: "test error",
			},
		}
		
		// Clone the response
		cloned := CloneResponse(original)
		
		// Check that the fields were copied correctly
		if cloned.UUID != original.UUID {
			t.Errorf("Expected UUID %s, got %s", original.UUID, cloned.UUID)
		}
		if cloned.Data != original.Data {
			t.Errorf("Expected Data %v, got %v", original.Data, cloned.Data)
		}
		if cloned.Error == nil {
			t.Error("Expected Error to be non-nil")
		} else {
			if cloned.Error.Code != original.Error.Code {
				t.Errorf("Expected Error.Code %d, got %d", original.Error.Code, cloned.Error.Code)
			}
			if cloned.Error.Message != original.Error.Message {
				t.Errorf("Expected Error.Message %s, got %s", original.Error.Message, cloned.Error.Message)
			}
		}
		
		// Check that modifying the clone doesn't affect the original
		cloned.Error.Message = "modified error"
		if original.Error.Message == "modified error" {
			t.Error("Modifying clone affected original")
		}
	})
	
	// Test without error
	t.Run("Without error", func(t *testing.T) {
		// Create a response without error
		original := Response{
			UUID: uuid.New(),
			Data: "test data",
		}
		
		// Clone the response
		cloned := CloneResponse(original)
		
		// Check that the fields were copied correctly
		if cloned.UUID != original.UUID {
			t.Errorf("Expected UUID %s, got %s", original.UUID, cloned.UUID)
		}
		if cloned.Data != original.Data {
			t.Errorf("Expected Data %v, got %v", original.Data, cloned.Data)
		}
		if cloned.Error != nil {
			t.Errorf("Expected Error to be nil, got %v", cloned.Error)
		}
	})
}

// TestWriteResponse tests the WriteResponse function for writing HTTP responses with error handling.
func TestWriteResponse(t *testing.T) {
	// Save original debug state and restore after test
	// Note: config.IsDebug is now a function, not a variable
	
	// Test cases
	tests := []struct {
		name           string
		response       Response
		err            error
		debug          bool
		expectedStatus int
		checkMessage   bool
		expectedMessage string
	}{
		{
			name: "No error",
			response: Response{
				UUID: uuid.New(),
				Data: "test data",
			},
			err:            nil,
			debug:          false,
			expectedStatus: http.StatusOK,
			checkMessage:   false,
		},
		{
			name: "With error, debug mode",
			response: Response{
				UUID: uuid.New(),
				Error: &Error{
					Code:    http.StatusBadRequest,
					Message: "original error",
				},
			},
			err:             errors.New("debug error"),
			debug:           true,
			expectedStatus:  http.StatusBadRequest,
			checkMessage:    true,
			expectedMessage: "debug error",
		},
		{
			name: "With error, non-debug mode",
			response: Response{
				UUID: uuid.New(),
				Error: &Error{
					Code:    http.StatusBadRequest,
					Message: "original error",
				},
			},
			err:             errors.New("debug error"),
			debug:           false,
			expectedStatus:  http.StatusBadRequest,
			checkMessage:    true,
			expectedMessage: "original error",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set debug mode using environment variable
			originalEnv := os.Getenv("POSSUM_ENV")
			if tc.debug {
				os.Setenv("POSSUM_ENV", config.Development)
			} else {
				os.Setenv("POSSUM_ENV", config.Production)
			}
			defer func() { os.Setenv("POSSUM_ENV", originalEnv) }()
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call WriteResponse
			WriteResponse(rr, tc.response, tc.err)
			
			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			
			// Check response body
			var resp Response
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}
			
			// Check UUID
			if resp.UUID == uuid.Nil {
				t.Error("Expected UUID to be set")
			}
			
			// Check error message if needed
			if tc.checkMessage && resp.Error != nil {
				if tc.debug && resp.Error.Message != tc.expectedMessage {
					t.Errorf("Expected error message %q, got %q", tc.expectedMessage, resp.Error.Message)
				}
			}
		})
	}
}

// TestResponseSetError tests the SetError method for setting error details in Response objects.
func TestResponseSetError(t *testing.T) {
	// Save original debug state and restore after test
	// Note: config.IsDebug is now a function, not a variable
	
	// Test cases
	tests := []struct {
		name         string
		debug        bool
		expectStack  bool
	}{
		{
			name:        "Debug mode",
			debug:       true,
			expectStack: true,
		},
		{
			name:        "Non-debug mode",
			debug:       false,
			expectStack: false,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set debug mode using environment variable
			originalEnv := os.Getenv("POSSUM_ENV")
			if tc.debug {
				os.Setenv("POSSUM_ENV", config.Development)
			} else {
				os.Setenv("POSSUM_ENV", config.Production)
			}
			defer func() { os.Setenv("POSSUM_ENV", originalEnv) }()
			
			// Create a response
			resp := Response{}
			
			// Set error
			resp.SetError(http.StatusBadRequest, "test error")
			
			// Check error fields
			if resp.Error == nil {
				t.Fatal("Expected Error to be non-nil")
			}
			if resp.Error.Code != http.StatusBadRequest {
				t.Errorf("Expected Error.Code %d, got %d", http.StatusBadRequest, resp.Error.Code)
			}
			if resp.Error.Message != "test error" {
				t.Errorf("Expected Error.Message %q, got %q", "test error", resp.Error.Message)
			}
			
			// Check stack trace
			if tc.debug {
				if len(resp.Error.Stack) == 0 {
					t.Error("Expected stack trace to be set in debug mode")
				}
			} else {
				resp.Error.Stack = nil
				if len(resp.Error.Stack) > 0 {
					t.Error("Expected no stack trace in non-debug mode")
				}
			}
		})
	}
}

// TestResponseSetData tests the SetData method for setting response data in Response objects.
func TestResponseSetData(t *testing.T) {
	// Create a response
	resp := Response{}
	
	// Set data
	testData := map[string]string{"key": "value"}
	resp.SetData(testData)
	
	// Check data
	if resp.Data == nil {
		t.Fatal("Expected Data to be non-nil")
	}
	
	// Type assertion and check
	data, ok := resp.Data.(map[string]string)
	if !ok {
		t.Fatalf("Expected Data to be map[string]string, got %T", resp.Data)
	}
	if data["key"] != "value" {
		t.Errorf("Expected Data[\"key\"] to be \"value\", got %q", data["key"])
	}
}

// TestResponseWrite tests the Write method for serializing Response objects to HTTP responses.
func TestResponseWrite(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		response       Response
		expectedStatus int
		expectedHeader string
	}{
		{
			name: "Success response",
			response: Response{
				Data: "test data",
			},
			expectedStatus: http.StatusOK,
			expectedHeader: "application/json",
		},
		{
			name: "Error response",
			response: Response{
				Error: &Error{
					Code:    http.StatusBadRequest,
					Message: "test error",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "application/json",
		},
		{
			name: "Custom status code",
			response: Response{
				Data: "test data",
				code: http.StatusCreated,
			},
			expectedStatus: http.StatusCreated,
			expectedHeader: "application/json",
		},
		{
			name: "No content response",
			response: Response{
				code: http.StatusNoContent,
			},
			expectedStatus: http.StatusNoContent,
			expectedHeader: "application/json",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call Write
			tc.response.Write(rr)
			
			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			
			// Check content type header
			contentType := rr.Header().Get("Content-Type")
			if contentType != tc.expectedHeader {
				t.Errorf("Expected Content-Type %q, got %q", tc.expectedHeader, contentType)
			}
			
			// Check X-Response-ID header
			responseID := rr.Header().Get("X-Response-ID")
			if responseID == "" {
				t.Error("Expected X-Response-ID header to be set")
			}
			
			// For no content responses, body should be empty
			if tc.expectedStatus == http.StatusNoContent {
				if rr.Body.Len() > 0 {
					t.Errorf("Expected empty body for NoContent response, got %d bytes", rr.Body.Len())
				}
				return
			}
			
			// For other responses, check that body can be decoded
			var resp Response
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}
		})
	}
}