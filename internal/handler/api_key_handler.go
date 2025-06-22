package handler

import (
	"net/http"
	"strconv"

	"go-grafana/internal/domain/models"
	"go-grafana/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APIKeyHandler handles HTTP requests for API key operations
type APIKeyHandler struct {
	apiKeyService service.APIKeyService
	logger        *zap.Logger
}

// NewAPIKeyHandler creates a new instance of APIKeyHandler
func NewAPIKeyHandler(apiKeyService service.APIKeyService, logger *zap.Logger) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
		logger:        logger,
	}
}

// CreateAPIKey godoc
// @Summary Create a new API key
// @Description Create a new API key with the provided information
// @Tags api-keys
// @Accept json
// @Produce json
// @Param api_key body models.CreateAPIKeyRequest true "API key information"
// @Success 201 {object} models.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api-keys [post]
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	var req models.CreateAPIKeyRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind create API key request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Create API key
	apiKey, err := h.apiKeyService.CreateAPIKey(&req)
	if err != nil {
		h.logger.Error("Failed to create API key", zap.Error(err), zap.String("name", req.Name))

		status := http.StatusInternalServerError
		if err.Error() == "API key already exists" {
			status = http.StatusConflict
		} else if err.Error() == "name is required" {
			status = http.StatusBadRequest
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to create API key",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("API key created successfully", zap.Uint("api_key_id", apiKey.ID), zap.String("name", apiKey.Name))
	c.JSON(http.StatusCreated, apiKey)
}

// GetAPIKeys godoc
// @Summary Get all API keys
// @Description Retrieve a list of all API keys (keys are masked for security)
// @Tags api-keys
// @Produce json
// @Param X-API-Key header string true "API Key" default(sk-1234567890abcdef)
// @Success 200 {array} models.APIKeyResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api-keys [get]
func (h *APIKeyHandler) GetAPIKeys(c *gin.Context) {
	apiKeys, err := h.apiKeyService.GetAllAPIKeys()
	if err != nil {
		h.logger.Error("Failed to get API keys", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve API keys",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("API keys retrieved successfully", zap.Int("count", len(apiKeys)))
	c.JSON(http.StatusOK, apiKeys)
}

// GetAPIKeyByID godoc
// @Summary Get API key by ID
// @Description Retrieve a specific API key by its ID (key is masked for security)
// @Tags api-keys
// @Produce json
// @Param id path int true "API Key ID"
// @Param X-API-Key header string true "API Key" default(sk-1234567890abcdef)
// @Success 200 {object} models.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api-keys/{id} [get]
func (h *APIKeyHandler) GetAPIKeyByID(c *gin.Context) {
	// Parse API key ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid API key ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid API key ID",
			Message: "API key ID must be a valid integer",
		})
		return
	}

	apiKey, err := h.apiKeyService.GetAPIKeyByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get API key by ID", zap.Uint64("id", id), zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "API key not found" || err.Error() == "invalid API key ID" {
			status = http.StatusNotFound
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to retrieve API key",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("API key retrieved successfully", zap.Uint("api_key_id", apiKey.ID))
	c.JSON(http.StatusOK, apiKey)
}

// UpdateAPIKey godoc
// @Summary Update API key
// @Description Update an existing API key's information
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Param api_key body models.UpdateAPIKeyRequest true "Updated API key information"
// @Success 200 {object} models.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api-keys/{id} [put]
func (h *APIKeyHandler) UpdateAPIKey(c *gin.Context) {
	// Parse API key ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid API key ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid API key ID",
			Message: "API key ID must be a valid integer",
		})
		return
	}

	var req models.UpdateAPIKeyRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind update API key request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Update API key
	apiKey, err := h.apiKeyService.UpdateAPIKey(uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update API key", zap.Uint64("id", id), zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "API key not found" || err.Error() == "invalid API key ID" {
			status = http.StatusNotFound
		} else if err.Error() == "name is required" {
			status = http.StatusBadRequest
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to update API key",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("API key updated successfully", zap.Uint("api_key_id", apiKey.ID))
	c.JSON(http.StatusOK, apiKey)
}

// DeleteAPIKey godoc
// @Summary Delete API key
// @Description Delete an existing API key
// @Tags api-keys
// @Produce json
// @Param id path int true "API Key ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api-keys/{id} [delete]
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	// Parse API key ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid API key ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid API key ID",
			Message: "API key ID must be a valid integer",
		})
		return
	}

	// Delete API key
	err = h.apiKeyService.DeleteAPIKey(uint(id))
	if err != nil {
		h.logger.Error("Failed to delete API key", zap.Uint64("id", id), zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "API key not found" || err.Error() == "invalid API key ID" {
			status = http.StatusNotFound
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to delete API key",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("API key deleted successfully", zap.Uint64("id", id))
	c.Status(http.StatusNoContent)
}
