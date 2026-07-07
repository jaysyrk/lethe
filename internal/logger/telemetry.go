package logger

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type ThreatLog struct {
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	UserAgent string    `json:"user_agent"`
	Strategy  string    `json:"deception_strategy"`
}

func LogThreat(r *http.Request, strategy string) {
	entry := ThreatLog{
		Timestamp: time.Now(),
		IP:        r.RemoteAddr,
		Method:    r.Method,
		Path:      r.URL.Path,
		UserAgent: r.UserAgent(),
		Strategy:  strategy,
	}

	logBytes, err := json.Marshal(entry)
	if err == nil {
		f, err := os.OpenFile("lethe_intel.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(string(logBytes) + "\n")
		}
		log.Println(string(logBytes))
	}
}
