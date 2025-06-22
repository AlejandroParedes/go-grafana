package handler

import (
	"net/http"
	"strconv"

	"go-grafana/internal/domain/models"
	"go-grafana/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key" default(sk-1234567890abcdef)
// @Param user body models.CreateUserRequest true "User information"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind create user request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Create user
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err), zap.String("email", req.Email))

		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			status = http.StatusConflict
		} else if err.Error() == "email is required" || err.Error() == "first name is required" ||
			err.Error() == "last name is required" || err.Error() == "age must be between 1 and 120" {
			status = http.StatusBadRequest
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to create user",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("User created successfully", zap.Uint("user_id", user.ID), zap.String("email", user.Email))
	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.logger.Error("Failed to get users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve users",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("Users retrieved successfully", zap.Int("count", len(users)))
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Parse user ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid integer",
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get user by ID", zap.Uint64("id", id), zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "invalid user ID" {
			status = http.StatusNotFound
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to retrieve user",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("User retrieved successfully", zap.Uint("user_id", user.ID))
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key" default(sk-1234567890abcdef)
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "Updated user information"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse user ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid integer",
		})
		return
	}

	var req models.UpdateUserRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind update user request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Uint64("id", id), zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "invalid user ID" {
			status = http.StatusNotFound
		} else if err.Error() == "user with this email already exists" {
			status = http.StatusConflict
		} else if err.Error() == "email is required" || err.Error() == "first name is required" ||
			err.Error() == "last name is required" || err.Error() == "age must be between 1 and 120" {
			status = http.StatusBadRequest
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to update user",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("User updated successfully", zap.Uint("user_id", user.ID))
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param X-API-Key header string true "API Key" default(sk-1234567890abcdef)
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Parse user ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid integer",
		})
		return
	}

	// Delete user
	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		h.logger.Error("Failed to delete user", zap.Uint64("id", id), zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "invalid user ID" {
			status = http.StatusNotFound
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to delete user",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("User deleted successfully", zap.Uint64("id", id))
	c.Status(http.StatusNoContent)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid request body"`
}
