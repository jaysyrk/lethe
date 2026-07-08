package detector

import (
	"net/http"
	"strings"

	"lethe/internal/config"
)

func IsMalicious(r *http.Request) bool {
	rules := config.GetRules()
	if rules == nil {
		return false
	}

	if rules.BadUAs.MatchString(r.UserAgent()) {
		return true
	}

	path := r.URL.Path
	for _, badPath := range rules.BadPaths {
		if strings.HasSuffix(path, badPath) {
			return true
		}
	}

	if rules.LFI.MatchString(path) || rules.SQLi.MatchString(path) || rules.XSS.MatchString(path) {
		return true
	}

	query := r.URL.RawQuery
	if rules.SQLi.MatchString(query) || rules.XSS.MatchString(query) || rules.LFI.MatchString(query) {
		return true
	}

	return false
}
