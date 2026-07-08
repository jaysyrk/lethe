package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"

	"lethe/internal/config"
	"lethe/internal/proxy"
)

func main() {
	port := flag.String("port", "8080", "Port for Lethe to listen on (use 443 if using a domain)")
	target := flag.String("target", "http://localhost:9000", "Target URL of your production backend")
	domain := flag.String("domain", "", "Your domain (e.g., 1cdn99.com) for auto-TLS. Leave empty for HTTP.")
	flag.Parse()

	config.LoadRules("rules.yaml")
	go config.WatchRules("rules.yaml")

	go func() {
		log.Println("Starting Prometheus metrics on :9090/metrics")
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", metricsMux); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

	proxyServer := proxy.NewProxy(*target)
	listenAddr := fmt.Sprintf(":%s", *port)

	log.Printf("Proxying clean traffic to %s\n", *target)

	if *domain != "" {
		log.Printf("Lethe Core Engine started on %s with auto-TLS for %s\n", listenAddr, *domain)

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*domain),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr:    listenAddr,
			Handler: proxyServer,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("TLS Server failed: %v", err)
		}
	} else {
		log.Printf("Lethe Core Engine started on %s (HTTP Only)\n", listenAddr)
		if err := http.ListenAndServe(listenAddr, proxyServer); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
