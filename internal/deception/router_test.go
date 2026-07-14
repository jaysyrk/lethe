package deception

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleMaliciousTraffic_Honeypot(t *testing.T) {
	req := httptest.NewRequest("GET", "/.env", nil)
	rr := httptest.NewRecorder()

	HandleMaliciousTraffic(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected OK status, got %d", rr.Code)
	}
}
