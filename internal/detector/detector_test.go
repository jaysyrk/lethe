package detector

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"lethe/internal/config"
)

func setupTestConfig(t *testing.T) {
	tempDir := t.TempDir()
	rulesFile := filepath.Join(tempDir, "test_rules.yaml")

	yamlContent := `threat_signatures:
  sqli: "(?i)(union.*select)"
  xss: "(?i)(<script>)"
  lfi: "(?i)(\\.\\./)"
  bad_user_agents: "(?i)(sqlmap)"
  bad_paths:
    - ".env"`

	err := os.WriteFile(rulesFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write temporary test rules file: %v", err)
	}

	config.LoadRules(rulesFile)
}

func TestIsMalicious(t *testing.T) {
	setupTestConfig(t)

	tests := []struct {
		name      string
		method    string
		url       string
		userAgent string
		expected  bool
	}{
		{
			name:      "Clean Request",
			method:    "GET",
			url:       "/home",
			userAgent: "Mozilla/5.0",
			expected:  false,
		},
		{
			name:      "Bad User Agent",
			method:    "GET",
			url:       "/home",
			userAgent: "sqlmap/1.0",
			expected:  true,
		},
		{
			name:      "Bad Path",
			method:    "GET",
			url:       "/admin/.env",
			userAgent: "Mozilla/5.0",
			expected:  true,
		},
		{
			name:      "SQLi in Path",
			method:    "GET",
			url:       "/union%20select",
			userAgent: "Mozilla/5.0",
			expected:  true,
		},
		{
			name:      "SQLi in Query",
			method:    "GET",
			url:       "/search?q=UNION+SELECT",
			userAgent: "Mozilla/5.0",
			expected:  true,
		},
		{
			name:      "XSS in Query",
			method:    "GET",
			url:       "/search?q=<script>alert(1)</script>",
			userAgent: "Mozilla/5.0",
			expected:  true,
		},
		{
			name:      "LFI in Query",
			method:    "GET",
			url:       "/view?file=../../etc/passwd",
			userAgent: "Mozilla/5.0",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Set("User-Agent", tt.userAgent)

			if got := IsMalicious(req); got != tt.expected {
				t.Errorf("IsMalicious() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestIsMalicious_NoRulesLoaded(t *testing.T) {

}
