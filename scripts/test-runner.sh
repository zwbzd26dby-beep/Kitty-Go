#!/usr/bin/env sh
set -e

# Run the full test suite.
cd "$(dirname "$0")/.."
go build ./...
go test ./...
