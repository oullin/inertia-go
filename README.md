# inertia-go

A Go server-side adapter for [Inertia.js](https://inertiajs.com). Build modern single-page apps using Vue, React, or Svelte with a Go backend -- no client-side routing or API layer required.

## Install

```bash
go get github.com/oullin/inertia-go/core
```

**Requires Go 1.25+**

## Quick Start

```go
package main

import (
    "log"
    "net/http"

    ihttp "github.com/oullin/inertia-go/core/http"
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
        i.Render(w, r, "Home", ihttp.Props{"message": "Hello from Go!"})
    })))

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Packages

| Package | Purpose |
|---------|---------|
| `inertia/` | Core engine -- constructors, `Render`, `Redirect`, shared props, context helpers |
| `http/` | Header constants (`X-Inertia-*`), shared types (`Props`, `ValidationErrors`) |
| `props/` | Prop types (`Always`, `Defer`, `Once`, `Merge`, `Optional`) and resolver |
| `middleware/` | HTTP middleware -- version checking, `Vary` header, redirect conversion |
| `response/` | Page object and HTML/JSON response rendering |
| `testing/` | `AssertableInertia` test helpers |

## Rendering Pages

```go
// JSON response for XHR visits, HTML for initial page loads.
i.Render(w, r, "Users/Index", ihttp.Props{
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

i.ShareProps(ihttp.Props{
    "app_name": "My App",
    "version":  "2.0",
})
```

## Prop Types

```go
import "github.com/oullin/inertia-go/core/props"

i.Render(w, r, "Dashboard", ihttp.Props{
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
        ctx = inertia.SetValidationErrors(ctx, ihttp.ValidationErrors{
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

## Redirects

```go
i.Redirect(w, r, "/dashboard")
i.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)

i.Back(w, r)                           // Redirect to Referer (fallback: "/")

i.Location(w, r, "https://external.com") // 409 + X-Inertia-Location for Inertia requests
```

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
import itesting "github.com/oullin/inertia-go/core/testing"

func TestUsersPage(t *testing.T) {
    req := httptest.NewRequest("GET", "/users", nil)
    req.RequestURI = "/users"

    result := itesting.AssertFromHandler(t, i, usersHandler, req)
    result.AssertComponent(t, "Users/Index")
    result.AssertURL(t, "/users")
    result.AssertVersion(t, "v1")
    result.AssertHasProp(t, "users")
    result.AssertPropEquals(t, "title", "User List")
    result.AssertMissingProp(t, "secret")
}

// Or decode raw JSON:
result := itesting.AssertFromBytes(t, responseBody)
result := itesting.AssertFromReader(t, resp.Body)
```

## Root Template

The root HTML template receives two special variables:

- `{{ .inertia }}` -- The container div with `data-page` attribute
- `{{ .inertiaHead }}` -- SSR head content (empty until SSR is configured)

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

## Example App

The `example/` directory contains a full working app with Go + Vue 3 + Vite:

```bash
# From the repository root
make example
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
