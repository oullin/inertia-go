ROOT_PATH := $(shell pwd)
GO_FMT := docker compose -f go-fmt.compose.yaml run --rm go-fmt

.PHONY: format test build tidy example

format:
	cd example/app && pnpm exec oxfmt --write src
	cd example/app && pnpm exec oxlint --fix src
	go vet $(ROOT_PATH)/...
	$(GO_FMT) format --host-path $(ROOT_PATH)/core
	$(GO_FMT) format --host-path $(ROOT_PATH)/example/api

test:
	go test $(ROOT_PATH)/...

build:
	go build $(ROOT_PATH)/...

tidy:
	go mod tidy

example:
	pnpm turbo build --filter=@inertia-go/example
	cd example/api && pnpm exec portless --force inertia-go go run ./cmd
