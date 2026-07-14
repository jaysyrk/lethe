package detector

import (
	"bytes"
	"io"
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

	for _, values := range r.URL.Query() {
		for _, v := range values {
			if rules.SQLi.MatchString(v) || rules.XSS.MatchString(v) || rules.LFI.MatchString(v) {
				return true
			}
		}
	}

	for key, values := range r.Header {
		for _, v := range values {
			if rules.SQLi.MatchString(v) || rules.XSS.MatchString(v) || rules.LFI.MatchString(v) {
				_ = key
				return true
			}
		}
	}

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 8192))
		if err == nil && len(bodyBytes) > 0 {
			r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes), r.Body))
			
			bodyStr := string(bodyBytes)
			if rules.SQLi.MatchString(bodyStr) || rules.XSS.MatchString(bodyStr) || rules.LFI.MatchString(bodyStr) {
				return true
			}
		}
	}

	return false
}
