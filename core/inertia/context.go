package inertia

import (
	"context"

	ihttp "github.com/oullin/inertia-go/core/http"
)

type contextKey struct{ name string }

var (
	ctxKeyProps            = &contextKey{"props"}
	ctxKeyTemplateData     = &contextKey{"templateData"}
	ctxKeyValidationErrors = &contextKey{"validationErrors"}
	ctxKeyEncryptHistory   = &contextKey{"encryptHistory"}
	ctxKeyClearHistory     = &contextKey{"clearHistory"}
)

// SetProp stores a single prop on the request context. Props set this
// way are merged into the response during Render, with higher priority
// than shared props but lower than props passed directly to Render.
func SetProp(ctx context.Context, key string, val any) context.Context {
	p := propsFromContext(ctx)
	p[key] = val

	return context.WithValue(ctx, ctxKeyProps, p)
}

// SetProps stores multiple props on the request context.
func SetProps(ctx context.Context, props ihttp.Props) context.Context {
	p := propsFromContext(ctx)

	for k, v := range props {
		p[k] = v
	}

	return context.WithValue(ctx, ctxKeyProps, p)
}

// SetValidationErrors stores validation errors on the request context.
// They are automatically added to the response props under the "errors" key.
func SetValidationErrors(ctx context.Context, errors ihttp.ValidationErrors) context.Context {
	return context.WithValue(ctx, ctxKeyValidationErrors, errors)
}

// SetEncryptHistory flags the response to encrypt the browser history state.
func SetEncryptHistory(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyEncryptHistory, true)
}

// SetClearHistory flags the response to clear any encrypted browser history.
func SetClearHistory(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyClearHistory, true)
}

// SetTemplateData stores additional data for the root HTML template
// used during initial (non-XHR) page visits.
func SetTemplateData(ctx context.Context, data ihttp.TemplateData) context.Context {
	existing := templateDataFromContext(ctx)

	for k, v := range data {
		existing[k] = v
	}

	return context.WithValue(ctx, ctxKeyTemplateData, existing)
}

// SetTemplateDatum stores a single template data value.
func SetTemplateDatum(ctx context.Context, key string, val any) context.Context {
	d := templateDataFromContext(ctx)
	d[key] = val

	return context.WithValue(ctx, ctxKeyTemplateData, d)
}

func propsFromContext(ctx context.Context) ihttp.Props {
	if p, ok := ctx.Value(ctxKeyProps).(ihttp.Props); ok {
		return p
	}

	return make(ihttp.Props)
}

func validationErrorsFromContext(ctx context.Context) ihttp.ValidationErrors {
	if v, ok := ctx.Value(ctxKeyValidationErrors).(ihttp.ValidationErrors); ok {
		return v
	}

	return nil
}

func encryptHistoryFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyEncryptHistory).(bool)

	return v
}

func clearHistoryFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyClearHistory).(bool)

	return v
}

func templateDataFromContext(ctx context.Context) ihttp.TemplateData {
	if d, ok := ctx.Value(ctxKeyTemplateData).(ihttp.TemplateData); ok {
		return d
	}

	return make(ihttp.TemplateData)
}
