GO_FMT := $(shell go env GOPATH)/bin/fmt

.PHONY: format format-check test lint build example

format:
	cd core && $(GO_FMT) format .

format-check:
	cd core && $(GO_FMT) check .

test:
	cd core && go test ./...

lint:
	cd core && go vet ./...

build:
	cd core && go build ./...

example:
	pnpm turbo build --filter=@inertia-go/example
	cd example && go run main.go
