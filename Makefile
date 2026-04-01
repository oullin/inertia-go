GO_FMT := docker compose -f go-fmt.compose.yaml run --rm go-fmt

.PHONY: format format-check test lint build tidy example

format:
	$(GO_FMT) format ./core
	cd example && npx oxfmt --write resources
	cd example && npx oxlint --fix resources

format-check:
	$(GO_FMT) check .

test:
	go test ./...

lint:
	go vet ./...

build:
	go build ./...

tidy:
	go mod tidy

example:
	pnpm turbo build --filter=@inertia-go/example
	cd example && go run main.go
