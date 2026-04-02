ROOT_PATH := $(shell pwd)
GO_FMT := docker compose -f go-fmt.compose.yaml run --rm go-fmt

.PHONY: format test build tidy demo seed

format:
	cd demo/app && npx oxfmt --write src
	cd demo/app && npx oxlint --fix src
	go vet $(ROOT_PATH)/...
	$(GO_FMT) format --host-path $(ROOT_PATH)/core
	$(GO_FMT) format --host-path $(ROOT_PATH)/demo/api

test:
	go test $(ROOT_PATH)/...

build:
	go build $(ROOT_PATH)/...

tidy:
	go mod tidy

demo:
	pnpm turbo build --filter=@inertia-go/demo
	cd demo/api && npx portless inertia-go --force go run ./cmd

seed:
	curl -s -X POST http://localhost:8080/dashboard/seed | python3 -m json.tool
