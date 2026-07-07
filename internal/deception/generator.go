package deception

import (
	"math/rand"
	"net/http"
)

func ServeInfiniteData(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	buf := make([]byte, 1024)
	for {
		rand.Read(buf)
		_, err := w.Write(buf)
		if err != nil {
			break
		}
		flusher.Flush()
	}
}
