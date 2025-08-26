package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the claims included in JWT tokens for user authentication.
type JWTClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token with user ID and expiration time claims.
// Returns the claims, token string, and any error that occurred during generation.
func GenerateJWT(secretKey []byte, userID uuid.UUID, expiresAt *time.Time) (*JWTClaims, string, error) {
	// Set default expiration time if not provided
	expTime := time.Now().Add(24 * time.Hour) // Default: 24 hours
	if expiresAt != nil {
		expTime = *expiresAt
	}

	// Create claims
	claims := &JWTClaims{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: expTime,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to sign token: %w", err)
	}

	return claims, tokenString, nil
}

func ParseToken(secret []byte, tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v\n", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenUnverifiable
}
