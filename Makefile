ROOT_PATH := $(shell pwd)
GO_FMT := docker compose -f go-fmt.compose.yaml run --rm go-fmt

.PHONY: format test build tidy example seed

format:
	cd example/app && npx oxfmt --write src
	cd example/app && npx oxlint --fix src
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
	cd example/api && npx portless inertia-go go run ./cmd

seed:
	curl -s -X POST http://localhost:8080/dashboard/seed | python3 -m json.tool
