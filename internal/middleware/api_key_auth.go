package middleware

import (
	"net/http"
	"strings"

	"go-grafana/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APIKeyAuthMiddleware creates middleware for API key authentication
func APIKeyAuthMiddleware(apiKeyService service.APIKeyService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKeyHeader := c.GetHeader("X-API-Key")
		if apiKeyHeader == "" {
			logger.Warn("Missing API key header", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "API key is required",
			})
			c.Abort()
			return
		}

		// Clean the API key (remove any whitespace or "Bearer " prefix)
		apiKey := strings.TrimSpace(apiKeyHeader)
		if strings.HasPrefix(apiKey, "Bearer ") {
			apiKey = strings.TrimSpace(strings.TrimPrefix(apiKey, "Bearer "))
		}

		if apiKey == "" {
			logger.Warn("Empty API key provided", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "API key cannot be empty",
			})
			c.Abort()
			return
		}

		// Validate the API key
		validatedAPIKey, err := apiKeyService.ValidateAPIKey(apiKey)
		if err != nil {
			logger.Warn("Invalid API key provided",
				zap.String("path", c.Request.URL.Path),
				zap.String("error", err.Error()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		// Store the validated API key in the context for potential use in handlers
		c.Set("api_key", validatedAPIKey)
		c.Set("api_key_id", validatedAPIKey.ID)
		c.Set("api_key_name", validatedAPIKey.Name)

		logger.Debug("API key validated successfully",
			zap.Uint("api_key_id", validatedAPIKey.ID),
			zap.String("api_key_name", validatedAPIKey.Name),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// GetAPIKeyFromContext retrieves the API key from the Gin context
func GetAPIKeyFromContext(c *gin.Context) (interface{}, bool) {
	return c.Get("api_key")
}

// GetAPIKeyIDFromContext retrieves the API key ID from the Gin context
func GetAPIKeyIDFromContext(c *gin.Context) (uint, bool) {
	apiKeyID, exists := c.Get("api_key_id")
	if !exists {
		return 0, false
	}

	if id, ok := apiKeyID.(uint); ok {
		return id, true
	}

	return 0, false
}

// GetAPIKeyNameFromContext retrieves the API key name from the Gin context
func GetAPIKeyNameFromContext(c *gin.Context) (string, bool) {
	apiKeyName, exists := c.Get("api_key_name")
	if !exists {
		return "", false
	}

	if name, ok := apiKeyName.(string); ok {
		return name, true
	}

	return "", false
}
