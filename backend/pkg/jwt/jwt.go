// Package jwt provides RS256 JWT signing, parsing, and claims extraction for WanderPlan.
package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the WanderPlan JWT payload.
type Claims struct {
	// UserID is the database UUID of the authenticated user.
	UserID string `json:"sub"`
	// Email is the user's email address.
	Email string `json:"email"`
	// Name is the user's display name.
	Name string `json:"name"`
	// AvatarURL is the user's profile picture URL.
	AvatarURL string `json:"avatar_url"`
	jwt.RegisteredClaims
}

// Manager handles RS256 JWT creation and verification.
type Manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	expiry     time.Duration
}

// NewManager creates a Manager from PEM-encoded RSA keys.
func NewManager(privatePEM, publicPEM []byte, expiry time.Duration) (*Manager, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, fmt.Errorf("parse rsa private key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, fmt.Errorf("parse rsa public key: %w", err)
	}
	return &Manager{privateKey: privateKey, publicKey: publicKey, expiry: expiry}, nil
}

// NewManagerFromBase64 creates a Manager from base64-encoded PEM strings (env-var friendly).
func NewManagerFromBase64(privateB64, publicB64 string, expiry time.Duration) (*Manager, error) {
	privPEM, err := base64.StdEncoding.DecodeString(privateB64)
	if err != nil {
		return nil, fmt.Errorf("decode private key base64: %w", err)
	}
	pubPEM, err := base64.StdEncoding.DecodeString(publicB64)
	if err != nil {
		return nil, fmt.Errorf("decode public key base64: %w", err)
	}
	return NewManager(privPEM, pubPEM, expiry)
}

// Sign creates a signed RS256 JWT for the given user attributes.
func (m *Manager) Sign(userID, email, name, avatarURL string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
			Issuer:    "wanderplan",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

// Parse validates a JWT string and returns its claims.
// Returns an error if the token is malformed, expired, or has an invalid signature.
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// PublicKey returns the RSA public key for external validators (e.g. api-gateway gRPC).
func (m *Manager) PublicKey() *rsa.PublicKey { return m.publicKey }
