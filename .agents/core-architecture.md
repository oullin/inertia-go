# Core Architecture

## Scope

This guide covers the non-obvious package boundaries and architectural rules for [`core/`](../core), the Go server-side adapter module at `github.com/oullin/inertia-go/core`.

## Package Map

```text
  cryptox  wayfinder  validation
     |         |          |
     +----+    |    +-----+
          v    v    v
httpx <-- props  response  config <-- i18n
  ^         |       |        |
  |         v       v        v
  +------ inertia <---- middleware
              ^
         +----+----+
       flash     assert
```

- `httpx` is the shared foundation for headers, request helpers, context helpers, and shared types.
- `props`, `response`, and `config` sit above `httpx` and should stay narrow in responsibility.
- `inertia` is the hub package that ties rendering, middleware integration, config, and prop resolution together.
- `flash` and `assert` depend on `inertia`; `inertia` must not import them back.
- `cryptox` and `wayfinder` are standalone utilities and should remain free of core-internal coupling.

## Dependency Rules

- Avoid circular imports. If a new dependency would force `httpx`, `props`, `response`, `config`, `cryptox`, or `wayfinder` to import `inertia`, the design is wrong.
- Prefer `httpx` for cross-cutting request metadata that middleware needs to share without pulling in render-specific state.
- Keep new external dependencies out of packages that are currently stdlib-only unless the benefit is clear and repo-wide.

## Context Boundaries

There are two context key systems and they serve different layers:

- `httpx.ctxKey` is for cross-cutting request concerns such as CSRF tokens, locale, and precognition flags.
- `inertia.contextKey` is for render-scoped data such as props, template data, validation errors, head metadata, and history flags.

Use `httpx` helpers in middleware that must not import `inertia`. Use `inertia` helpers in handlers or higher-level integration code. The `inertia` package re-exports `SetCSRFToken`, `SetPrecognition`, and `SetLocale` for convenience when the higher-level package boundary is acceptable.

## Config And Protocol Constraints

- Config loading flows from defaults to YAML file to env vars with the `INERTIA_*` prefix.
- Head meta entries with an empty `Content` value act as placeholder slots and are omitted from rendered HTML.
- [`response.Page`](../core/response/page.go) must keep JSON tags aligned with the Inertia.js page object protocol. Do not rename tags or add `omitempty` to required fields without verifying the client contract.

## External Dependencies

`core/` currently relies on only two non-stdlib packages:

- `github.com/go-playground/validator/v10` in `validation/`
- `github.com/spf13/viper` in `config/`

Treat the small dependency surface as part of the package design.
