package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"lethe/internal/deception"
	"lethe/internal/detector"
	"lethe/internal/jail"
	"lethe/internal/metrics"
)

type Proxy struct {
	targetURL *url.URL
	proxy     *httputil.ReverseProxy
}

func NewProxy(target string) *Proxy {
	url, _ := url.Parse(target)
	return &Proxy{
		targetURL: url,
		proxy:     httputil.NewSingleHostReverseProxy(url),
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	if jail.IsBanned(ip) {
		metrics.JailDrops.Inc()
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		return
	}

	if detector.IsMalicious(r) {
		jail.BanIP(ip)
		deception.HandleMaliciousTraffic(w, r)
		return
	}

	metrics.CleanRequests.Inc()
	p.proxy.ServeHTTP(w, r)
}
