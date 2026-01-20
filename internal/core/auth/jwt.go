package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims for authentication.
type Claims struct {
	jwt.RegisteredClaims
	Type string `json:"type,omitempty"`
}

// TokenType defines the type of JWT token.
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// JWTManager handles JWT token operations.
type JWTManager struct {
	secret      []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
	issuer      string
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     "trelay",
	}
}

// GenerateAccessToken creates a new access token.
func (m *JWTManager) GenerateAccessToken() (string, error) {
	return m.generateToken(TokenTypeAccess, m.accessTTL)
}

// GenerateRefreshToken creates a new refresh token.
func (m *JWTManager) GenerateRefreshToken() (string, error) {
	return m.generateToken(TokenTypeRefresh, m.refreshTTL)
}

func (m *JWTManager) generateToken(tokenType TokenType, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
		Type: string(tokenType),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken validates a JWT token and returns its claims.
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// IsAccessToken checks if the claims are for an access token.
func (c *Claims) IsAccessToken() bool {
	return c.Type == string(TokenTypeAccess)
}

// IsRefreshToken checks if the claims are for a refresh token.
func (c *Claims) IsRefreshToken() bool {
	return c.Type == string(TokenTypeRefresh)
}
