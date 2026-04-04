package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/middleware"
)

func TestPrecognition_SetsVaryHeader(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Precognition" {
		t.Errorf("Vary = %q, want %q", vary, "Precognition")
	}
}

func TestPrecognition_NonPrecognitive_PassesThrough(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	var called bool

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		if httpx.IsPrecognition(r.Context()) {
			t.Error("context should not be marked as precognition")
		}

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if !called {
		t.Error("handler was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestPrecognition_SetsContextFlag(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	var isPrecognition bool

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPrecognition = httpx.IsPrecognition(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if !isPrecognition {
		t.Error("context should be marked as precognition")
	}

	if got := w.Header().Get(httpx.HeaderPrecognition); got != "true" {
		t.Errorf("Precognition header = %q, want %q", got, "true")
	}
}

func TestPrecognition_VaryHeader_AlwaysSet(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Precognition request should also have Vary header.
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Precognition" {
		t.Errorf("Vary = %q, want %q", vary, "Precognition")
	}
}
