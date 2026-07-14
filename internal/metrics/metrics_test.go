package metrics

import (
	"testing"
)

func TestMetrics(t *testing.T) {

	CleanRequests.Inc()
	MaliciousRequests.WithLabelValues("test_strategy").Inc()
	JailDrops.Inc()
	TarpitConnections.Inc()
	TarpitConnections.Dec()
}
