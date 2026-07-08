package deception

import (
	"net/http"
	"time"

	"lethe/internal/metrics"
)

var tarpitSemaphore = make(chan struct{}, 1000)

func ServeTarpit(w http.ResponseWriter, r *http.Request) {
	select {
	case tarpitSemaphore <- struct{}{}:
		defer func() { <-tarpitSemaphore }()
	default:
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		} else {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
		return
	}

	metrics.TarpitConnections.Inc()
	defer metrics.TarpitConnections.Dec()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	for {
		_, err := w.Write([]byte("0"))
		if err != nil {
			break
		}
		flusher.Flush()
		time.Sleep(10 * time.Second)
	}
}
