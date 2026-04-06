# Core Props And Rendering

## Prop Merge Order

The render pipeline merges props in this order:

```text
shared < context (SetProp / SetProps) < render-time < validation errors
```

Later sources win on key collisions. Validation errors are applied last so they reliably override any earlier `"errors"` value.

## Wrapper Chain Semantics

Prop wrappers compose by nesting. `walkPropChain()` in [`core/props/resolver.go`](../core/props/resolver.go) walks outward to inward and records the first instance of each trait, so the outermost wrapper wins when duplicate traits appear.

Examples:

- `Always(Defer(value))` is treated as always-included, because `Always` overrides defer filtering.
- `Merge(Defer(value))` preserves merge metadata while still following deferred inclusion rules.

## Resolver Pipeline

[`props.Resolve()`](../core/props/resolver.go) runs in two distinct stages:

1. Filter included props and collect metadata from request headers.
2. Evaluate only the props that survived filtering.

This split is intentional. Lazy values such as `func() any`, `Proper`, and `TryProper` must not execute during the filter phase.

## Partial Reload And Deferred Behavior

- Partial reloads activate only when the request is an Inertia request and `X-Inertia-Partial-Component` matches the current component.
- `X-Inertia-Partial-Data` narrows the included keys.
- `X-Inertia-Partial-Except` excludes keys even during a partial reload.
- Deferred props are omitted on initial visits and grouped for later fetches.
- Optional props are omitted on initial visits and only included when explicitly requested.
- Once props record metadata for the client and can be suppressed through the except-once header.
- Merge, deep-merge, and scroll metadata are recorded during filtering so the response payload carries the right client instructions.

## Render Flow

[`Inertia.Render()`](../core/inertia/inertia.go) follows this sequence:

1. Handle precognition early and short-circuit when appropriate.
2. Merge shared, context, render-time, and validation props.
3. Resolve props through the resolver pipeline.
4. Build a `response.Page`.
5. Write JSON for Inertia requests or render the root HTML template for initial visits.

For HTML responses, head metadata is layered in this order:

```text
global head < locale head < request head < auto-appended csrf meta
```

That layering matters when changing locale middleware, template rendering, or head helper behavior.
