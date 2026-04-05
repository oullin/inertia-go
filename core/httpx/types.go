package httpx

import "errors"

// ErrNotFound indicates the requested resource could not be found.

// Props holds page properties passed to the frontend component.

// TemplateData holds additional data passed to the root HTML template
// during initial (non-XHR) page visits.

// ValidationErrors holds field-level validation errors to be shared
// with the frontend via the "errors" prop.

// JSONMarshaler abstracts JSON encoding/decoding so callers can swap
// in a faster implementation (e.g. github.com/goccy/go-json) without
// changing any Inertia code.

// Logger is a minimal logging interface compatible with the standard
// library's log.Logger and most structured logging packages.

var ErrNotFound = errors.New("not found")

type Props map[string]any

type TemplateData map[string]any

type ValidationErrors map[string]any

type JSONMarshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type Logger interface {
	Printf(format string, v ...any)
}
