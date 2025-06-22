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

// MockUserService is a mock of UserService
type MockUserService struct {
	CreateUserFunc   func(req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByIDFunc  func(id uint) (*models.UserResponse, error)
	GetAllUsersFunc  func() ([]models.UserResponse, error)
	UpdateUserFunc   func(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUserFunc   func(id uint) error
	GetUserCountFunc func() (int64, error)
}

func (m *MockUserService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	return m.CreateUserFunc(req)
}
func (m *MockUserService) GetUserByID(id uint) (*models.UserResponse, error) {
	return m.GetUserByIDFunc(id)
}
func (m *MockUserService) GetAllUsers() ([]models.UserResponse, error) {
	return m.GetAllUsersFunc()
}
func (m *MockUserService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	return m.UpdateUserFunc(id, req)
}
func (m *MockUserService) DeleteUser(id uint) error {
	return m.DeleteUserFunc(id)
}
func (m *MockUserService) GetUserCount() (int64, error) {
	return m.GetUserCountFunc()
}

func setupUserTestRouter() (*gin.Engine, *MockUserService, *UserHandler) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	logger := zap.NewNop()
	// We pass nil for metrics here as the handler itself doesn't use it directly.
	// The service, which is mocked, is responsible for metrics.
	handler := NewUserHandler(mockService, logger)
	router := gin.Default()
	return router, mockService, handler
}

func TestUserHandler_CreateUser(t *testing.T) {
	router, mockService, handler := setupUserTestRouter()
	router.POST("/users", handler.CreateUser)

	t.Run("success", func(t *testing.T) {
		reqBody := models.CreateUserRequest{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Age:       30,
		}
		jsonBody, _ := json.Marshal(reqBody)
		mockService.CreateUserFunc = func(req *models.CreateUserRequest) (*models.UserResponse, error) {
			return &models.UserResponse{ID: 1, FirstName: req.FirstName}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

func TestUserHandler_GetUsers(t *testing.T) {
	router, mockService, handler := setupUserTestRouter()
	router.GET("/users", handler.GetUsers)

	t.Run("success", func(t *testing.T) {
		mockService.GetAllUsersFunc = func() ([]models.UserResponse, error) {
			return []models.UserResponse{{ID: 1}}, nil
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resps []models.UserResponse
		json.Unmarshal(w.Body.Bytes(), &resps)
		if len(resps) != 1 {
			t.Errorf("expected 1 user, got %d", len(resps))
		}
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	router, mockService, handler := setupUserTestRouter()
	router.GET("/users/:id", handler.GetUserByID)

	t.Run("success", func(t *testing.T) {
		mockService.GetUserByIDFunc = func(id uint) (*models.UserResponse, error) {
			return &models.UserResponse{ID: id}, nil
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockService.GetUserByIDFunc = func(id uint) (*models.UserResponse, error) {
			return nil, errors.New("user not found")
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/99", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	router, mockService, handler := setupUserTestRouter()
	router.PUT("/users/:id", handler.UpdateUser)

	t.Run("success", func(t *testing.T) {
		reqBody := models.UpdateUserRequest{
			Email:     "test@example.com",
			FirstName: "Updated",
			LastName:  "User",
			Age:       31,
			Active:    true,
		}
		jsonBody, _ := json.Marshal(reqBody)
		mockService.UpdateUserFunc = func(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
			return &models.UserResponse{ID: id, FirstName: req.FirstName}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	router, mockService, handler := setupUserTestRouter()
	router.DELETE("/users/:id", handler.DeleteUser)

	t.Run("success", func(t *testing.T) {
		mockService.DeleteUserFunc = func(id uint) error {
			return nil
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
