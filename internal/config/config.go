package config

import (
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type RulesYAML struct {
	ThreatSignatures struct {
		SQLi          string   `yaml:"sqli"`
		XSS           string   `yaml:"xss"`
		LFI           string   `yaml:"lfi"`
		BadUserAgents string   `yaml:"bad_user_agents"`
		BadPaths      []string `yaml:"bad_paths"`
	} `yaml:"threat_signatures"`
}

type CompiledRules struct {
	SQLi     *regexp.Regexp
	XSS      *regexp.Regexp
	LFI      *regexp.Regexp
	BadUAs   *regexp.Regexp
	BadPaths []string
}

var (
	CurrentRules *CompiledRules
	mu           sync.RWMutex
)

func LoadRules(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read rules file: %v", err)
		return
	}

	var raw RulesYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Printf("Failed to parse YAML: %v", err)
		return
	}

	compiled := &CompiledRules{
		SQLi:     regexp.MustCompile(raw.ThreatSignatures.SQLi),
		XSS:      regexp.MustCompile(raw.ThreatSignatures.XSS),
		LFI:      regexp.MustCompile(raw.ThreatSignatures.LFI),
		BadUAs:   regexp.MustCompile(raw.ThreatSignatures.BadUserAgents),
		BadPaths: raw.ThreatSignatures.BadPaths,
	}

	mu.Lock()
	CurrentRules = compiled
	mu.Unlock()
}

func GetRules() *CompiledRules {
	mu.RLock()
	defer mu.RUnlock()
	return CurrentRules
}

func WatchRules(filename string) {
	for {
		time.Sleep(10 * time.Second)
		LoadRules(filename)
	}
}
