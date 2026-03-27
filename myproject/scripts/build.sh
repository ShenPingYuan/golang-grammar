#!/bin/bash
set -e

echo "==> Building binaries..."
mkdir -p bin

go build -o bin/server    ./cmd/server
go build -o bin/grpcserver ./cmd/grpcserver
go build -o bin/worker    ./cmd/worker
go build -o bin/scheduler ./cmd/scheduler
go build -o bin/migrate   ./cmd/migrate
go build -o bin/cli       ./cmd/cli

echo "==> Done. Binaries in ./bin/"
ls -la bin/