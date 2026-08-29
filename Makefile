.PHONY: build repl cli api test unit integration clean fmt vet

build:
	go build ./...

repl:
	go run ./cmd/kitty-repl

cli:
	go run ./cmd/kitty-cli "Smoke test"

api:
	go run ./cmd/kitty-api

test:
	go test ./...

unit:
	go test ./tests/unit/... ./pkg/... ./cmd/...

integration:
	go test ./tests/integration/... ./tests/cli/...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin dist
