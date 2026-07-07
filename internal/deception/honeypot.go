package deception

import (
	"net/http"
)

func ServeGhostHoneypot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DB_PASSWORD=fake_password_123\nAWS_KEY=AKIAIOSFODNN7EXAMPLE\n"))
}
