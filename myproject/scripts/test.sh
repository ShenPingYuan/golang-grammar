#!/bin/bash
set -e

echo "==> Running tests..."
go test ./... -race -cover -count=1 -timeout 60s

echo "==> Running vet..."
go vet ./...

echo "==> All checks passed."