package repository

import (
	"errors"

	"go-grafana/internal/domain/models"

	"gorm.io/gorm"
)

// APIKeyRepository defines the interface for API key data operations
type APIKeyRepository interface {
	Create(apiKey *models.APIKey) error
	GetByID(id uint) (*models.APIKey, error)
	GetByKey(key string) (*models.APIKey, error)
	GetAll() ([]*models.APIKey, error)
	Update(apiKey *models.APIKey) error
	Delete(id uint) error
	ExistsByKey(key string) bool
}

// apiKeyRepository implements APIKeyRepository
type apiKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository creates a new instance of APIKeyRepository
func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{
		db: db,
	}
}

// Create creates a new API key in the database
func (r *apiKeyRepository) Create(apiKey *models.APIKey) error {
	if apiKey.Name == "" {
		return errors.New("name is required")
	}

	if apiKey.Key == "" {
		return errors.New("key is required")
	}

	// Check if key already exists
	if r.ExistsByKey(apiKey.Key) {
		return errors.New("API key already exists")
	}

	result := r.db.Create(apiKey)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetByID retrieves an API key by its ID
func (r *apiKeyRepository) GetByID(id uint) (*models.APIKey, error) {
	if id == 0 {
		return nil, errors.New("invalid API key ID")
	}

	var apiKey models.APIKey
	result := r.db.First(&apiKey, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("API key not found")
		}
		return nil, result.Error
	}

	return &apiKey, nil
}

// GetByKey retrieves an API key by its key value
func (r *apiKeyRepository) GetByKey(key string) (*models.APIKey, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}

	var apiKey models.APIKey
	result := r.db.Where("key = ?", key).First(&apiKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("API key not found")
		}
		return nil, result.Error
	}

	return &apiKey, nil
}

// GetAll retrieves all API keys from the database
func (r *apiKeyRepository) GetAll() ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	result := r.db.Find(&apiKeys)
	if result.Error != nil {
		return nil, result.Error
	}

	return apiKeys, nil
}

// Update updates an existing API key in the database
func (r *apiKeyRepository) Update(apiKey *models.APIKey) error {
	if apiKey.ID == 0 {
		return errors.New("invalid API key ID")
	}

	if apiKey.Name == "" {
		return errors.New("name is required")
	}

	// Check if API key exists
	existing, err := r.GetByID(apiKey.ID)
	if err != nil {
		return err
	}

	// Update only allowed fields (don't update the key itself)
	existing.Name = apiKey.Name
	existing.Description = apiKey.Description
	existing.Active = apiKey.Active
	existing.ExpiresAt = apiKey.ExpiresAt

	result := r.db.Save(existing)
	if result.Error != nil {
		return result.Error
	}

	// Copy updated data back to the original object
	*apiKey = *existing

	return nil
}

// Delete removes an API key from the database
func (r *apiKeyRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid API key ID")
	}

	// Check if API key exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	result := r.db.Delete(&models.APIKey{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ExistsByKey checks if an API key exists by its key value
func (r *apiKeyRepository) ExistsByKey(key string) bool {
	if key == "" {
		return false
	}

	var count int64
	r.db.Model(&models.APIKey{}).Where("key = ?", key).Count(&count)
	return count > 0
}
