package deception

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeGhostHoneypot(t *testing.T) {
	req := httptest.NewRequest("GET", "/.env", nil)
	rr := httptest.NewRecorder()

	ServeGhostHoneypot(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "DB_PASSWORD=fake_password_123\nAWS_KEY=AKIAIOSFODNN7EXAMPLE\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
