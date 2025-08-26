package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestGenerateJWT tests the GenerateJWT function for creating JWT tokens with claims.
func TestGenerateJWT(t *testing.T) {
	// Test secret
	secret := []byte("test-secret-key")
	
	// Test user ID
	userID := uuid.New()
	
	// Test cases
	tests := []struct {
		name        string
		userID      uuid.UUID
		expiresAt   *time.Time
		expectError bool
	}{
		{
			name:        "Default expiration",
			userID:      userID,
			expiresAt:   nil,
			expectError: false,
		},
		{
			name:        "Custom expiration",
			userID:      userID,
			expiresAt:   func() *time.Time { t := time.Now().Add(1 * time.Hour); return &t }(),
			expectError: false,
		},
		{
			name:        "Zero UUID",
			userID:      uuid.Nil,
			expiresAt:   nil,
			expectError: false,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Generate JWT
			claims, tokenString, err := GenerateJWT(secret, tc.userID, tc.expiresAt)
			
			// Check error
			if tc.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			
			// If no error, check claims and token
			if err == nil {
				// Check claims
				if claims.UserID != tc.userID {
					t.Errorf("Expected UserID %v, got %v", tc.userID, claims.UserID)
				}
				
				// Check token string
				if tokenString == "" {
					t.Error("Expected non-empty token string")
				}
				
				// Check expiration time
				if tc.expiresAt != nil {
					if !claims.ExpiresAt.Equal(*tc.expiresAt) {
						t.Errorf("Expected ExpiresAt %v, got %v", *tc.expiresAt, claims.ExpiresAt)
					}
				} else {
					// Default expiration should be ~24 hours from now
					expectedExp := time.Now().Add(24 * time.Hour)
					diff := claims.ExpiresAt.Sub(expectedExp)
					if diff < -5*time.Second || diff > 5*time.Second {
						t.Errorf("Expected ExpiresAt to be ~24 hours from now, got %v (diff: %v)", 
							claims.ExpiresAt, diff)
					}
				}
			}
		})
	}
}

// TestParseToken tests the ParseToken function for validating and parsing JWT tokens.
func TestParseToken(t *testing.T) {
	// Test secret
	secret := []byte("test-secret-key")
	
	// Test user ID
	userID := uuid.New()
	
	// Generate a valid token
	expTime := time.Now().Add(1 * time.Hour)
	_, validToken, err := GenerateJWT(secret, userID, &expTime)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}
	
	// Generate an expired token
	expiredTime := time.Now().Add(-1 * time.Hour)
	_, expiredToken, err := GenerateJWT(secret, userID, &expiredTime)
	if err != nil {
		t.Fatalf("Failed to generate expired test token: %v", err)
	}
	
	// Test cases
	tests := []struct {
		name        string
		secret      []byte
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			secret:      secret,
			token:       validToken,
			expectError: false,
		},
		{
			name:        "Expired token",
			secret:      secret,
			token:       expiredToken,
			expectError: true,
		},
		{
			name:        "Invalid token format",
			secret:      secret,
			token:       "invalid-token",
			expectError: true,
		},
		{
			name:        "Wrong secret",
			secret:      []byte("wrong-secret"),
			token:       validToken,
			expectError: true,
		},
		{
			name:        "Empty token",
			secret:      secret,
			token:       "",
			expectError: true,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Parse token
			claims, err := ParseToken(tc.secret, tc.token)
			
			// Check error
			if tc.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			
			// If no error, check claims
			if err == nil {
				if claims == nil {
					t.Error("Expected non-nil claims")
				} else {
					if claims.UserID != userID {
						t.Errorf("Expected UserID %v, got %v", userID, claims.UserID)
					}
				}
			}
		})
	}
}