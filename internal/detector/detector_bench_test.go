package detector

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"lethe/internal/config"
)

func setupBenchConfig(b *testing.B) {
	tempDir := b.TempDir()
	rulesFile := filepath.Join(tempDir, "bench_rules.yaml")

	yamlContent := `threat_signatures:
  sqli: "(?i)(union.*select)"
  xss: "(?i)(<script>)"
  lfi: "(?i)(\\.\\./)"
  bad_user_agents: "(?i)(sqlmap)"
  bad_paths:
    - ".env"`

	err := os.WriteFile(rulesFile, []byte(yamlContent), 0644)
	if err != nil {
		b.Fatalf("Failed to write temporary rules file: %v", err)
	}

	config.LoadRules(rulesFile)
}

func BenchmarkIsMalicious_Clean(b *testing.B) {
	setupBenchConfig(b)
	req := httptest.NewRequest("GET", "/home/dashboard?id=123", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsMalicious(req)
	}
}

func BenchmarkIsMalicious_Malicious(b *testing.B) {
	setupBenchConfig(b)
	req := httptest.NewRequest("GET", "/search?q=UNION+SELECT", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsMalicious(req)
	}
}
