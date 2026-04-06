# Core Testing

## Baseline Rules

- Use the Go standard library test stack: `testing`, `net/http/httptest`, and ordinary helpers.
- Do not introduce `testify`, `gomock`, or similar assertion or mocking frameworks.
- Prefer `t.Parallel()` in package tests unless shared mutable state makes it unsafe.
- Helpers should accept `testing.TB` and call `t.Helper()`.

## HTTP And Inertia Assertions

Use `httptest.ResponseRecorder` for HTTP-level assertions. When testing rendered Inertia responses, prefer [`assert.AssertFromHandler()`](../core/assert/assert.go) so the request runs through the middleware path instead of bypassing protocol behavior.

`AssertableInertia` is the main response helper and currently exposes:

- `AssertComponent`
- `AssertURL`
- `AssertVersion`
- `AssertHasProp`
- `AssertPropEquals`
- `AssertMissingProp`

That makes it the default tool when the test needs to validate serialized page payloads rather than internal structs.

## What To Cover

For changes in `core/`, tests should generally cover the behavior that agents are most likely to break:

- prop merge precedence
- partial reload filtering
- deferred, optional, once, merge, and scroll metadata
- middleware protocol behavior such as `Vary`, version conflicts, redirect normalization, and precognition handling
- head layering, template data, and context helpers
- validation error shape and field-name resolution

## Validation Notes

[`validation.Validate()`](../core/validation/validator.go) resolves field names in this order:

```text
json tag -> form tag -> Go field name
```

When a test depends on returned error keys, assert against the resolved field name rather than the struct field name unless the test intentionally covers fallback behavior.
