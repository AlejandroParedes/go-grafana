package service

import (
	"errors"
	"fmt"

	"go-grafana/internal/domain/models"
	"go-grafana/internal/domain/repository"
	"go-grafana/pkg/metrics"
)

// UserService defines the interface for user business operations
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	UpdateUser(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(id uint) error
	GetUserCount() (int64, error)
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
	metrics  *metrics.PrometheusMetrics
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository, prometheusMetrics *metrics.PrometheusMetrics) UserService {
	return &userService{
		userRepo: userRepo,
		metrics:  prometheusMetrics,
	}
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user := &models.User{}
	user.FromCreateRequest(req)

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Record metrics
	s.metrics.RecordUserCreation()
	s.metrics.RecordUserAge(user.Age)

	// Update active users count
	if count, err := s.userRepo.Count(); err == nil {
		s.metrics.SetActiveUsers(count)
	}

	return user.ToResponse(), nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Convert to response format
	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}

	return responses, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if it conflicts with existing user
	if user.Email != req.Email {
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.New("user with this email already exists")
		}
	}

	// Update user data
	user.FromUpdateRequest(req)

	// Save to database
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Record metrics
	s.metrics.RecordUserUpdate()
	s.metrics.RecordUserAge(user.Age)

	return user.ToResponse(), nil
}

// DeleteUser removes a user from the system
func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete user
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Record metrics
	s.metrics.RecordUserDeletion()

	// Update active users count
	if count, err := s.userRepo.Count(); err == nil {
		s.metrics.SetActiveUsers(count)
	}

	return nil
}

// GetUserCount returns the total number of users
func (s *userService) GetUserCount() (int64, error) {
	count, err := s.userRepo.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}
	return count, nil
}

// validateCreateRequest validates the create user request
func (s *userService) validateCreateRequest(req *models.CreateUserRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.FirstName == "" {
		return errors.New("first name is required")
	}

	if req.LastName == "" {
		return errors.New("last name is required")
	}

	if req.Age <= 0 || req.Age > 120 {
		return errors.New("age must be between 1 and 120")
	}

	return nil
}

// validateUpdateRequest validates the update user request
func (s *userService) validateUpdateRequest(req *models.UpdateUserRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.FirstName == "" {
		return errors.New("first name is required")
	}

	if req.LastName == "" {
		return errors.New("last name is required")
	}

	if req.Age <= 0 || req.Age > 120 {
		return errors.New("age must be between 1 and 120")
	}

	return nil
}
