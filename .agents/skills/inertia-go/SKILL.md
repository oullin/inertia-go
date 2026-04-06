---
name: inertia-go
description: Use when working on the core/ package of inertia-go. Covers architecture, package boundaries, prop system, middleware, context conventions, and testing patterns for the Inertia.js Go adapter.
---

# Inertia Go

Guidance for working on the `core/` module of inertia-go, the Go server-side adapter for Inertia.js.

## When to Use

- Adding or modifying packages in `core/`
- Working with the prop system (`Always`, `Defer`, `Once`, `Optional`, `Merge`, `Scroll`)
- Touching middleware, context helpers, or the response pipeline
- Writing or fixing tests for core packages
- Changing config loading or wire format

## Quick Reference

Run from the repo root:

```sh
cd core && go build ./...
cd core && go test ./...
cd core && go test -race ./...
cd core && go vet ./...
make format
make tidy
```

Formatting runs through Docker Compose via `go-fmt.compose.yaml`. Do not use `gofmt` directly.

## Package Map

```text
httpx         <- leaf foundation (types, headers, context, form parsing)
props         <- prop wrappers + resolver (imports httpx)
response      <- page object + rendering (imports httpx)
config        <- YAML/env config loading via viper (imports httpx)
cryptox       <- AES-256-CBC encryption (standalone)
wayfinder     <- named route registry + codegen (standalone)
validation    <- struct validation via go-playground/validator (imports httpx)
i18n          <- URL-prefix locale detection (imports config, httpx)
middleware    <- Inertia protocol, CSRF, precognition (imports httpx, config, cryptox)
inertia       <- hub: ties props, response, middleware, config together
flash         <- flash messages (imports inertia and calls inertia.SetProp)
assert        <- test helpers (imports inertia)
```

Dependency rule: `httpx` is the shared foundation. `inertia` is the hub. Do not create circular imports.

## Architecture Rules

### Prop Precedence

In `inertia/inertia.go` -> `mergeProps()`:

```text
shared < context (SetProp) < render-time < validation errors
```

Later sources override earlier ones for the same key.

### Prop Wrapper Nesting

Wrappers compose by struct nesting such as `Always(Defer(value))`. The resolver in `props/resolver.go` -> `walkPropChain()` walks the chain, and the outermost wrapper wins when duplicate traits appear.

### Two Context Key Systems

1. `httpx.ctxKey` for cross-cutting concerns such as CSRF tokens, locale, and precognition.
2. `inertia.contextKey` for render-scoped data such as props, template data, validation errors, head metadata, and history flags.

`inertia` re-exports `SetCSRFToken`, `SetPrecognition`, and `SetLocale` as convenience wrappers over `httpx`.

### Resolver Pipeline

`props.Resolve()` runs two stages:

1. Filter props and collect metadata from request headers.
2. Evaluate only the props that survived filtering.

Lazy props must not execute during filtering.

### Config Layering

`config/` loads in this order: defaults, YAML file, then `INERTIA_*` env vars. Empty `Content` values in head meta tags act as placeholders and are omitted from rendered HTML.

### Wire Format

`response.Page` JSON tags must match the Inertia.js page object spec. Do not rename tags or add `omitempty` to required fields without checking the client contract.

## Testing Conventions

- Stdlib only: `testing` and `net/http/httptest`
- Use `t.Parallel()` where safe
- Use `httptest.ResponseRecorder` for HTTP assertions
- Use `assert.AssertFromHandler()` when tests should exercise middleware and serialized Inertia responses
- Test helpers should accept `testing.TB` and call `t.Helper()`
- `validation.Validate()` resolves field names as `json` tag, then `form` tag, then the Go field name

## Must

- Preserve package boundaries before adding imports
- Run `make format` before committing
- Run `cd core && go test -race ./...` for CI-aligned coverage
- Keep `response.Page` JSON tags in sync with the protocol
- Use functional options for constructor configuration
- Use context for per-request data and mutexes for shared state

## Must Not

- Import `inertia` from `httpx`, `props`, `response`, `config`, `cryptox`, or `wayfinder`
- Add external dependencies to packages that are currently stdlib-only without a strong reason
- Use `omitempty` on required `response.Page` fields
- Skip the filter stage when resolving props
- Use `gofmt` directly instead of `make format`

## External Dependencies

`core/` currently relies on only two non-stdlib packages:

- `go-playground/validator/v10` in `validation/`
- `spf13/viper` in `config/`

Keep the dependency surface small.
