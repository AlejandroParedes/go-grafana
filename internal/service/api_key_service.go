package service

import (
	"errors"

	"go-grafana/internal/domain/models"
	"go-grafana/internal/domain/repository"
	"go-grafana/internal/util"
)

// APIKeyService defines the interface for API key business operations
type APIKeyService interface {
	CreateAPIKey(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error)
	GetAPIKeyByID(id uint) (*models.APIKeyResponse, error)
	GetAllAPIKeys() ([]*models.APIKeyResponse, error)
	UpdateAPIKey(id uint, req *models.UpdateAPIKeyRequest) (*models.APIKeyResponse, error)
	DeleteAPIKey(id uint) error
	ValidateAPIKey(key string) (*models.APIKey, error)
}

// apiKeyService implements APIKeyService
type apiKeyService struct {
	apiKeyRepo repository.APIKeyRepository
}

// NewAPIKeyService creates a new instance of APIKeyService
func NewAPIKeyService(apiKeyRepo repository.APIKeyRepository) APIKeyService {
	return &apiKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

// CreateAPIKey creates a new API key
func (s *apiKeyService) CreateAPIKey(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	apiKey := &models.APIKey{}
	plainTextKey, err := apiKey.FromCreateRequest(req)
	if err != nil {
		return nil, err
	}

	err = s.apiKeyRepo.Create(apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey.ToResponseWithKey(plainTextKey), nil
}

// GetAPIKeyByID retrieves an API key by its ID
func (s *apiKeyService) GetAPIKeyByID(id uint) (*models.APIKeyResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid API key ID")
	}

	apiKey, err := s.apiKeyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return apiKey.ToResponseWithoutKey(), nil
}

// GetAllAPIKeys retrieves all API keys
func (s *apiKeyService) GetAllAPIKeys() ([]*models.APIKeyResponse, error) {
	apiKeys, err := s.apiKeyRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*models.APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		responses[i] = apiKey.ToResponseWithoutKey()
	}

	return responses, nil
}

// UpdateAPIKey updates an existing API key
func (s *apiKeyService) UpdateAPIKey(id uint, req *models.UpdateAPIKeyRequest) (*models.APIKeyResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid API key ID")
	}

	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	// Get existing API key
	existing, err := s.apiKeyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update with new data
	existing.FromUpdateRequest(req)

	err = s.apiKeyRepo.Update(existing)
	if err != nil {
		return nil, err
	}

	return existing.ToResponseWithoutKey(), nil
}

// DeleteAPIKey deletes an API key
func (s *apiKeyService) DeleteAPIKey(id uint) error {
	if id == 0 {
		return errors.New("invalid API key ID")
	}

	return s.apiKeyRepo.Delete(id)
}

// ValidateAPIKey validates an API key and returns the API key object if valid
func (s *apiKeyService) ValidateAPIKey(key string) (*models.APIKey, error) {
	if key == "" {
		return nil, errors.New("API key is required")
	}

	hashedKey := util.HashAPIKey(key)

	apiKey, err := s.apiKeyRepo.GetByKey(hashedKey)
	if err != nil {
		return nil, errors.New("invalid API key")
	}

	if !apiKey.IsValid() {
		return nil, errors.New("API key is not valid")
	}

	return apiKey, nil
}
