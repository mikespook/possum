package log

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// TestNew tests the New function for creating a logger with configuration.
func TestNew(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		config         *Config
		expectedLevel  zerolog.Level
		checkFileError bool
	}{
		{
			name: "Default config",
			config: &Config{
				Level:    "",
				Filename: "",
			},
			expectedLevel:  zerolog.TraceLevel,
			checkFileError: false,
		},
		{
			name: "Custom level",
			config: &Config{
				Level:    "info",
				Filename: "",
			},
			expectedLevel:  zerolog.InfoLevel,
			checkFileError: false,
		},
		{
			name: "Invalid level",
			config: &Config{
				Level:    "invalid",
				Filename: "",
			},
			expectedLevel:  zerolog.TraceLevel,
			checkFileError: false,
		},
		{
			name: "Non-existent file path",
			config: &Config{
				Level:    "debug",
				Filename: "/non/existent/path/log.txt",
			},
			expectedLevel:  zerolog.DebugLevel,
			checkFileError: true,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create logger
			logger := New(tc.config)
			
			// Check level
			if logger.GetLevel() != tc.expectedLevel {
				t.Errorf("Expected level %v, got %v", tc.expectedLevel, logger.GetLevel())
			}
		})
	}
}

// TestInit tests the Init function for initializing the global logger.
func TestInit(t *testing.T) {
	// Save original logger and restore after test
	originalLogger := logger
	defer func() { logger = originalLogger }()
	
	// Test config
	config := &Config{
		Level:    "info",
		Filename: "",
	}
	
	// Initialize logger
	Init(config)
	
	// Check level
	if logger.GetLevel() != zerolog.InfoLevel {
		t.Errorf("Expected level %v, got %v", zerolog.InfoLevel, logger.GetLevel())
	}
}

// TestOutput tests the Output function for redirecting logger output.
func TestOutput(t *testing.T) {
	// Create a buffer to capture output
	buf := &bytes.Buffer{}
	
	// Create a logger with the buffer as output
	testLogger := Output(buf)
	
	// Write a message
	testLogger.Info().Msg("test message")
	
	// Check output
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got %q", output)
	}
	if !strings.Contains(output, "\"level\":\"info\"") {
		t.Errorf("Expected output to contain info level, got %q", output)
	}
}

// TestWith tests the With function for adding contextual fields to the logger.
func TestWith(t *testing.T) {
	// Create a buffer to capture output
	buf := &bytes.Buffer{}
	
	// Create a logger with the buffer as output and a field
	testLogger := With().Str("field", "value").Logger().Output(buf)
	
	// Write a message
	testLogger.Info().Msg("test message")
	
	// Check output
	output := buf.String()
	if !strings.Contains(output, "\"field\":\"value\"") {
		t.Errorf("Expected output to contain field, got %q", output)
	}
}

// TestLevel tests the Level function for setting the logger's minimum level.
func TestLevel(t *testing.T) {
	// Create a buffer to capture output
	buf := &bytes.Buffer{}
	
	// Create a logger with debug level
	testLogger := Level(zerolog.DebugLevel).Output(buf)
	
	// Write messages at different levels
	testLogger.Debug().Msg("debug message")
	testLogger.Info().Msg("info message")
	testLogger.Trace().Msg("trace message")
	
	// Check output
	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected output to contain debug message, got %q", output)
	}
	if !strings.Contains(output, "info message") {
		t.Errorf("Expected output to contain info message, got %q", output)
	}
	if strings.Contains(output, "trace message") {
		t.Errorf("Expected output to not contain trace message, got %q", output)
	}
}

// TestLogLevels tests the log level functions (Trace, Debug, Info, etc.) for proper level logging.
func TestLogLevels(t *testing.T) {
	// Test cases for different log levels
	tests := []struct {
		name     string
		logFunc  func() *zerolog.Event
		expected string
	}{
		{
			name:     "Trace",
			logFunc:  Trace,
			expected: "\"level\":\"trace\"",
		},
		{
			name:     "Debug",
			logFunc:  Debug,
			expected: "\"level\":\"debug\"",
		},
		{
			name:     "Info",
			logFunc:  Info,
			expected: "\"level\":\"info\"",
		},
		{
			name:     "Warn",
			logFunc:  Warn,
			expected: "\"level\":\"warn\"",
		},
		{
			name:     "Error",
			logFunc:  Error,
			expected: "\"level\":\"error\"",
		},
		{
			name:     "Log",
			logFunc:  Log,
			expected: "\"message\":",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to capture output
			buf := &bytes.Buffer{}
			
			// Set output to buffer
			originalLogger := logger
			logger.Logger = logger.Output(buf)
			defer func() { logger = originalLogger }()
			
			// Write message
			tc.logFunc().Msg("test message")
			
			// Check output
			output := buf.String()
			if !strings.Contains(output, tc.expected) {
				t.Errorf("Expected output to contain %q, got %q", tc.expected, output)
			}
			if !strings.Contains(output, "test message") {
				t.Errorf("Expected output to contain message, got %q", output)
			}
		})
	}
}

// TestPrint tests the Print function for debug-level logging with fmt.Print semantics.
func TestPrint(t *testing.T) {
	// Create a buffer to capture output
	buf := &bytes.Buffer{}
	
	// Set output to buffer
	originalLogger := logger
	logger.Logger = logger.Output(buf)
	defer func() { logger = originalLogger }()
	
	// Print message
	Print("test", "message")
	
	// Check output
	output := buf.String()
	if !strings.Contains(output, "testmessage") {
		t.Errorf("Expected output to contain 'testmessage', got %q", output)
	}
	if !strings.Contains(output, "\"level\":\"debug\"") {
		t.Errorf("Expected output to contain debug level, got %q", output)
	}
}

// TestPrintf tests the Printf function for debug-level logging with fmt.Printf semantics.
func TestPrintf(t *testing.T) {
	// Create a buffer to capture output
	buf := &bytes.Buffer{}
	
	// Set output to buffer
	originalLogger := logger
	logger.Logger = logger.Output(buf)
	defer func() { logger = originalLogger }()
	
	// Print message
	Printf("test %s", "message")
	
	// Check output
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got %q", output)
	}
	if !strings.Contains(output, "\"level\":\"debug\"") {
		t.Errorf("Expected output to contain debug level, got %q", output)
	}
}

// TestCtx tests the Ctx function for retrieving a logger from context.
func TestCtx(t *testing.T) {
	// Create a context with logger
	buf := &bytes.Buffer{}
	testLogger := zerolog.New(buf).With().Str("context_field", "context_value").Logger()
	ctx := testLogger.WithContext(context.Background())
	
	// Get logger from context
	loggerFromCtx := Ctx(ctx)
	
	// Write message
	loggerFromCtx.Info().Msg("context test")
	
	// Check output
	output := buf.String()
	if !strings.Contains(output, "context_field") {
		t.Errorf("Expected output to contain context field, got %q", output)
	}
	if !strings.Contains(output, "context test") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
}