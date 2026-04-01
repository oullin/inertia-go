GO_FMT := $(shell go env GOPATH)/bin/fmt

.PHONY: format format-check test lint build tidy example

format:
	cd core && $(GO_FMT) format .
	cd example && npx oxfmt --write resources
	cd example && npx oxlint --fix resources

format-check:
	cd core && $(GO_FMT) check .

test:
	cd core && go test ./...

lint:
	cd core && go vet ./...

build:
	cd core && go build ./...

tidy:
	cd core && go mod tidy
	cd example && go mod tidy

example:
	pnpm turbo build --filter=@inertia-go/example
	cd example && go run main.go
