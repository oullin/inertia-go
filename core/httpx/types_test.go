package httpx_test

import (
	"context"
	"net/http/httptest"
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

func TestPrecognition_ContextRoundTrip(t *testing.T) {
	ctx := httpx.SetPrecognition(context.Background())

	if !httpx.IsPrecognition(ctx) {
		t.Error("IsPrecognition() = false, want true")
	}
}

func TestIsPrecognition_Missing(t *testing.T) {
	if httpx.IsPrecognition(context.Background()) {
		t.Error("IsPrecognition() = true, want false")
	}
}

func TestIsPrecognitionRequest(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("Precognition", "true")

	if !httpx.IsPrecognitionRequest(r) {
		t.Error("IsPrecognitionRequest() = false, want true")
	}
}

func TestIsPrecognitionRequest_Missing(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)

	if httpx.IsPrecognitionRequest(r) {
		t.Error("IsPrecognitionRequest() = true, want false")
	}
}

func TestValidateOnly(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("Validate-Only", "name,email,phone")

	fields := httpx.ValidateOnly(r)

	if len(fields) != 3 {
		t.Fatalf("len = %d, want 3", len(fields))
	}

	expected := []string{"name", "email", "phone"}

	for i, f := range fields {
		if f != expected[i] {
			t.Errorf("fields[%d] = %q, want %q", i, f, expected[i])
		}
	}
}

func TestValidateOnly_WithSpaces(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("Validate-Only", " name , email ")

	fields := httpx.ValidateOnly(r)

	if len(fields) != 2 {
		t.Fatalf("len = %d, want 2", len(fields))
	}

	if fields[0] != "name" || fields[1] != "email" {
		t.Errorf("fields = %v, want [name email]", fields)
	}
}

func TestValidateOnly_Missing(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)

	if fields := httpx.ValidateOnly(r); fields != nil {
		t.Errorf("ValidateOnly() = %v, want nil", fields)
	}
}

func TestValidateOnly_Empty(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("Validate-Only", "")

	if fields := httpx.ValidateOnly(r); fields != nil {
		t.Errorf("ValidateOnly() = %v, want nil", fields)
	}
}
