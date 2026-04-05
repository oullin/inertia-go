# core — Inertia.js Go Adapter

Go server-side adapter for [Inertia.js](https://inertiajs.com). Module: `github.com/oullin/inertia-go/core`. Requires Go 1.26+.

## Quick Reference

All commands run from the **repo root**:

| Task | Command |
|------|---------|
| Build | `cd core && go build ./...` |
| Test | `cd core && go test ./...` |
| Test (match CI) | `cd core && go test -race ./...` |
| Vet | `cd core && go vet ./...` |
| Format | `make format` (requires Docker) |
| Tidy | `make tidy` |

> **Formatting uses Docker Compose** (`go-fmt.compose.yaml`). Running `gofmt` directly will not match CI.

## Package Dependency Graph

```
  cryptox  wayfinder  validation
     │         │          │
     └────┐    │    ┌─────┘
          ▼    ▼    ▼
httpx ◄── props  response  config ◄── i18n
  ▲         │       │        │
  │         ▼       ▼        ▼
  └────── inertia ◄──── middleware
              ▲
         ┌────┴────┐
       flash     assert
```

- **httpx** — leaf foundation; nearly everything imports it (types, headers, context, form parsing)
- **cryptox**, **wayfinder** — standalone; no core-internal imports
- **inertia** — hub package; ties props, response, middleware, config together
- **flash** — depends on inertia (calls `inertia.SetProp()`), not the reverse

## Architecture — Non-Obvious Details

### Prop precedence

Defined in `inertia/inertia.go` → `mergeProps()` via `props.MergeAll()`:

```
shared < context (SetProp) < render-time < validation errors
```

### Prop wrapper nesting

Wrappers compose by struct nesting: `Always(Defer(value))`. The resolver (`props/resolver.go` → `walkPropChain()`) walks the chain — **outermost wrapper wins** when duplicate traits appear.

### Two context key systems

- **`httpx.ctxKey`** — cross-cutting concerns (CSRF token, locale, precognition). Use in middleware that must not import `inertia`.
- **`inertia.contextKey`** — render-scoped data (props, template data, validation errors, head, history flags). Use in application handlers.
- `inertia` re-exports `SetCSRFToken`, `SetPrecognition`, `SetLocale` as convenience wrappers.

### Resolver two-stage pipeline

`props/resolver.go` → `Resolve()`:

1. **Filter** — header-driven inclusion (`X-Inertia-Partial-Data`, etc.) + metadata collection (merge keys, deferred groups, once tracking)
2. **Evaluate** — unwrap `Proper`/`TryProper` interfaces, call lazy `func() any` closures

Lazy props only execute if they survive filtering. This is performance-critical.

### Config layering

`config/` packages load: **defaults → YAML file → env vars** (`INERTIA_*` prefix). Empty `Content` in meta tags serves as a placeholder slot and is excluded from rendered HTML.

### Wire format

`response.Page` JSON tags must match the [Inertia.js page object spec](https://inertiajs.com/the-protocol#the-page-object) exactly. Do not add `omitempty` to required fields or rename JSON tags without checking the JS client.

## Testing Conventions

- **Stdlib only** — `testing` + `net/http/httptest`. No testify, no gomock.
- Tests use `t.Parallel()`.
- HTTP tests use the `httptest.ResponseRecorder` pattern.
- `assert.AssertFromHandler()` runs a handler through Inertia middleware and captures the response.
- `assert.AssertableInertia` provides: `AssertComponent`, `AssertPropEquals`, `AssertHasProp`, `AssertMissingProp`.
- Test helpers accept `testing.TB` and call `t.Helper()`.
- `validation.Validate()` resolves field names via: `json` tag → `form` tag → Go field name.

## External Dependencies

Only two:

- **go-playground/validator/v10** — used in `validation/`
- **spf13/viper** — used in `config/`

Everything else is stdlib.
