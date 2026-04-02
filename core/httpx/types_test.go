package httpx_test

import (
	"context"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
)

func TestCSRFToken_ContextRoundTrip(t *testing.T) {
	ctx := httpx.SetCSRFToken(context.Background(), "test-token-abc")
	got := httpx.CSRFTokenFromContext(ctx)

	if got != "test-token-abc" {
		t.Errorf("CSRFTokenFromContext() = %q, want %q", got, "test-token-abc")
	}
}

func TestCSRFTokenFromContext_Missing(t *testing.T) {
	got := httpx.CSRFTokenFromContext(context.Background())

	if got != "" {
		t.Errorf("CSRFTokenFromContext() = %q, want empty string", got)
	}
}

func TestLocale_ContextRoundTrip(t *testing.T) {
	locale := &httpx.Locale{
		Code:      "es",
		Name:      "Español",
		Direction: "ltr",
	}

	ctx := httpx.SetLocale(context.Background(), locale)
	got := httpx.LocaleFromContext(ctx)

	if got == nil {
		t.Fatal("LocaleFromContext() = nil, want locale")
	}

	if got.Code != "es" {
		t.Errorf("Code = %q, want %q", got.Code, "es")
	}

	if got.Name != "Español" {
		t.Errorf("Name = %q, want %q", got.Name, "Español")
	}
}

func TestLocaleFromContext_Missing(t *testing.T) {
	got := httpx.LocaleFromContext(context.Background())

	if got != nil {
		t.Errorf("LocaleFromContext() = %v, want nil", got)
	}
}
