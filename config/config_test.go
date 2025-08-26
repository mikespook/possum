package config

import (
	"os"
	"testing"
)

// TestIsDebug tests the IsDebug function for determining debug mode based on environment.
func TestIsDebug(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Getenv("POSSUM_ENV")
	defer os.Setenv("POSSUM_ENV", originalEnv)
	
	// Test cases
	tests := []struct {
		name           string
		env            string
		expectedResult bool
	}{
		{
			name:           "Production environment",
			env:            Production,
			expectedResult: false,
		},
		{
			name:           "Development environment",
			env:            Development,
			expectedResult: true,
		},
		{
			name:           "Test environment",
			env:            Test,
			expectedResult: true,
		},
		{
			name:           "Empty environment",
			env:            "",
			expectedResult: true,
		},
		{
			name:           "Custom environment",
			env:            "custom",
			expectedResult: true,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("POSSUM_ENV", tc.env)
			
			// Reset defaultEnv to simulate init()
			defaultEnv = tc.env
			
			// Test IsDebug
			result := IsDebug()
			
			// Check result
			if result != tc.expectedResult {
				t.Errorf("Expected IsDebug() to be %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

// TestIsDev tests the IsDev function for determining development environment mode.
func TestIsDev(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Getenv("POSSUM_ENV")
	defer os.Setenv("POSSUM_ENV", originalEnv)
	
	// Test cases
	tests := []struct {
		name           string
		env            string
		expectedResult bool
	}{
		{
			name:           "Production environment",
			env:            Production,
			expectedResult: false,
		},
		{
			name:           "Development environment",
			env:            Development,
			expectedResult: true,
		},
		{
			name:           "Test environment",
			env:            Test,
			expectedResult: false,
		},
		{
			name:           "Empty environment",
			env:            "",
			expectedResult: false,
		},
		{
			name:           "Custom environment",
			env:            "custom",
			expectedResult: false,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("POSSUM_ENV", tc.env)
			
			// Reset defaultEnv to simulate init()
			defaultEnv = tc.env
			
			// Test IsDev
			result := IsDev()
			
			// Check result
			if result != tc.expectedResult {
				t.Errorf("Expected IsDev() to be %v, got %v", tc.expectedResult, result)
			}
		})
	}
}