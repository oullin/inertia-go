package httpx_test

import (
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
)

func TestHeaderConstants(t *testing.T) {
	t.Parallel()

	// Verify header names match the Inertia.js protocol spec.
	if httpx.HeaderInertia != "X-Inertia" {
		t.Errorf("HeaderInertia = %q, want %q", httpx.HeaderInertia, "X-Inertia")
	}

	if httpx.HeaderVersion != "X-Inertia-Version" {
		t.Errorf("HeaderVersion = %q, want %q", httpx.HeaderVersion, "X-Inertia-Version")
	}

	if httpx.HeaderPartialComponent != "X-Inertia-Partial-Component" {
		t.Errorf("HeaderPartialComponent = %q, want %q", httpx.HeaderPartialComponent, "X-Inertia-Partial-Component")
	}

	if httpx.HeaderPartialData != "X-Inertia-Partial-Data" {
		t.Errorf("HeaderPartialData = %q, want %q", httpx.HeaderPartialData, "X-Inertia-Partial-Data")
	}

	if httpx.HeaderPartialExcept != "X-Inertia-Partial-Except" {
		t.Errorf("HeaderPartialExcept = %q, want %q", httpx.HeaderPartialExcept, "X-Inertia-Partial-Except")
	}

	if httpx.HeaderLocation != "X-Inertia-Location" {
		t.Errorf("HeaderLocation = %q, want %q", httpx.HeaderLocation, "X-Inertia-Location")
	}

	if httpx.HeaderInfiniteScroll != "X-Inertia-Infinite-Scroll-Merge-Intent" {
		t.Errorf("HeaderInfiniteScroll = %q, want %q", httpx.HeaderInfiniteScroll, "X-Inertia-Infinite-Scroll-Merge-Intent")
	}

	if httpx.HeaderExceptOnceProps != "X-Inertia-Except-Once-Props" {
		t.Errorf("HeaderExceptOnceProps = %q, want %q", httpx.HeaderExceptOnceProps, "X-Inertia-Except-Once-Props")
	}

	if httpx.HeaderReset != "X-Inertia-Reset" {
		t.Errorf("HeaderReset = %q, want %q", httpx.HeaderReset, "X-Inertia-Reset")
	}

	if httpx.HeaderErrorBag != "X-Inertia-Error-Bag" {
		t.Errorf("HeaderErrorBag = %q, want %q", httpx.HeaderErrorBag, "X-Inertia-Error-Bag")
	}

	if httpx.HeaderRedirect != "X-Inertia-Redirect" {
		t.Errorf("HeaderRedirect = %q, want %q", httpx.HeaderRedirect, "X-Inertia-Redirect")
	}
}
