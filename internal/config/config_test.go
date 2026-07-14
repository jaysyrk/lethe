package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRules(t *testing.T) {

	tempDir := t.TempDir()
	rulesFile := filepath.Join(tempDir, "test_rules.yaml")

	yamlContent := `threat_signatures:
  sqli: "(?i)(union.*select|select.*from)"
  xss: "(?i)(<script>)"
  lfi: "(?i)(\\.\\./)"
  bad_user_agents: "(?i)(sqlmap)"
  bad_paths:
    - ".env"`

	err := os.WriteFile(rulesFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write temporary test rules file: %v", err)
	}


	LoadRules(rulesFile)

	rules := GetRules()
	if rules == nil {
		t.Fatal("Expected rules to be loaded, got nil")
	}


	if !rules.SQLi.MatchString("UNION SELECT") {
		t.Errorf("Expected SQLi regex to match 'UNION SELECT'")
	}


	if !rules.XSS.MatchString("<script>") {
		t.Errorf("Expected XSS regex to match '<script>'")
	}


	if !rules.LFI.MatchString("../") {
		t.Errorf("Expected LFI regex to match '../'")
	}


	if !rules.BadUAs.MatchString("sqlmap") {
		t.Errorf("Expected BadUAs regex to match 'sqlmap'")
	}


	if len(rules.BadPaths) != 1 || rules.BadPaths[0] != ".env" {
		t.Errorf("Expected BadPaths to contain '.env', got %v", rules.BadPaths)
	}
}

func TestLoadRulesInvalidFile(t *testing.T) {

	mu.Lock()
	CurrentRules = nil
	mu.Unlock()

	LoadRules("non_existent_file.yaml")

	rules := GetRules()
	if rules != nil {
		t.Errorf("Expected rules to be nil when loading invalid file, got %v", rules)
	}
}
