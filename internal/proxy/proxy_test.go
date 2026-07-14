package proxy

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"lethe/internal/config"
)

func TestProxy_ServeHTTP_Clean(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend Response"))
	}))
	defer backend.Close()

	p := NewProxy(backend.URL)
	
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	p.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}
	if rr.Body.String() != "Backend Response" {
		t.Errorf("Expected backend response, got %s", rr.Body.String())
	}
}

func TestProxy_ServeHTTP_Malicious(t *testing.T) {
	defer os.Remove("lethe_intel.log")


	tempDir := t.TempDir()
	rulesFile := tempDir + "/test_rules.yaml"
	yamlContent := `threat_signatures:
  bad_paths:
    - ".env"`
	os.WriteFile(rulesFile, []byte(yamlContent), 0644)
	
	config.LoadRules(rulesFile)

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend Response"))
	}))
	defer backend.Close()

	p := NewProxy(backend.URL)
	
	req := httptest.NewRequest("GET", "/.env", nil)
	rr := httptest.NewRecorder()

	p.ServeHTTP(rr, req)


	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK from honeypot, got %d", rr.Code)
	}
	if rr.Body.String() == "Backend Response" {
		t.Errorf("Expected honeypot response, got backend response")
	}
}
