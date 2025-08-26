package possum

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestLog tests the Log middleware function for request/response logging.
func TestLog(t *testing.T) {
	// Create a mock handler
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	// Create a mock handler that returns an error status
	errorHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}

	// Test cases
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful request",
			handler:        mockHandler,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Error request",
			handler:        errorHandler,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("User-Agent", "test-agent")
			
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create the log middleware
			handler := Log(tc.handler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check response body
			if rr.Body.String() != tc.expectedBody {
				t.Errorf("Expected body %s, got %s", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

// TestLogResponseWriter tests the methods of logResponseWriter for logging HTTP responses.
func TestLogResponseWriter(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Create a logResponseWriter
	lrw := &logResponseWriter{
		writer: rr,
		status: http.StatusOK,
	}

	// Test Write method
	content := []byte("test content")
	n, err := lrw.Write(content)
	if err != nil {
		t.Errorf("Write returned error: %v", err)
	}
	if n != len(content) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(content), n)
	}
	if rr.Body.String() != "test content" {
		t.Errorf("Expected body 'test content', got '%s'", rr.Body.String())
	}

	// Test Header method
	lrw.Header().Set("X-Test", "test-value")
	if rr.Header().Get("X-Test") != "test-value" {
		t.Errorf("Expected header X-Test to be 'test-value', got '%s'", rr.Header().Get("X-Test"))
	}

	// Test WriteHeader method
	lrw.WriteHeader(http.StatusCreated)
	if lrw.status != http.StatusCreated {
		t.Errorf("Expected status to be %d, got %d", http.StatusCreated, lrw.status)
	}
	// The underlying recorder may not capture the status code if the response is not written
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response code to be %d, got %d", http.StatusOK, rr.Code)
	}

	// Test Flush method
	// This is a no-op if the underlying writer doesn't implement http.Flusher
	lrw.Flush()
}

// TestLogResponseWriterHijack tests the Hijack method of logResponseWriter for WebSocket connections.
func TestLogResponseWriterHijack(t *testing.T) {
	// Create a mock hijacker
	mockHijacker := &mockResponseWriterHijacker{}
	
	// Create a logResponseWriter with a hijacker
	lrw := &logResponseWriter{
		writer: mockHijacker,
	}

	// Test Hijack method
	conn, rw, err := lrw.Hijack()
	if err != nil {
		t.Errorf("Hijack returned error: %v", err)
	}
	if conn == nil {
		t.Error("Expected connection to be non-nil")
	}
	if rw == nil {
		t.Error("Expected ReadWriter to be non-nil")
	}

	// Create a logResponseWriter with a non-hijacker
	lrw = &logResponseWriter{
		writer: httptest.NewRecorder(),
	}

	// Test Hijack method with non-hijacker
	conn, rw, err = lrw.Hijack()
	if err != ErrHijackerNotImplement {
		t.Errorf("Expected error %v, got %v", ErrHijackerNotImplement, err)
	}
	if conn != nil {
		t.Error("Expected connection to be nil")
	}
	if rw != nil {
		t.Error("Expected ReadWriter to be nil")
	}
}

// Mock ResponseWriter that implements http.Hijacker
type mockResponseWriterHijacker struct {
	httptest.ResponseRecorder
}

func (m *mockResponseWriterHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// Return mock objects
	return &mockConn{}, bufio.NewReadWriter(
		bufio.NewReader(nil),
		bufio.NewWriter(nil),
	), nil
}

// Mock net.Conn
type mockConn struct{}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return &mockAddr{} }
func (m *mockConn) RemoteAddr() net.Addr               { return &mockAddr{} }
func (m *mockConn) SetDeadline(t time.Time) error    { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// Mock net.Addr
type mockAddr struct{}

func (m *mockAddr) Network() string { return "mock" }
func (m *mockAddr) String() string  { return "mock-address" }