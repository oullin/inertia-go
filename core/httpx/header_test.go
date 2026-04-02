package httpx_test

import (
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
)

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

	if httpx.HeaderInfiniteScroll != "X-Inertia-Infinite-Scroll-Merge-Intent" {
		t.Errorf("HeaderInfiniteScroll = %q", httpx.HeaderInfiniteScroll)
	}

	if httpx.HeaderExceptOnceProps != "X-Inertia-Except-Once-Props" {
		t.Errorf("HeaderExceptOnceProps = %q", httpx.HeaderExceptOnceProps)
	}

	if httpx.HeaderReset != "X-Inertia-Reset" {
		t.Errorf("HeaderReset = %q", httpx.HeaderReset)
	}

	if httpx.HeaderErrorBag != "X-Inertia-Error-Bag" {
		t.Errorf("HeaderErrorBag = %q", httpx.HeaderErrorBag)
	}

	if httpx.HeaderRedirect != "X-Inertia-Redirect" {
		t.Errorf("HeaderRedirect = %q", httpx.HeaderRedirect)
	}
}
