package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

const (
	apiKeyLength = 32
	apiKeyPrefix = "tr_"
)

// GenerateAPIKey creates a new random API key.
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, apiKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return apiKeyPrefix + hex.EncodeToString(bytes), nil
}

// HashAPIKey creates a SHA-256 hash of an API key for storage.
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// ValidateAPIKey checks if an API key matches its hash.
func ValidateAPIKey(key, hash string) bool {
	keyHash := HashAPIKey(key)
	return ConstantTimeCompare(keyHash, hash)
}
