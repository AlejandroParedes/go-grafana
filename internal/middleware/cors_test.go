package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestNewCORSMiddleware(t *testing.T) {
	logger := zap.NewNop()
	corsMiddleware := NewCORSMiddleware(logger)
	if corsMiddleware.logger == nil {
		t.Error("expected logger to be initialized")
	}
}

func TestCORSMiddleware_Handle(t *testing.T) {
	logger := zap.NewNop()
	corsMiddleware := NewCORSMiddleware(logger)
	handler := corsMiddleware.Handle()

	if handler == nil {
		t.Fatal("expected handler to be non-nil")
	}

	// Check that it sets CORS headers
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(handler)
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")

	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected Access-Control-Allow-Origin header to be set")
	}
}
