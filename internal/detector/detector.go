package detector

import (
	"net/http"
	"regexp"
	"strings"
)

var (
	sqliRegex  = regexp.MustCompile(`(?i)(union.*select|select.*from|insert.*into|drop\s+table)`)
	xssRegex   = regexp.MustCompile(`(?i)(<script>|javascript:|onerror=)`)
	lfiRegex   = regexp.MustCompile(`(?i)(\.\./\.\./|/etc/passwd|/windows/win\.ini)`)
	badUARegex = regexp.MustCompile(`(?i)(sqlmap|nmap|nikto|masscan|zgrab)`)
)

func IsMalicious(r *http.Request) bool {
	userAgent := r.UserAgent()
	if badUARegex.MatchString(userAgent) {
		return true
	}

	path := r.URL.Path
	if strings.HasSuffix(path, ".env") || strings.HasSuffix(path, ".git") || strings.HasSuffix(path, "config.json") {
		return true
	}

	if lfiRegex.MatchString(path) || sqliRegex.MatchString(path) || xssRegex.MatchString(path) {
		return true
	}

	query := r.URL.RawQuery
	if sqliRegex.MatchString(query) || xssRegex.MatchString(query) || lfiRegex.MatchString(query) {
		return true
	}

	return false
}
