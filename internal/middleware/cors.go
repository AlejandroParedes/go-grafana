package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CORSMiddleware provides CORS configuration
type CORSMiddleware struct {
	logger *zap.Logger
}

// NewCORSMiddleware creates a new CORS middleware instance
func NewCORSMiddleware(logger *zap.Logger) CORSMiddleware {
	return CORSMiddleware{
		logger: logger,
	}
}

// Handle returns a Gin middleware function for CORS
func (m CORSMiddleware) Handle() gin.HandlerFunc {
	config := cors.DefaultConfig()

	// Allow all origins for development
	// In production, you should specify allowed origins
	config.AllowAllOrigins = true

	// Allow specific methods
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	// Allow specific headers
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}

	// Allow credentials
	config.AllowCredentials = true

	// Expose headers
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
	}

	m.logger.Info("CORS middleware configured",
		zap.Bool("allow_all_origins", config.AllowAllOrigins),
		zap.Strings("allow_methods", config.AllowMethods),
	)

	return cors.New(config)
}
