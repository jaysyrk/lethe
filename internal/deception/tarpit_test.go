package deception

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeTarpit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ServeTarpit))
	defer ts.Close()


	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	
	if err == nil {
		defer resp.Body.Close()
		buf := make([]byte, 1)
		n, _ := resp.Body.Read(buf)
		if n > 0 && string(buf[:n]) != "0" {
			t.Errorf("Expected '0' from tarpit, got %s", string(buf[:n]))
		}
	}
}
