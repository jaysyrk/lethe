package deception

import (
	"net/http"
	"strings"

	"lethe/internal/logger"
	"lethe/internal/metrics"
)

func HandleMaliciousTraffic(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".env") {
		logger.LogThreat(r, "synthetic_schemas")
		metrics.MaliciousRequests.WithLabelValues("synthetic_schemas").Inc()
		ServeGhostHoneypot(w, r)
		return
	}

	if r.Method == http.MethodPost {
		logger.LogThreat(r, "infinite_data_stream")
		metrics.MaliciousRequests.WithLabelValues("infinite_data_stream").Inc()
		ServeInfiniteData(w, r)
		return
	}

	logger.LogThreat(r, "tcp_tarpit")
	metrics.MaliciousRequests.WithLabelValues("tcp_tarpit").Inc()
	ServeTarpit(w, r)
}
