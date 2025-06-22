package models

import (
	"time"

	"go-grafana/internal/util"

	"gorm.io/gorm"
)

// APIKey represents an API key entity in the system
type APIKey struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	Name        string         `json:"name" gorm:"not null" validate:"required,min=2,max=100" example:"My API Key"`
	Key         string         `json:"key" gorm:"uniqueIndex;not null" example:"sk-1234567890abcdef"`
	Description string         `json:"description" gorm:"type:text" example:"API key for external service"`
	Active      bool           `json:"active" gorm:"default:true" example:"true"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for the APIKey model
func (APIKey) TableName() string {
	return "api_keys"
}

// CreateAPIKeyRequest represents the request payload for creating an API key
type CreateAPIKeyRequest struct {
	Name        string     `json:"name" binding:"required,min=2,max=100" example:"My API Key"`
	Description string     `json:"description" example:"API key for external service"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// UpdateAPIKeyRequest represents the request payload for updating an API key
type UpdateAPIKeyRequest struct {
	Name        string     `json:"name" binding:"required,min=2,max=100" example:"My API Key"`
	Description string     `json:"description" example:"API key for external service"`
	Active      bool       `json:"active" example:"true"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// APIKeyResponse represents the response payload for API key data
type APIKeyResponse struct {
	ID          uint       `json:"id" example:"1"`
	Name        string     `json:"name" example:"My API Key"`
	Key         string     `json:"key,omitempty" example:"sk-1234567890abcdef"`
	Description string     `json:"description" example:"API key for external service"`
	Active      bool       `json:"active" example:"true"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ToResponseWithKey converts an APIKey model to APIKeyResponse, including the plaintext key.
// This should only be used when creating a new key.
func (ak *APIKey) ToResponseWithKey(plainTextKey string) *APIKeyResponse {
	return &APIKeyResponse{
		ID:          ak.ID,
		Name:        ak.Name,
		Key:         plainTextKey,
		Description: ak.Description,
		Active:      ak.Active,
		ExpiresAt:   ak.ExpiresAt,
		CreatedAt:   ak.CreatedAt,
		UpdatedAt:   ak.UpdatedAt,
	}
}

// ToResponseWithoutKey converts an APIKey model to APIKeyResponse without exposing the key
func (ak *APIKey) ToResponseWithoutKey() *APIKeyResponse {
	return &APIKeyResponse{
		ID:          ak.ID,
		Name:        ak.Name,
		Key:         "***", // Mask the key for security
		Description: ak.Description,
		Active:      ak.Active,
		ExpiresAt:   ak.ExpiresAt,
		CreatedAt:   ak.CreatedAt,
		UpdatedAt:   ak.UpdatedAt,
	}
}

// FromCreateRequest populates an APIKey from CreateAPIKeyRequest and generates a new key
func (ak *APIKey) FromCreateRequest(req *CreateAPIKeyRequest) (string, error) {
	ak.Name = req.Name
	ak.Description = req.Description
	ak.ExpiresAt = req.ExpiresAt
	ak.Active = true // Default to active when creating

	plainTextKey, err := util.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	ak.Key = util.HashAPIKey(plainTextKey)

	return plainTextKey, nil
}

// FromUpdateRequest populates an APIKey from UpdateAPIKeyRequest
func (ak *APIKey) FromUpdateRequest(req *UpdateAPIKeyRequest) {
	ak.Name = req.Name
	ak.Description = req.Description
	ak.Active = req.Active
	ak.ExpiresAt = req.ExpiresAt
}

// IsExpired returns true if the API key has expired
func (ak *APIKey) IsExpired() bool {
	if ak.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ak.ExpiresAt)
}

// IsValid returns true if the API key is active and not expired
func (ak *APIKey) IsValid() bool {
	return ak.Active && !ak.IsExpired()
}
