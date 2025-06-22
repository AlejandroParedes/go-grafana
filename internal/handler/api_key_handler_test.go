package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-grafana/internal/domain/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MockAPIKeyService is a mock of APIKeyService
type MockAPIKeyService struct {
	CreateAPIKeyFunc   func(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error)
	GetAPIKeyByIDFunc  func(id uint) (*models.APIKeyResponse, error)
	GetAllAPIKeysFunc  func() ([]*models.APIKeyResponse, error)
	UpdateAPIKeyFunc   func(id uint, req *models.UpdateAPIKeyRequest) (*models.APIKeyResponse, error)
	DeleteAPIKeyFunc   func(id uint) error
	ValidateAPIKeyFunc func(key string) (*models.APIKey, error)
}

func (m *MockAPIKeyService) CreateAPIKey(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error) {
	return m.CreateAPIKeyFunc(req)
}
func (m *MockAPIKeyService) GetAPIKeyByID(id uint) (*models.APIKeyResponse, error) {
	return m.GetAPIKeyByIDFunc(id)
}
func (m *MockAPIKeyService) GetAllAPIKeys() ([]*models.APIKeyResponse, error) {
	return m.GetAllAPIKeysFunc()
}
func (m *MockAPIKeyService) UpdateAPIKey(id uint, req *models.UpdateAPIKeyRequest) (*models.APIKeyResponse, error) {
	return m.UpdateAPIKeyFunc(id, req)
}
func (m *MockAPIKeyService) DeleteAPIKey(id uint) error {
	return m.DeleteAPIKeyFunc(id)
}
func (m *MockAPIKeyService) ValidateAPIKey(key string) (*models.APIKey, error) {
	return m.ValidateAPIKeyFunc(key)
}

func setupTestRouter() (*gin.Engine, *MockAPIKeyService, *APIKeyHandler) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAPIKeyService{}
	logger := zap.NewNop()
	handler := NewAPIKeyHandler(mockService, logger)
	router := gin.Default()
	return router, mockService, handler
}

func TestAPIKeyHandler_CreateAPIKey(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.POST("/api-keys", handler.CreateAPIKey)

	t.Run("success", func(t *testing.T) {
		reqBody := models.CreateAPIKeyRequest{Name: "test-key"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.CreateAPIKeyFunc = func(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error) {
			return &models.APIKeyResponse{ID: 1, Name: req.Name}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api-keys", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
		var resp models.APIKeyResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Name != "test-key" {
			t.Errorf("expected name %s, got %s", "test-key", resp.Name)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api-keys", bytes.NewBuffer([]byte("{")))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		reqBody := models.CreateAPIKeyRequest{Name: "test-key"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.CreateAPIKeyFunc = func(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error) {
			return nil, errors.New("some service error")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api-keys", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestAPIKeyHandler_GetAPIKeys(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api-keys", handler.GetAPIKeys)

	t.Run("success", func(t *testing.T) {
		mockService.GetAllAPIKeysFunc = func() ([]*models.APIKeyResponse, error) {
			return []*models.APIKeyResponse{{ID: 1, Name: "key1"}, {ID: 2, Name: "key2"}}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api-keys", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resps []*models.APIKeyResponse
		json.Unmarshal(w.Body.Bytes(), &resps)
		if len(resps) != 2 {
			t.Errorf("expected %d keys, got %d", 2, len(resps))
		}
	})
}

func TestAPIKeyHandler_GetAPIKeyByID(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api-keys/:id", handler.GetAPIKeyByID)

	t.Run("success", func(t *testing.T) {
		mockService.GetAPIKeyByIDFunc = func(id uint) (*models.APIKeyResponse, error) {
			if id != 1 {
				t.Fatalf("mock service called with unexpected id: %d", id)
			}
			return &models.APIKeyResponse{ID: 1, Name: "test-key"}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api-keys/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockService.GetAPIKeyByIDFunc = func(id uint) (*models.APIKeyResponse, error) {
			return nil, errors.New("API key not found")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api-keys/99", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestAPIKeyHandler_DeleteAPIKey(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.DELETE("/api-keys/:id", handler.DeleteAPIKey)

	t.Run("success", func(t *testing.T) {
		mockService.DeleteAPIKeyFunc = func(id uint) error {
			return nil
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api-keys/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
