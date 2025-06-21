package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware provides structured logging for HTTP requests
type LoggingMiddleware struct {
	logger *zap.Logger
}

// NewLoggingMiddleware creates a new logging middleware instance
func NewLoggingMiddleware(logger *zap.Logger) LoggingMiddleware {
	return LoggingMiddleware{
		logger: logger,
	}
}

// Handle returns a Gin middleware function for logging
func (m LoggingMiddleware) Handle() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log structured data using Zap
		m.logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Int("status_code", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("error", param.ErrorMessage),
			zap.Time("timestamp", param.TimeStamp),
		)

		// Return empty string as we're using Zap for logging
		return ""
	})
}
