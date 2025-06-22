package service

import (
	"errors"
	"testing"

	"go-grafana/internal/domain/models"
	"go-grafana/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	CreateFunc     func(user *models.User) error
	GetByIDFunc    func(id uint) (*models.User, error)
	GetAllFunc     func() ([]models.User, error)
	UpdateFunc     func(user *models.User) error
	DeleteFunc     func(id uint) error
	GetByEmailFunc func(email string) (*models.User, error)
	CountFunc      func() (int64, error)
}

func (m *MockUserRepository) Create(user *models.User) error        { return m.CreateFunc(user) }
func (m *MockUserRepository) GetByID(id uint) (*models.User, error) { return m.GetByIDFunc(id) }
func (m *MockUserRepository) GetAll() ([]models.User, error)        { return m.GetAllFunc() }
func (m *MockUserRepository) Update(user *models.User) error        { return m.UpdateFunc(user) }
func (m *MockUserRepository) Delete(id uint) error                  { return m.DeleteFunc(id) }
func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	return m.GetByEmailFunc(email)
}
func (m *MockUserRepository) Count() (int64, error) { return m.CountFunc() }

func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo, metrics.NewPrometheusMetrics(zap.NewNop(), prometheus.NewRegistry()))
	if service == nil {
		t.Error("NewUserService() returned nil")
	}
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo, metrics.NewPrometheusMetrics(zap.NewNop(), prometheus.NewRegistry()))

	t.Run("success", func(t *testing.T) {
		req := &models.CreateUserRequest{Email: "test@example.com", FirstName: "Test", LastName: "User", Age: 30}
		mockRepo.GetByEmailFunc = func(email string) (*models.User, error) {
			return nil, errors.New("not found")
		}
		mockRepo.CreateFunc = func(user *models.User) error {
			user.ID = 1
			return nil
		}
		mockRepo.CountFunc = func() (int64, error) {
			return 1, nil
		}

		resp, err := service.CreateUser(req)
		if err != nil {
			t.Fatalf("CreateUser() error = %v", err)
		}
		if resp.Email != req.Email {
			t.Errorf("expected email %s, got %s", req.Email, resp.Email)
		}
	})

	t.Run("email exists", func(t *testing.T) {
		req := &models.CreateUserRequest{Email: "test@example.com", FirstName: "Test", LastName: "User", Age: 30}
		mockRepo.GetByEmailFunc = func(email string) (*models.User, error) {
			return &models.User{ID: 1, Email: email}, nil
		}
		_, err := service.CreateUser(req)
		if err == nil {
			t.Error("expected an error for existing email, got nil")
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		req := &models.CreateUserRequest{Email: ""} // Invalid
		_, err := service.CreateUser(req)
		if err == nil {
			t.Error("expected an error for invalid request, got nil")
		}
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo, metrics.NewPrometheusMetrics(zap.NewNop(), prometheus.NewRegistry()))

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.User{ID: 1}
		mockRepo.GetByIDFunc = func(id uint) (*models.User, error) {
			if id == 1 {
				return expectedUser, nil
			}
			return nil, errors.New("not found")
		}
		user, err := service.GetUserByID(1)
		if err != nil {
			t.Fatalf("GetUserByID() error = %v", err)
		}
		if user.ID != 1 {
			t.Errorf("expected user ID 1, got %d", user.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(id uint) (*models.User, error) {
			return nil, errors.New("not found")
		}
		_, err := service.GetUserByID(99)
		if err == nil {
			t.Error("expected error for user not found, got nil")
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo, metrics.NewPrometheusMetrics(zap.NewNop(), prometheus.NewRegistry()))

	t.Run("success", func(t *testing.T) {
		req := &models.UpdateUserRequest{Email: "new@example.com", FirstName: "New", LastName: "Name", Age: 40}
		existingUser := &models.User{ID: 1, Email: "old@example.com"}
		mockRepo.GetByIDFunc = func(id uint) (*models.User, error) {
			return existingUser, nil
		}
		mockRepo.GetByEmailFunc = func(email string) (*models.User, error) {
			return nil, errors.New("not found") // No conflict
		}
		mockRepo.UpdateFunc = func(user *models.User) error {
			return nil
		}

		resp, err := service.UpdateUser(1, req)
		if err != nil {
			t.Fatalf("UpdateUser() error = %v", err)
		}
		if resp.Email != req.Email {
			t.Errorf("expected email %s, got %s", req.Email, resp.Email)
		}
	})

	t.Run("email conflict", func(t *testing.T) {
		req := &models.UpdateUserRequest{Email: "conflict@example.com", FirstName: "A", LastName: "B", Age: 10}
		mockRepo.GetByIDFunc = func(id uint) (*models.User, error) {
			return &models.User{ID: 1, Email: "original@example.com"}, nil
		}
		mockRepo.GetByEmailFunc = func(email string) (*models.User, error) {
			return &models.User{ID: 2, Email: "conflict@example.com"}, nil // Other user has this email
		}
		_, err := service.UpdateUser(1, req)
		if err == nil {
			t.Error("expected error for email conflict, got nil")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo, metrics.NewPrometheusMetrics(zap.NewNop(), prometheus.NewRegistry()))

	t.Run("success", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(id uint) (*models.User, error) {
			return &models.User{ID: 1}, nil
		}
		mockRepo.DeleteFunc = func(id uint) error {
			return nil
		}
		mockRepo.CountFunc = func() (int64, error) {
			return 0, nil
		}

		err := service.DeleteUser(1)
		if err != nil {
			t.Fatalf("DeleteUser() error = %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(id uint) (*models.User, error) {
			return nil, errors.New("not found")
		}
		err := service.DeleteUser(1)
		if err == nil {
			t.Error("expected error for user not found, got nil")
		}
	})
}
