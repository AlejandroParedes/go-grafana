package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"go.uber.org/zap"
)

// This is a re-implementation of the middleware creation with a specific registry
// to avoid global state issues during testing.
func newTestMetricsMiddleware(reg *prometheus.Registry) MetricsMiddleware {
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total"},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds"},
		[]string{"method", "endpoint"},
	)
	httpRequestsInFlight := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "http_requests_in_flight"},
		[]string{"method", "endpoint"},
	)

	reg.MustRegister(httpRequestsTotal, httpRequestDuration, httpRequestsInFlight)

	return MetricsMiddleware{
		logger:               zap.NewNop(),
		httpRequestsTotal:    httpRequestsTotal,
		httpRequestDuration:  httpRequestDuration,
		httpRequestsInFlight: httpRequestsInFlight,
	}
}

func TestMetricsMiddleware_Handle(t *testing.T) {
	reg := prometheus.NewRegistry()
	metricsMiddleware := newTestMetricsMiddleware(reg)
	handler := metricsMiddleware.Handle()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(handler)
	router.GET("/test-metrics", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-metrics", nil)
	router.ServeHTTP(w, req)

	// Check counter
	err := testutil.CollectAndCompare(reg, strings.NewReader(`
		# HELP http_requests_total 
		# TYPE http_requests_total counter
		http_requests_total{endpoint="/test-metrics",method="GET",status="200"} 1
	`), "http_requests_total")
	if err != nil {
		t.Errorf("metric http_requests_total did not match expected value: %v", err)
	}

	// Check histogram
	// We just check that it has been observed once, not the value.
	metricFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("could not gather metrics: %v", err)
	}

	var histo *dto.MetricFamily
	for _, mf := range metricFamilies {
		if mf.GetName() == "http_request_duration_seconds" {
			histo = mf
			break
		}
	}

	if histo == nil || len(histo.GetMetric()) != 1 || histo.GetMetric()[0].GetHistogram().GetSampleCount() != 1 {
		t.Errorf("metric http_request_duration_seconds was not observed correctly")
	}
}

func TestMetricsMiddleware_MetricsHandler(t *testing.T) {
	logger := zap.NewNop()
	// The real NewMetricsMiddleware will use the default registry, which is fine for this test.
	metricsMiddleware := NewMetricsMiddleware(logger)
	handler := metricsMiddleware.MetricsHandler()

	if handler == nil {
		t.Fatal("expected handler to be non-nil")
	}

	// Make a request to see if it returns a 200
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/metrics", handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d for metrics handler, got %d", http.StatusOK, w.Code)
	}
}
