ROOT_PATH := $(shell pwd)
GO_FMT := docker compose -f go-fmt.compose.yaml run --rm go-fmt

.PHONY: format test build tidy demo seed

format:
	cd demo/app && npx oxfmt --write src
	cd demo/app && npx oxlint --fix src
	go vet ./core/... ./demo/api/...
	$(GO_FMT) format --host-path $(ROOT_PATH)/core
	$(GO_FMT) format --host-path $(ROOT_PATH)/demo/api

test:
	go test ./core/... ./demo/api/...

build:
	go build ./core/... ./demo/api/...

tidy:
	cd core && go mod tidy
	cd demo/api && go mod tidy
	go work sync

demo:
	pnpm turbo build --filter=@inertia-go/demo
	cd demo/api && npx portless inertia-go --force go run ./cmd

seed:
	curl -s -X POST http://localhost:8080/dashboard/seed | python3 -m json.tool
