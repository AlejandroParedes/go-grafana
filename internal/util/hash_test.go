package util

import (
	"strings"
	"testing"
)

func TestGenerateAPIKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key1, err := GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey() error = %v, wantErr %v", err, false)
		}

		if !strings.HasPrefix(key1, "sk-") {
			t.Errorf("GenerateAPIKey() key = %v, want prefix %v", key1, "sk-")
		}

		if len(key1) != 67 { // sk- + 64 hex chars
			t.Errorf("GenerateAPIKey() key length = %v, want %v", len(key1), 67)
		}

		key2, err := GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey() error = %v, wantErr %v", err, false)
		}

		if key1 == key2 {
			t.Errorf("GenerateAPIKey() generated two identical keys: %v", key1)
		}
	})
}

func TestHashAPIKey(t *testing.T) {
	t.Run("hash consistency", func(t *testing.T) {
		key := "my-secret-key"
		hash1 := HashAPIKey(key)
		hash2 := HashAPIKey(key)

		if hash1 != hash2 {
			t.Errorf("HashAPIKey() produced different hashes for the same key: hash1 = %v, hash2 = %v", hash1, hash2)
		}
	})

	t.Run("hash correctness", func(t *testing.T) {
		key := "my-secret-key"
		expectedHash := "1311f8fc80a7ea28d78dd7723f09c44c1754cd35160ca8e7133ae3d7f636a19a"
		hash := HashAPIKey(key)

		if hash != expectedHash {
			t.Errorf("HashAPIKey() hash = %v, want %v", hash, expectedHash)
		}
	})

	t.Run("different keys have different hashes", func(t *testing.T) {
		key1 := "my-secret-key-1"
		key2 := "my-secret-key-2"
		hash1 := HashAPIKey(key1)
		hash2 := HashAPIKey(key2)

		if hash1 == hash2 {
			t.Errorf("HashAPIKey() produced the same hash for different keys")
		}
	})
}
