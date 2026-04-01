GO_FMT := $(shell go env GOPATH)/bin/fmt

.PHONY: format format-check test lint build example

format:
	$(GO_FMT) format .

format-check:
	$(GO_FMT) check .

test:
	go test ./...

lint:
	go vet ./...

build:
	go build ./...

example:
	pnpm turbo build --filter=@inertia-go/example
	cd example && go run main.go
