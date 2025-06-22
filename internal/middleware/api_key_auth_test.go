package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-grafana/internal/domain/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MockAPIKeyService is a mock of APIKeyService for middleware tests
type MockAPIKeyService struct {
	ValidateAPIKeyFunc func(key string) (*models.APIKey, error)
}

func (m *MockAPIKeyService) CreateAPIKey(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error) {
	return nil, nil
}
func (m *MockAPIKeyService) GetAPIKeyByID(id uint) (*models.APIKeyResponse, error) { return nil, nil }
func (m *MockAPIKeyService) GetAllAPIKeys() ([]*models.APIKeyResponse, error)      { return nil, nil }
func (m *MockAPIKeyService) UpdateAPIKey(id uint, req *models.UpdateAPIKeyRequest) (*models.APIKeyResponse, error) {
	return nil, nil
}
func (m *MockAPIKeyService) DeleteAPIKey(id uint) error { return nil }
func (m *MockAPIKeyService) ValidateAPIKey(key string) (*models.APIKey, error) {
	return m.ValidateAPIKeyFunc(key)
}

func TestAPIKeyAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAPIKeyService{}
	logger := zap.NewNop()
	middleware := APIKeyAuthMiddleware(mockService, logger)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("valid key", func(t *testing.T) {
		mockService.ValidateAPIKeyFunc = func(key string) (*models.APIKey, error) {
			if key == "valid-key" {
				return &models.APIKey{ID: 1, Name: "test-key"}, nil
			}
			return nil, errors.New("invalid key")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-API-Key", "valid-key")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		mockService.ValidateAPIKeyFunc = func(key string) (*models.APIKey, error) {
			return nil, errors.New("invalid key")
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-API-Key", "invalid-key")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetAPIKeyFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedKey := &models.APIKey{ID: 1}
	c.Set("api_key", expectedKey)

	key, exists := GetAPIKeyFromContext(c)
	if !exists {
		t.Fatal("expected api_key to exist in context")
	}
	if key.(*models.APIKey).ID != expectedKey.ID {
		t.Errorf("retrieved key incorrect")
	}
}

func TestGetAPIKeyIDFromContext(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("api_key_id", uint(1))

	id, exists := GetAPIKeyIDFromContext(c)
	if !exists || id != 1 {
		t.Errorf("expected id 1, got %d (exists: %t)", id, exists)
	}
}

func TestGetAPIKeyNameFromContext(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("api_key_name", "test-key")

	name, exists := GetAPIKeyNameFromContext(c)
	if !exists || name != "test-key" {
		t.Errorf("expected name 'test-key', got '%s' (exists: %t)", name, exists)
	}
}
