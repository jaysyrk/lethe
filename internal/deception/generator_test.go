package deception

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeInfiniteData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ServeInfiniteData))
	defer ts.Close()


	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)


	if err == nil {
		defer resp.Body.Close()
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			t.Errorf("Expected some data from infinite stream, got none")
		}
	}
}
