package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
)

func TestIsInertiaRequest(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{"with header", "true", true},
		{"without header", "", false},
		{"wrong value", "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.header != "" {
				r.Header.Set(httpx.HeaderInertia, tt.header)
			}

			if got := httpx.IsInertiaRequest(r); got != tt.want {
				t.Errorf("IsInertiaRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeaderConstants(t *testing.T) {
	// Verify header names match the Inertia.js protocol spec.
	if httpx.HeaderInertia != "X-Inertia" {
		t.Errorf("HeaderInertia = %q", httpx.HeaderInertia)
	}

	if httpx.HeaderVersion != "X-Inertia-Version" {
		t.Errorf("HeaderVersion = %q", httpx.HeaderVersion)
	}

	if httpx.HeaderPartialComponent != "X-Inertia-Partial-Component" {
		t.Errorf("HeaderPartialComponent = %q", httpx.HeaderPartialComponent)
	}

	if httpx.HeaderPartialData != "X-Inertia-Partial-Data" {
		t.Errorf("HeaderPartialData = %q", httpx.HeaderPartialData)
	}

	if httpx.HeaderPartialExcept != "X-Inertia-Partial-Except" {
		t.Errorf("HeaderPartialExcept = %q", httpx.HeaderPartialExcept)
	}

	if httpx.HeaderLocation != "X-Inertia-Location" {
		t.Errorf("HeaderLocation = %q", httpx.HeaderLocation)
	}
}
