package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager creates new JWT manager for the gateway
func NewJWTManager(secret string, duration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secret,
		tokenDuration: duration,
	}
}

// Generate token with userID (used by auth service, but kept here for consistency)
func (jm *JWTManager) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(jm.tokenDuration).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "backend-infra-gateway", // Issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

// Validate parses and validates a JWT token
func (jm *JWTManager) Validate(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jm.secretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, jwt.ErrTokenInvalidClaims
	}

	return token, claims, nil
}

// Secret returns the secret key for middleware usage
func (jm *JWTManager) Secret() string {
	return jm.secretKey
}

// GetTokenDuration returns the token duration
func (jm *JWTManager) GetTokenDuration() time.Duration {
	return jm.tokenDuration
}