package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CleanRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lethe_clean_requests_total",
		Help: "The total number of clean requests routed to production",
	})

	MaliciousRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lethe_malicious_requests_total",
		Help: "The total number of malicious requests caught by Lethe",
	}, []string{"strategy"})

	JailDrops = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lethe_jail_drops_total",
		Help: "The total number of TCP connections instantly dropped by the IP jail",
	})

	TarpitConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "lethe_tarpit_active_connections",
		Help: "The current number of attackers stuck in the tarpit",
	})
)
