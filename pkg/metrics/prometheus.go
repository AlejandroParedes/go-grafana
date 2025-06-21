package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// PrometheusMetrics provides custom business metrics
type PrometheusMetrics struct {
	logger *zap.Logger
	// Business metrics
	userCreationTotal prometheus.Counter
	userDeletionTotal prometheus.Counter
	userUpdateTotal   prometheus.Counter
	activeUsersGauge  prometheus.Gauge
	userAgeHistogram  prometheus.Histogram
}

// NewPrometheusMetrics creates a new Prometheus metrics instance
func NewPrometheusMetrics(logger *zap.Logger) *PrometheusMetrics {
	// Define business metrics
	userCreationTotal := promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_creation_total",
		Help: "Total number of users created",
	})

	userDeletionTotal := promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_deletion_total",
		Help: "Total number of users deleted",
	})

	userUpdateTotal := promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_update_total",
		Help: "Total number of user updates",
	})

	activeUsersGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_users_total",
		Help: "Total number of active users",
	})

	userAgeHistogram := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "user_age_distribution",
		Help:    "Distribution of user ages",
		Buckets: prometheus.LinearBuckets(0, 10, 13), // 0-120 years in 10-year buckets
	})

	logger.Info("Prometheus metrics initialized")

	return &PrometheusMetrics{
		logger:            logger,
		userCreationTotal: userCreationTotal,
		userDeletionTotal: userDeletionTotal,
		userUpdateTotal:   userUpdateTotal,
		activeUsersGauge:  activeUsersGauge,
		userAgeHistogram:  userAgeHistogram,
	}
}

// RecordUserCreation increments the user creation counter
func (m *PrometheusMetrics) RecordUserCreation() {
	m.userCreationTotal.Inc()
	m.logger.Debug("User creation metric recorded")
}

// RecordUserDeletion increments the user deletion counter
func (m *PrometheusMetrics) RecordUserDeletion() {
	m.userDeletionTotal.Inc()
	m.logger.Debug("User deletion metric recorded")
}

// RecordUserUpdate increments the user update counter
func (m *PrometheusMetrics) RecordUserUpdate() {
	m.userUpdateTotal.Inc()
	m.logger.Debug("User update metric recorded")
}

// SetActiveUsers sets the active users gauge
func (m *PrometheusMetrics) SetActiveUsers(count int64) {
	m.activeUsersGauge.Set(float64(count))
	m.logger.Debug("Active users metric updated", zap.Int64("count", count))
}

// RecordUserAge records a user's age in the histogram
func (m *PrometheusMetrics) RecordUserAge(age int) {
	m.userAgeHistogram.Observe(float64(age))
	m.logger.Debug("User age metric recorded", zap.Int("age", age))
}
