package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLoggingMiddleware(t *testing.T) {
	logger := zap.NewNop()
	loggingMiddleware := NewLoggingMiddleware(logger)
	if loggingMiddleware.logger == nil {
		t.Error("expected logger to be initialized")
	}
}

func TestLoggingMiddleware_Handle(t *testing.T) {
	var buffer bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buffer), zap.InfoLevel)
	logger := zap.New(core)

	loggingMiddleware := NewLoggingMiddleware(logger)
	handler := loggingMiddleware.Handle()

	if handler == nil {
		t.Fatal("expected handler to be non-nil")
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(handler)
	router.GET("/test-log", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-log", nil)

	router.ServeHTTP(w, req)

	logOutput := buffer.String()
	if !strings.Contains(logOutput, `"msg":"HTTP Request"`) {
		t.Errorf("log output should contain 'HTTP Request', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"path":"/test-log"`) {
		t.Errorf("log output should contain the path, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"status_code":200`) {
		t.Errorf("log output should contain the status code, got: %s", logOutput)
	}
}
