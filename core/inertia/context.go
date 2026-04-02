package inertia

import (
	"context"

	"github.com/oullin/inertia-go/core/httpx"
)

type contextKey struct{ name string }

var (
	ctxKeyProps            = &contextKey{"props"}
	ctxKeyTemplateData     = &contextKey{"templateData"}
	ctxKeyValidationErrors = &contextKey{"validationErrors"}
	ctxKeyEncryptHistory   = &contextKey{"encryptHistory"}
	ctxKeyClearHistory     = &contextKey{"clearHistory"}
	ctxKeyHead             = &contextKey{"head"}
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
func SetProps(ctx context.Context, props httpx.Props) context.Context {
	p := propsFromContext(ctx)

	for k, v := range props {
		p[k] = v
	}

	return context.WithValue(ctx, ctxKeyProps, p)
}

// SetValidationErrors stores validation errors in the request context.
// They are automatically added to the response props under the "errors" key.
func SetValidationErrors(ctx context.Context, errors httpx.ValidationErrors) context.Context {
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
func SetTemplateData(ctx context.Context, data httpx.TemplateData) context.Context {
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

func propsFromContext(ctx context.Context) httpx.Props {
	if p, ok := ctx.Value(ctxKeyProps).(httpx.Props); ok {
		return p
	}

	return make(httpx.Props)
}

func validationErrorsFromContext(ctx context.Context) httpx.ValidationErrors {
	if v, ok := ctx.Value(ctxKeyValidationErrors).(httpx.ValidationErrors); ok {
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

func templateDataFromContext(ctx context.Context) httpx.TemplateData {
	if d, ok := ctx.Value(ctxKeyTemplateData).(httpx.TemplateData); ok {
		return d
	}

	return make(httpx.TemplateData)
}

// SetHead stores head elements on the request context. These are rendered
// into {{ .inertiaHead }} on initial page loads. Per-request head elements
// are merged with (and override) default head elements set via WithHead.
func SetHead(ctx context.Context, head httpx.Head) context.Context {
	existing := headFromContext(ctx)
	merged := httpx.MergeHead(existing, head)

	return context.WithValue(ctx, ctxKeyHead, merged)
}

// SetTitle is a convenience helper that sets only the <title> element
// on the request context.
func SetTitle(ctx context.Context, title string) context.Context {
	return SetHead(ctx, httpx.Head{Title: title})
}

// SetLang is a convenience helper that sets only the lang attribute
// on the request context.
func SetLang(ctx context.Context, lang string) context.Context {
	return SetHead(ctx, httpx.Head{Lang: lang})
}

// SetMeta is a convenience helper that adds meta tags to the request context.
func SetMeta(ctx context.Context, tags ...httpx.MetaTag) context.Context {
	return SetHead(ctx, httpx.Head{Meta: tags})
}

// SetLinks is a convenience helper that adds link tags to the request context.
func SetLinks(ctx context.Context, links ...httpx.LinkTag) context.Context {
	return SetHead(ctx, httpx.Head{Links: links})
}

// SetCSRFToken stores a CSRF token in the request context. When present,
// Render automatically adds <meta name="csrf-token" content="TOKEN"> to
// the head on initial page loads. Delegates to httpx.SetCSRFToken.
func SetCSRFToken(ctx context.Context, token string) context.Context {
	return httpx.SetCSRFToken(ctx, token)
}

// SetPrecognition marks the request context as a precognition request.
// Delegates to httpx.SetPrecognition.
func SetPrecognition(ctx context.Context) context.Context {
	return httpx.SetPrecognition(ctx)
}

// SetLocale stores the resolved locale in the request context. This is
// typically called by the i18n middleware, not by application code.
// Delegates to httpx.SetLocale.
func SetLocale(ctx context.Context, locale *httpx.Locale) context.Context {
	return httpx.SetLocale(ctx, locale)
}

func headFromContext(ctx context.Context) httpx.Head {
	if h, ok := ctx.Value(ctxKeyHead).(httpx.Head); ok {
		return h
	}

	return httpx.Head{}
}
