package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkProxy_ServeHTTP_Clean(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	p := NewProxy(backend.URL)
	req := httptest.NewRequest("GET", "/", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		p.ServeHTTP(rr, req)
	}
}
