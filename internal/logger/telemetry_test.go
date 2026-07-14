package logger

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLogThreat(t *testing.T) {
	tempFile := "lethe_intel.log"
	defer os.Remove(tempFile)

	req := httptest.NewRequest("GET", "/badpath", nil)
	req.RemoteAddr = "192.168.1.50:12345"
	req.Header.Set("User-Agent", "sqlmap")

	LogThreat(req, "test_strategy")

	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		t.Fatal("Expected at least one log entry")
	}

	lastLine := lines[len(lines)-1]
	var logEntry ThreatLog
	if err := json.Unmarshal([]byte(lastLine), &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if logEntry.IP != "192.168.1.50:12345" {
		t.Errorf("Expected IP 192.168.1.50:12345, got %s", logEntry.IP)
	}
	if logEntry.Strategy != "test_strategy" {
		t.Errorf("Expected Strategy test_strategy, got %s", logEntry.Strategy)
	}
}
