package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// MetricsMiddleware provides Prometheus metrics collection
type MetricsMiddleware struct {
	logger *zap.Logger
	// HTTP request metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight *prometheus.GaugeVec
}

// NewMetricsMiddleware creates a new metrics middleware instance
func NewMetricsMiddleware(logger *zap.Logger) MetricsMiddleware {
	// Define metrics
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestsInFlight := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
		[]string{"method", "endpoint"},
	)

	// Register metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestsInFlight)

	return MetricsMiddleware{
		logger:               logger,
		httpRequestsTotal:    httpRequestsTotal,
		httpRequestDuration:  httpRequestDuration,
		httpRequestsInFlight: httpRequestsInFlight,
	}
}

// Handle returns a Gin middleware function for metrics collection
func (m MetricsMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Increment in-flight requests
		m.httpRequestsInFlight.WithLabelValues(c.Request.Method, path).Inc()
		defer m.httpRequestsInFlight.WithLabelValues(c.Request.Method, path).Dec()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Record request duration
		m.httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

		// Record total requests
		m.httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()

		// Log metrics for debugging
		m.logger.Debug("Request metrics recorded",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("status", status),
			zap.Float64("duration_seconds", duration),
		)
	}
}

// MetricsHandler returns the Prometheus metrics handler
func (m MetricsMiddleware) MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
