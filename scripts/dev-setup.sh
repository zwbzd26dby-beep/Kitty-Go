#!/usr/bin/env sh
set -e

# Development setup for Kitty-Go (Go toolchain check).
command -v go >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
echo "Go: $(go version)"
go env GOPATH GOOS GOARCH
echo "Setup OK"
