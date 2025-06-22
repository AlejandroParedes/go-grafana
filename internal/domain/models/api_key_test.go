package models

import (
	"testing"
	"time"
)

func TestAPIKey_TableName(t *testing.T) {
	var apiKey APIKey
	if apiKey.TableName() != "api_keys" {
		t.Errorf("expected table name 'api_keys', got '%s'", apiKey.TableName())
	}
}

func TestAPIKey_ToResponseWithKey(t *testing.T) {
	now := time.Now()
	apiKey := &APIKey{
		ID:          1,
		Name:        "test key",
		Description: "test description",
		Active:      true,
		ExpiresAt:   &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	plainTextKey := "plain-text"

	resp := apiKey.ToResponseWithKey(plainTextKey)

	if resp.ID != apiKey.ID {
		t.Errorf("expected ID %d, got %d", apiKey.ID, resp.ID)
	}
	if resp.Name != apiKey.Name {
		t.Errorf("expected Name '%s', got '%s'", apiKey.Name, resp.Name)
	}
	if resp.Key != plainTextKey {
		t.Errorf("expected Key '%s', got '%s'", plainTextKey, resp.Key)
	}
}

func TestAPIKey_ToResponseWithoutKey(t *testing.T) {
	apiKey := &APIKey{ID: 1, Name: "test key"}
	resp := apiKey.ToResponseWithoutKey()
	if resp.Key != "***" {
		t.Errorf("expected masked key '***', got '%s'", resp.Key)
	}
}

func TestAPIKey_FromCreateRequest(t *testing.T) {
	req := &CreateAPIKeyRequest{
		Name:        "new key",
		Description: "new description",
	}
	apiKey := &APIKey{}
	plainTextKey, err := apiKey.FromCreateRequest(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apiKey.Name != req.Name {
		t.Errorf("expected Name '%s', got '%s'", req.Name, apiKey.Name)
	}
	if !apiKey.Active {
		t.Error("expected key to be active")
	}
	if plainTextKey == "" {
		t.Error("expected a plain text key to be generated")
	}
	if apiKey.Key == "" || apiKey.Key == plainTextKey {
		t.Error("expected key to be hashed")
	}
}

func TestAPIKey_FromUpdateRequest(t *testing.T) {
	now := time.Now()
	req := &UpdateAPIKeyRequest{
		Name:        "updated key",
		Description: "updated description",
		Active:      false,
		ExpiresAt:   &now,
	}
	apiKey := &APIKey{}
	apiKey.FromUpdateRequest(req)

	if apiKey.Name != req.Name {
		t.Errorf("expected Name '%s', got '%s'", req.Name, apiKey.Name)
	}
	if apiKey.Description != req.Description {
		t.Errorf("expected Description '%s', got '%s'", req.Description, apiKey.Description)
	}
	if apiKey.Active != req.Active {
		t.Errorf("expected Active %t, got %t", req.Active, apiKey.Active)
	}
	if apiKey.ExpiresAt != req.ExpiresAt {
		t.Errorf("expected ExpiresAt %v, got %v", req.ExpiresAt, apiKey.ExpiresAt)
	}
}

func TestAPIKey_IsExpired(t *testing.T) {
	t.Run("not expired if no expiry date", func(t *testing.T) {
		apiKey := &APIKey{ExpiresAt: nil}
		if apiKey.IsExpired() {
			t.Error("expected not expired")
		}
	})

	t.Run("not expired if expiry date is in the future", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour)
		apiKey := &APIKey{ExpiresAt: &future}
		if apiKey.IsExpired() {
			t.Error("expected not expired")
		}
	})

	t.Run("expired if expiry date is in the past", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		apiKey := &APIKey{ExpiresAt: &past}
		if !apiKey.IsExpired() {
			t.Error("expected expired")
		}
	})
}

func TestAPIKey_IsValid(t *testing.T) {
	t.Run("valid if active and not expired", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour)
		apiKey := &APIKey{Active: true, ExpiresAt: &future}
		if !apiKey.IsValid() {
			t.Error("expected valid")
		}
	})

	t.Run("invalid if not active", func(t *testing.T) {
		apiKey := &APIKey{Active: false}
		if apiKey.IsValid() {
			t.Error("expected invalid")
		}
	})

	t.Run("invalid if expired", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		apiKey := &APIKey{Active: true, ExpiresAt: &past}
		if apiKey.IsValid() {
			t.Error("expected invalid")
		}
	})

	t.Run("valid if active and no expiry date", func(t *testing.T) {
		apiKey := &APIKey{Active: true, ExpiresAt: nil}
		if !apiKey.IsValid() {
			t.Error("expected valid")
		}
	})
}
