package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"lethe/internal/proxy"
)

func main() {
	port := flag.String("port", "8080", "Port for the Lethe proxy to listen on")
	target := flag.String("target", "http://localhost:9000", "Target URL of your production backend")
	flag.Parse()

	proxyServer := proxy.NewProxy(*target)

	listenAddr := fmt.Sprintf(":%s", *port)
	log.Printf("Lethe Core Engine started on %s\n", listenAddr)
	log.Printf("Proxying clean traffic to %s\n", *target)

	if err := http.ListenAndServe(listenAddr, proxyServer); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
