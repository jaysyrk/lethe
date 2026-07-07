package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"lethe/internal/deception"
	"lethe/internal/detector"
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
	if detector.IsMalicious(r) {
		deception.HandleMaliciousTraffic(w, r)
		return
	}

	p.proxy.ServeHTTP(w, r)
}
