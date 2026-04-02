ROOT_PATH := $(shell pwd)
GO_FMT := docker compose -f go-fmt.compose.yaml run --rm go-fmt

.PHONY: format test build tidy demo seed

format:
	cd demo/app && npx oxfmt --write src
	cd demo/app && npx oxlint --fix src
	cd core && go vet ./...
	cd demo/api && go vet ./...
	$(GO_FMT) format --host-path $(ROOT_PATH)/core
	$(GO_FMT) format --host-path $(ROOT_PATH)/demo/api

test:
	cd core && go test -coverprofile=$(ROOT_PATH)/storage/.cache/coverage.out ./...
	cd core && go test -coverprofile=$(ROOT_PATH)/storage/.cache/inertia_cov.out ./inertia/...
	cd core && go test -coverprofile=$(ROOT_PATH)/storage/.cache/mw_cov.out ./middleware/...
	cd core && go test -coverprofile=$(ROOT_PATH)/storage/.cache/props_cov.out ./props/...
	cd core && go test -coverprofile=$(ROOT_PATH)/storage/.cache/assert_cov.out ./assert/...
	cd demo/api && go test ./...

build:
	cd core && go build ./...
	cd demo/api && go build -o $(ROOT_PATH)/storage/dist/api/demo-api ./cmd

tidy:
	cd core && go mod tidy
	cd demo/api && go mod tidy
	go work sync

demo:
	pnpm turbo build --filter=@inertia-go/demo
	cd demo/api && npx portless inertia-go --force go run ./cmd

seed:
	curl -s -X POST http://localhost:8080/dashboard/seed | python3 -m json.tool
