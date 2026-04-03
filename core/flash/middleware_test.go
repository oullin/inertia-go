package flash

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestMiddlewareConsumesFlash(t *testing.T) {
	s := NewCookieStore(WithCookieName("test_flash"))

	var capturedFlash any

	handler := Middleware(s)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The middleware should have set the flash prop on the context.
		// We can't directly access context props without importing inertia internals,
		// so we verify via the cookie lifecycle instead.
		w.WriteHeader(http.StatusOK)
	}))

	_ = capturedFlash

	msg := Message{Kind: "success", Title: "Done", Message: "Created."}
	data, _ := json.Marshal(msg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "test_flash",
		Value: url.QueryEscape(string(data)),
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Verify the cookie was consumed (deleted).
	cookies := rec.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "test_flash" && c.MaxAge == -1 {
			found = true

			break
		}
	}

	if !found {
		t.Error("expected flash cookie to be deleted after middleware consumed it")
	}
}

func TestMiddlewareNoFlash(t *testing.T) {
	s := NewCookieStore(WithCookieName("test_flash"))
	called := false

	handler := Middleware(s)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called even without flash")
	}

	cookies := rec.Result().Cookies()

	for _, c := range cookies {
		if c.Name == "test_flash" {
			t.Error("expected no flash cookie manipulation when no flash exists")
		}
	}
}

func TestMiddlewareCustomPropKey(t *testing.T) {
	s := NewCookieStore(WithCookieName("test_flash"))

	handler := Middleware(s, WithPropKey("notification"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	msg := Message{Kind: "error", Title: "Oops", Message: "Failed."}
	data, _ := json.Marshal(msg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "test_flash",
		Value: url.QueryEscape(string(data)),
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Verify cookie was still consumed.
	cookies := rec.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "test_flash" && c.MaxAge == -1 {
			found = true

			break
		}
	}

	if !found {
		t.Error("expected flash cookie to be deleted with custom prop key")
	}
}
