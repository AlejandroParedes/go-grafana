package service

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"go-grafana/internal/domain/models"
	"go-grafana/internal/util"
)

// MockAPIKeyRepository is a mock implementation of APIKeyRepository for testing
type MockAPIKeyRepository struct {
	CreateFunc      func(apiKey *models.APIKey) error
	GetByIDFunc     func(id uint) (*models.APIKey, error)
	GetByKeyFunc    func(key string) (*models.APIKey, error)
	GetAllFunc      func() ([]*models.APIKey, error)
	UpdateFunc      func(apiKey *models.APIKey) error
	DeleteFunc      func(id uint) error
	ExistsByKeyFunc func(key string) bool
}

func (m *MockAPIKeyRepository) Create(apiKey *models.APIKey) error {
	return m.CreateFunc(apiKey)
}
func (m *MockAPIKeyRepository) GetByID(id uint) (*models.APIKey, error) {
	return m.GetByIDFunc(id)
}
func (m *MockAPIKeyRepository) GetByKey(key string) (*models.APIKey, error) {
	return m.GetByKeyFunc(key)
}
func (m *MockAPIKeyRepository) GetAll() ([]*models.APIKey, error) {
	return m.GetAllFunc()
}
func (m *MockAPIKeyRepository) Update(apiKey *models.APIKey) error {
	return m.UpdateFunc(apiKey)
}
func (m *MockAPIKeyRepository) Delete(id uint) error {
	return m.DeleteFunc(id)
}
func (m *MockAPIKeyRepository) ExistsByKey(key string) bool {
	return m.ExistsByKeyFunc(key)
}

func TestNewAPIKeyService(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)
	if service == nil {
		t.Error("NewAPIKeyService() returned nil")
	}
}

func TestAPIKeyService_CreateAPIKey(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)

	t.Run("success", func(t *testing.T) {
		req := &models.CreateAPIKeyRequest{Name: "test key"}
		mockRepo.CreateFunc = func(apiKey *models.APIKey) error {
			apiKey.ID = 1
			return nil
		}

		resp, err := service.CreateAPIKey(req)
		if err != nil {
			t.Fatalf("CreateAPIKey() error = %v, wantErr %v", err, false)
		}
		if resp.Name != req.Name {
			t.Errorf("expected name %s, got %s", req.Name, resp.Name)
		}
		if resp.Key == "" {
			t.Error("expected a non-empty key")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		req := &models.CreateAPIKeyRequest{Name: ""}
		_, err := service.CreateAPIKey(req)
		if err == nil {
			t.Error("expected an error for empty name, got nil")
		}
	})

	t.Run("repository create error", func(t *testing.T) {
		req := &models.CreateAPIKeyRequest{Name: "test key"}
		mockRepo.CreateFunc = func(apiKey *models.APIKey) error {
			return errors.New("db error")
		}
		_, err := service.CreateAPIKey(req)
		if err == nil {
			t.Error("expected a repository error, got nil")
		}
	})
}

func TestAPIKeyService_GetAPIKeyByID(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedAPIKey := &models.APIKey{ID: 1, Name: "test"}
		mockRepo.GetByIDFunc = func(id uint) (*models.APIKey, error) {
			if id == 1 {
				return expectedAPIKey, nil
			}
			return nil, errors.New("not found")
		}
		resp, err := service.GetAPIKeyByID(1)
		if err != nil {
			t.Fatalf("GetAPIKeyByID() error = %v", err)
		}
		if resp.ID != expectedAPIKey.ID {
			t.Errorf("got %d, want %d", resp.ID, expectedAPIKey.ID)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := service.GetAPIKeyByID(0)
		if err == nil {
			t.Error("expected error for invalid id, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(id uint) (*models.APIKey, error) {
			return nil, errors.New("not found")
		}
		_, err := service.GetAPIKeyByID(99)
		if err == nil {
			t.Error("expected error for not found, got nil")
		}
	})
}

func TestAPIKeyService_GetAllAPIKeys(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)

	t.Run("success", func(t *testing.T) {
		keys := []*models.APIKey{{ID: 1}, {ID: 2}}
		mockRepo.GetAllFunc = func() ([]*models.APIKey, error) {
			return keys, nil
		}
		resps, err := service.GetAllAPIKeys()
		if err != nil {
			t.Fatalf("GetAllAPIKeys() error = %v", err)
		}
		if len(resps) != 2 {
			t.Fatalf("expected 2 keys, got %d", len(resps))
		}
	})

	t.Run("db error", func(t *testing.T) {
		mockRepo.GetAllFunc = func() ([]*models.APIKey, error) {
			return nil, errors.New("db error")
		}
		_, err := service.GetAllAPIKeys()
		if err == nil {
			t.Error("expected db error, got nil")
		}
	})
}

func TestAPIKeyService_UpdateAPIKey(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)

	t.Run("success", func(t *testing.T) {
		existingKey := &models.APIKey{ID: 1, Name: "old name"}
		mockRepo.GetByIDFunc = func(id uint) (*models.APIKey, error) {
			return existingKey, nil
		}
		mockRepo.UpdateFunc = func(apiKey *models.APIKey) error {
			return nil
		}

		req := &models.UpdateAPIKeyRequest{Name: "new name"}
		resp, err := service.UpdateAPIKey(1, req)
		if err != nil {
			t.Fatalf("UpdateAPIKey() error = %v", err)
		}
		if resp.Name != "new name" {
			t.Errorf("expected name to be updated to 'new name', got '%s'", resp.Name)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := &models.UpdateAPIKeyRequest{Name: "new name"}
		_, err := service.UpdateAPIKey(0, req)
		if err == nil {
			t.Error("expected error for invalid id, got nil")
		}
	})
}

func TestAPIKeyService_DeleteAPIKey(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.DeleteFunc = func(id uint) error {
			return nil
		}
		err := service.DeleteAPIKey(1)
		if err != nil {
			t.Fatalf("DeleteAPIKey() error = %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		err := service.DeleteAPIKey(0)
		if err == nil {
			t.Error("expected error for invalid id, got nil")
		}
	})
}

func TestAPIKeyService_ValidateAPIKey(t *testing.T) {
	mockRepo := &MockAPIKeyRepository{}
	service := NewAPIKeyService(mockRepo)

	plainTextKey := "valid-key"
	hashedKey := util.HashAPIKey(plainTextKey)

	t.Run("valid key", func(t *testing.T) {
		validKey := &models.APIKey{ID: 1, Key: hashedKey, Active: true}
		mockRepo.GetByKeyFunc = func(key string) (*models.APIKey, error) {
			if key == hashedKey {
				return validKey, nil
			}
			return nil, errors.New("not found")
		}

		apiKey, err := service.ValidateAPIKey(plainTextKey)
		if err != nil {
			t.Fatalf("ValidateAPIKey() error = %v", err)
		}
		if !reflect.DeepEqual(apiKey, validKey) {
			t.Errorf("got %+v, want %+v", apiKey, validKey)
		}
	})

	t.Run("inactive key", func(t *testing.T) {
		inactiveKey := &models.APIKey{ID: 1, Key: hashedKey, Active: false}
		mockRepo.GetByKeyFunc = func(key string) (*models.APIKey, error) {
			return inactiveKey, nil
		}
		_, err := service.ValidateAPIKey(plainTextKey)
		if err == nil {
			t.Error("expected error for inactive key, got nil")
		}
	})

	t.Run("expired key", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)
		expiredKey := &models.APIKey{ID: 1, Key: hashedKey, Active: true, ExpiresAt: &pastTime}
		mockRepo.GetByKeyFunc = func(key string) (*models.APIKey, error) {
			return expiredKey, nil
		}
		_, err := service.ValidateAPIKey(plainTextKey)
		if err == nil {
			t.Error("expected error for expired key, got nil")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		_, err := service.ValidateAPIKey("")
		if err == nil {
			t.Error("expected error for empty key, got nil")
		}
	})
}
