package deception

import (
	"net/http"
	"time"
)

func ServeTarpit(w http.ResponseWriter, r *http.Request) {
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
