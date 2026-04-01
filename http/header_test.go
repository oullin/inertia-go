package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ihttp "github.com/oullin/inertia-go/http"
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
				r.Header.Set(ihttp.HeaderInertia, tt.header)
			}

			if got := ihttp.IsInertiaRequest(r); got != tt.want {
				t.Errorf("IsInertiaRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeaderConstants(t *testing.T) {
	// Verify header names match the Inertia.js protocol spec.
	if ihttp.HeaderInertia != "X-Inertia" {
		t.Errorf("HeaderInertia = %q", ihttp.HeaderInertia)
	}

	if ihttp.HeaderVersion != "X-Inertia-Version" {
		t.Errorf("HeaderVersion = %q", ihttp.HeaderVersion)
	}

	if ihttp.HeaderPartialComponent != "X-Inertia-Partial-Component" {
		t.Errorf("HeaderPartialComponent = %q", ihttp.HeaderPartialComponent)
	}

	if ihttp.HeaderPartialData != "X-Inertia-Partial-Data" {
		t.Errorf("HeaderPartialData = %q", ihttp.HeaderPartialData)
	}

	if ihttp.HeaderPartialExcept != "X-Inertia-Partial-Except" {
		t.Errorf("HeaderPartialExcept = %q", ihttp.HeaderPartialExcept)
	}

	if ihttp.HeaderLocation != "X-Inertia-Location" {
		t.Errorf("HeaderLocation = %q", ihttp.HeaderLocation)
	}
}
