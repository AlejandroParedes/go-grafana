package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateAPIKey generates a new secure API key.
// The key is a 32-byte random string, hex-encoded, and prefixed with "sk-".
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk-" + hex.EncodeToString(bytes), nil
}

// HashAPIKey hashes an API key using SHA-256.
// This is used to store keys securely in the database.
func HashAPIKey(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
