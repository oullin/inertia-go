package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/middleware"
)

func TestMiddleware_SetsVaryHeader(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Vary"); got != httpx.HeaderInertia {
		t.Errorf("Vary = %q, want %q", got, httpx.HeaderInertia)
	}
}

func TestMiddleware_NonInertiaPassesThrough(t *testing.T) {
	called := false
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if !called {
		t.Error("next handler was not called for non-Inertia request")
	}

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestMiddleware_VersionMatch(t *testing.T) {
	called := false
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if !called {
		t.Error("next handler was not called when versions match")
	}
}

func TestMiddleware_VersionMismatch(t *testing.T) {
	called := false
	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	r.RequestURI = "/test"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if called {
		t.Error("next handler should not be called on version mismatch")
	}

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}

	if loc := w.Header().Get(httpx.HeaderLocation); loc != "/test" {
		t.Errorf("X-Inertia-Location = %q, want %q", loc, "/test")
	}
}

func TestMiddleware_VersionMismatchOnlyGET(t *testing.T) {
	called := false
	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if !called {
		t.Error("next handler should be called for POST even on version mismatch")
	}
}

func TestMiddleware_RedirectConversion_PUTTo303(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodPut, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d (303)", w.Code, http.StatusSeeOther)
	}
}

func TestMiddleware_RedirectConversion_PATCHTo303(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodPatch, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d (303)", w.Code, http.StatusSeeOther)
	}
}

func TestMiddleware_RedirectConversion_DELETETo303(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d (303)", w.Code, http.StatusSeeOther)
	}
}

func TestMiddleware_RedirectConversion_GETStays302(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d (302)", w.Code, http.StatusFound)
	}
}

func TestMiddleware_WriteWithoutExplicitWriteHeader(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write without calling WriteHeader first triggers implicit 200.
		w.Write([]byte("hello"))
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "hello" {
		t.Errorf("body = %q", w.Body.String())
	}
}

func TestMiddleware_DoubleWriteHeaderIgnored(t *testing.T) {
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(http.StatusInternalServerError) // Second call ignored.
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d (first call wins)", w.Code, http.StatusCreated)
	}
}
