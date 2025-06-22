package metrics

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.uber.org/zap"
)

func TestPrometheusMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	logger := zap.NewNop()
	metrics := NewPrometheusMetrics(logger, reg)

	metrics.RecordUserCreation()
	metrics.RecordUserDeletion()
	metrics.RecordUserUpdate()
	metrics.SetActiveUsers(42)
	metrics.RecordUserAge(30)

	expected := `
		# HELP active_users_total Total number of active users
		# TYPE active_users_total gauge
		active_users_total 42
		# HELP user_age_distribution Distribution of user ages
		# TYPE user_age_distribution histogram
		user_age_distribution_bucket{le="0"} 0
		user_age_distribution_bucket{le="10"} 0
		user_age_distribution_bucket{le="20"} 0
		user_age_distribution_bucket{le="30"} 1
		user_age_distribution_bucket{le="40"} 1
		user_age_distribution_bucket{le="50"} 1
		user_age_distribution_bucket{le="60"} 1
		user_age_distribution_bucket{le="70"} 1
		user_age_distribution_bucket{le="80"} 1
		user_age_distribution_bucket{le="90"} 1
		user_age_distribution_bucket{le="100"} 1
		user_age_distribution_bucket{le="110"} 1
		user_age_distribution_bucket{le="120"} 1
		user_age_distribution_bucket{le="+Inf"} 1
		user_age_distribution_sum 30
		user_age_distribution_count 1
		# HELP user_creation_total Total number of users created
		# TYPE user_creation_total counter
		user_creation_total 1
		# HELP user_deletion_total Total number of users deleted
		# TYPE user_deletion_total counter
		user_deletion_total 1
		# HELP user_update_total Total number of user updates
		# TYPE user_update_total counter
		user_update_total 1
	`
	err := testutil.CollectAndCompare(reg, strings.NewReader(expected))
	if err != nil {
		t.Errorf("unexpected metrics collection result:\n%v", err)
	}
}
