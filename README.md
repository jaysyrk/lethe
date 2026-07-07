# Lethe Core Engine

Lethe is a high-performance Deception Proxy and Honeypot built in Go. It acts as a reverse proxy that routes clean traffic to a production backend while silently migrating malicious traffic into an isolated deception layer to waste attacker resources and harvest threat intelligence.

## Architecture

1. **Threat Detection Engine**: Inspects incoming traffic for known malicious signatures (SQLi, XSS, LFI) and malicious user agents.
2. **Production Route**: Clean traffic is silently forwarded to the real backend.
3. **Isolated Deception Layer**: Malicious traffic is trapped using various strategies:
   - **TCP Tarpit (1-Byte Window Stalling)**: Holds the attacker's connection open and trickles 1 byte of data every 10 seconds to exhaust their resources.
   - **Infinite Data Stream Generator**: Floods data exfiltration attempts (e.g., automated POST requests) with endless streams of garbage data.
   - **Ephemeral Ghost Honeypot**: Serves synthetic schemas and fake flag files (like fake `.env` files) to poison the attacker's intel.
4. **Zetetic Log Chain**: Continuously fingerprints attackers and logs their IP, tooling, and behavior to a structured telemetry file (`lethe_intel.log`).

## How to Use It (Integration)

Lethe is designed to act as a **drop-in Edge Proxy** (similar to NGINX, HAProxy, or Traefik) that sits directly in front of your actual application. It can protect any tech stack (Node.js, Python, Rust, PHP, etc.) because it operates at the network/HTTP level.

1. **Hide your real application**: Bind your real backend server (e.g., a Django or Express app) to `localhost` or an internal network IP so it cannot be accessed directly from the internet.
2. **Configure Lethe**: Start Lethe using command-line flags to point it at your hidden backend.
3. **Expose Lethe**: Run Lethe on your public edge port (e.g., port `80` or `443`). 

Normal users will browse your site seamlessly, completely unaware of the proxy. However, the moment a scanner or attacker fires a malicious payload, their connection is instantly hijacked away from your backend and thrown into the deception layer.

## Getting Started

Initialize the module and compile the proxy:
```bash
go mod tidy
go build -o lethe cmd/lethe/main.go
```

Start the proxy using CLI flags:
```bash
# Example: Protect a python app running on port 5000, and expose Lethe on port 80 (HTTP)
./lethe -port 80 -target http://localhost:5000
```
*If no flags are provided, it defaults to listening on `:8080` and proxying to `http://localhost:9000`.*

## Testing the Deception Layer

Assuming Lethe is listening on port `8080`:

**1. Test the Ghost Honeypot:**
```bash
curl -v http://localhost:8080/.env
```

**2. Test the SQL Injection Trap (Tarpit):**
```bash
curl -v "http://localhost:8080/search?q=UNION+SELECT+*+FROM+users"
```

**3. Test the Infinite Data Stream (POST request):**
```bash
curl -v -X POST -H "User-Agent: sqlmap" http://localhost:8080/
```
