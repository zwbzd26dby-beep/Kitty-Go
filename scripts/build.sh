#!/usr/bin/env sh
set -e

# Build all command binaries.
cd "$(dirname "$0")/.."
mkdir -p bin
for cmd in kitty-cli kitty-repl kitty-api kitty-gui; do
    go build -o "bin/$cmd" "./cmd/$cmd"
done
echo "Build complete: bin/"
