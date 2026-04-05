# inertia-go

[![Go Reference](https://pkg.go.dev/badge/github.com/oullin/inertia-go.svg)](https://pkg.go.dev/github.com/oullin/inertia-go)
[![go 1.26](https://img.shields.io/badge/go-1.26-00ADD8.svg)](https://go.dev/dl/)
[![CI](https://github.com/oullin/inertia-go/actions/workflows/ci.yml/badge.svg)](https://github.com/oullin/inertia-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oullin/inertia-go/branch/main/graph/badge.svg)](https://codecov.io/gh/oullin/inertia-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/oullin/inertia-go)](https://goreportcard.com/report/github.com/oullin/inertia-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go server-side adapter for [Inertia.js](https://inertiajs.com). Build modern single-page apps using Vue, React, or Svelte with a Go backend -- no client-side routing or API layer required.

## Why

Go teams building modern UIs face two bad options: build a full REST/GraphQL API plus a separate SPA with its own routing, state management, and deployment pipeline -- or stick with server-rendered templates and accept the UX tradeoff.

Inertia removes this choice. Your Go handlers render frontend components directly. No API endpoints to maintain, no client-side router to configure, no request/response contracts to keep in sync. A handler calls `Render("Users/Index", props)` and the client gets a Vue, React, or Svelte page with exactly the data it needs. On first visit it's server-rendered HTML; on navigation it's a JSON payload that swaps the component in place. Same handler, both cases.

The result: SPA-quality UX with the simplicity of a traditional server-rendered app. One router, one source of truth, one deployment.

## Who it's for

- **Go backend developers** who want a modern frontend without building and maintaining a separate API layer.
- **Teams moving to Go** from Laravel, Rails, or Django stacks that already use Inertia and want to keep the same workflow.
- **Fullstack developers** who prefer server-driven control over routing, authorization, and data loading -- but don't want to give up single-page app UX.

## Install

```bash
go get github.com/oullin/inertia-go/core
```

**Requires Go 1.26+**

## Development Formatting

From the repo root, format both Go trees with the repo-local Compose file:

```bash
docker compose -f go-fmt.compose.yaml run --rm go-fmt
```

Check both Go trees explicitly with repeated `--host-path` values:

```bash
docker compose -f go-fmt.compose.yaml run --rm go-fmt check --host-path "$PWD/core" --host-path "$PWD/demo/api"
```

The Compose service bakes in the default `format` command for `core` and `demo/api`. When you override that command with `check` or another subcommand, you must pass both `--host-path` arguments again, matching the upstream consumer workflow.

## Quick Start

```go
package main

import (
    "log"
    "net/http"

    "github.com/oullin/inertia-go/core/httpx"
    "github.com/oullin/inertia-go/core/inertia"
)

func main() {
    i, err := inertia.New(`<!DOCTYPE html>
<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}<script src="/app.js"></script></body>
</html>`, inertia.WithVersion("v1"))

    if err != nil {
        log.Fatal(err)
    }

    mux := http.NewServeMux()
    mux.Handle("/", i.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        i.Render(w, r, "Home", httpx.Props{"message": "Hello from Go!"})
    })))

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Packages

| Package | Purpose |
|---------|---------|
| `inertia/` | Core engine -- constructors, `Render`, `Redirect`, shared props, context helpers |
| `httpx/` | Header constants (`X-Inertia-*`), shared types (`Props`, `ValidationErrors`) |
| `props/` | Prop types (`Always`, `Defer`, `Once`, `Merge`, `Optional`) and resolver |
| `middleware/` | HTTP middleware -- version checking, `Vary` header, redirect conversion |
| `response/` | Page object and HTML/JSON response rendering |
| `assert/` | `AssertableInertia` test helpers |

## SEO / Head Management

Server-rendered head content can be configured globally, per locale, and per request:

```go
i, err := inertia.New(templateHTML,
    inertia.WithHeadDefaults(),
    inertia.WithHeadFromFile("config/seo.yml"),
)

ctx := inertia.SetTitle(r.Context(), "Dashboard")
ctx = inertia.SetMeta(ctx, httpx.MetaTag{Name: "description", Content: "Ops overview"})
```

Precedence is:

- built-in defaults
- file config
- env overrides
- explicit Go-side overrides (`WithHead`, `SetHead`, `SetTitle`, `SetMeta`, `SetLinks`)

## Rendering Pages

```go
// JSON response for XHR visits, HTML for initial page loads.
i.Render(w, r, "Users/Index", httpx.Props{
    "users": users,
    "title": "User List",
})
```

## Shared Props

Props available on every page:

```go
i.ShareProp("auth", map[string]any{
    "user": currentUser,
})

i.ShareProps(httpx.Props{
    "app_name": "My App",
    "version":  "2.0",
})
```

## Prop Types

```go
import "github.com/oullin/inertia-go/core/props"

i.Render(w, r, "Dashboard", httpx.Props{
    // Always included, even in partial reloads.
    "flash": props.Always(flashMessages),

    // Excluded from initial load, fetched async by the client.
    "stats": props.Defer(func() any {
        return db.GetStats()
    }, "sidebar"),

    // Resolved once, client reuses on subsequent visits.
    "config": props.Once(appConfig),

    // Only included when explicitly requested in partial reloads.
    "debug": props.Optional(debugInfo),

    // Client merges new data with existing (infinite scroll).
    "posts": props.Merge(nextPage),

    // Deep merge for nested structures.
    "settings": props.DeepMerge(updatedSettings),
})
```

## Middleware

The middleware handles the Inertia protocol automatically:

- Sets `Vary: X-Inertia` for HTTP caching
- Checks `X-Inertia-Version` and returns `409` on mismatch
- Converts `302` to `303` for `PUT`/`PATCH`/`DELETE` redirects

```go
// Works with any router that accepts func(http.Handler) http.Handler
mux.Handle("/", i.Middleware(handler))

// chi
r.Use(i.Middleware)

// alice
chain := alice.New(i.Middleware)
```

## Context Helpers

Set per-request data from middleware or handlers:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := inertia.SetProp(r.Context(), "user", currentUser)
        ctx = inertia.SetValidationErrors(ctx, httpx.ValidationErrors{
            "email": "This email is already taken.",
        })
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

Available helpers:

| Function | Purpose |
|----------|---------|
| `SetProp(ctx, key, val)` | Add a per-request prop |
| `SetProps(ctx, props)` | Add multiple per-request props |
| `SetValidationErrors(ctx, errors)` | Set validation errors (added as `"errors"` prop) |
| `SetEncryptHistory(ctx)` | Flag response to encrypt browser history |
| `SetClearHistory(ctx)` | Flag response to clear encrypted history |
| `SetTemplateData(ctx, data)` | Extra data for the root HTML template |
| `SetTemplateDatum(ctx, key, val)` | Add a single template data value |
| `SetHead(ctx, head)` | Set per-request head tags and attributes |
| `SetTitle(ctx, title)` | Set the page title |
| `SetMeta(ctx, tags...)` | Add or override meta tags |
| `SetLinks(ctx, links...)` | Add or override link tags |

## Redirects

```go
i.Redirect(w, r, "/dashboard")
i.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)

i.Back(w, r)                           // Redirect to Referer (fallback: "/")

i.Location(w, r, "https://external.com") // 409 + X-Inertia-Location for Inertia requests
```

## Precognition

Precognitive requests use the `Precognition: true` and `Validate-Only` headers. The middleware marks the request, while `Render` and `HandlePrecognition` write the validation-only response:

```go
handler := middleware.Precognition()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if len(errors) > 0 {
        ctx := inertia.SetValidationErrors(r.Context(), errors)
        _ = i.Render(w, r.WithContext(ctx), "Users/Form")
        return
    }

    if handled, err := i.HandlePrecognition(w, r); handled {
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // perform mutation side effects here
}))
```

Use `HandlePrecognition` before mutating state in successful precognitive flows.

## Request Forgery Protection

The CSRF middleware now mirrors Laravel 13's two-layer approach:

- allow same-origin requests immediately via `Sec-Fetch-Site`
- optionally allow same-site requests
- fall back to `_token`, `X-CSRF-TOKEN`, or `X-XSRF-TOKEN`

`XSRF-TOKEN` cookies use Laravel-compatible encrypted values and accept URL-encoded `X-XSRF-TOKEN` headers.

## Options

```go
i, err := inertia.New(templateHTML,
    inertia.WithVersion("v1"),                  // Static version string
    inertia.WithVersionFromFile("manifest.json"), // Hash file for version
    inertia.WithContainerID("app"),             // Root div ID (default: "app")
    inertia.WithEncryptHistory(),               // Encrypt history by default
    inertia.WithJSONMarshaler(customMarshaler), // Swap JSON encoder
    inertia.WithLogger(logger),                 // Diagnostic logging
    inertia.WithTemplateFuncs(template.FuncMap{ // Extra template funcs
        "upper": strings.ToUpper,
    }),
)
```

Constructors:

```go
inertia.New(htmlString, opts...)           // Parse template string
inertia.NewFromFile("app.html", opts...)   // Parse template file
inertia.NewFromReader(reader, opts...)     // Parse from io.Reader
inertia.NewFromTemplate(tmpl, opts...)     // Use pre-parsed template
```

## Testing

```go
import "github.com/oullin/inertia-go/core/assert"

func TestUsersPage(t *testing.T) {
    req := httptest.NewRequest("GET", "/users", nil)
    req.RequestURI = "/users"

    result := assert.AssertFromHandler(t, i, usersHandler, req)
    result.AssertComponent(t, "Users/Index")
    result.AssertURL(t, "/users")
    result.AssertVersion(t, "v1")
    result.AssertHasProp(t, "users")
    result.AssertPropEquals(t, "title", "User List")
    result.AssertMissingProp(t, "secret")
}

// Or decode raw JSON:
result := assert.AssertFromBytes(t, responseBody)
result := assert.AssertFromReader(t, resp.Body)
```

## Root Template

The root HTML template receives two special variables:

- `{{ .inertia }}` -- A `<script type="application/json">` element with page data, followed by the container div
- `{{ .inertiaHead }}` -- Server-rendered head content assembled from defaults, locale config, and per-request overrides

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    {{ .inertiaHead }}
    <link rel="stylesheet" href="/app.css">
</head>
<body>
    {{ .inertia }}
    <script type="module" src="/app.js"></script>
</body>
</html>
```

## Demo App

The `demo/` directory contains a full working app with Go + Vue 3 + Vite:

```bash
# From the repository root
make demo
```

Then visit `http://localhost:8080`.

## Protocol Reference

This adapter implements the [Inertia.js protocol](https://inertiajs.com/docs/v3/core-concepts/the-protocol), including:

- Initial HTML visits with embedded page data
- XHR JSON responses for subsequent navigation
- Asset version checking with `409 Conflict`
- Partial reloads via `X-Inertia-Partial-*` headers
- Deferred, once, merge, always, and optional prop types
- Automatic `302` to `303` redirect conversion for `PUT`/`PATCH`/`DELETE`

## Licence

MIT
