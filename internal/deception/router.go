package deception

import (
	"net/http"
	"strings"

	"lethe/internal/logger"
)

func HandleMaliciousTraffic(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".env") {
		logger.LogThreat(r, "synthetic_schemas")
		ServeGhostHoneypot(w, r)
		return
	}

	if r.Method == http.MethodPost {
		logger.LogThreat(r, "infinite_data_stream")
		ServeInfiniteData(w, r)
		return
	}

	logger.LogThreat(r, "tcp_tarpit")
	ServeTarpit(w, r)
}
