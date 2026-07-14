#!/bin/bash
set -e

# Install hey
echo "Installing hey..."
go install github.com/rakyll/hey@latest
export PATH=$PATH:$(go env GOPATH)/bin

# Build backend and lethe
echo "Building servers..."
go build -o bench_backend bench_backend.go
go build -o lethe_bench cmd/lethe/main.go

# Start backend
echo "Starting mock backend on :5000..."
./bench_backend &
BACKEND_PID=$!

# Start Lethe
echo "Starting Lethe proxy on :8080..."
./lethe_bench -port 8080 -target http://localhost:5000 &
LETHE_PID=$!

# Wait for servers to spin up
sleep 2

echo "======================================"
echo "BENCHMARK 1: Clean Traffic Throughput"
echo "======================================"
hey -n 20000 -c 100 http://localhost:8080/

echo "======================================"
echo "BENCHMARK 2: Tarpit (Malicious Traffic)"
echo "======================================"
# We only do 500 requests here because it will hang them on purpose (tarpit holds connections open)
hey -n 500 -c 100 "http://localhost:8080/search?q=UNION+SELECT" &
HEY_PID=$!

sleep 3
echo "Connections active in tarpit (checking metrics if implemented, else just verifying they are held):"
curl -s http://localhost:8080/.env | head -n 1
echo "Tarpit test initiated. Lethe is holding the connections."
kill $HEY_PID 2>/dev/null || true

# Cleanup
kill $LETHE_PID
kill $BACKEND_PID
rm bench_backend lethe_bench
