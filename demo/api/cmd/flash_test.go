package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSetFlashAndConsumeFlash(t *testing.T) {
	w := httptest.NewRecorder()
	setFlash(w, flashPayload{Kind: "success", Title: "Done", Message: "It worked"})

	flashCookie := findCookie(t, w, flashCookieName)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.AddCookie(flashCookie)
	w2 := httptest.NewRecorder()

	flash := consumeFlash(w2, req)

	if flash == nil {
		t.Fatal("flash should be present")
	}

	if flash["kind"] != "success" {
		t.Fatalf("kind = %v, want %q", flash["kind"], "success")
	}

	if flash["title"] != "Done" {
		t.Fatalf("title = %v, want %q", flash["title"], "Done")
	}

	cleared := findCookie(t, w2, flashCookieName)

	if cleared.MaxAge != -1 {
		t.Fatalf("MaxAge = %d, want -1", cleared.MaxAge)
	}
}

func TestConsumeFlashReturnsNilForMissingOrInvalidCookie(t *testing.T) {
	tests := []struct {
		name   string
		cookie *http.Cookie
	}{
		{name: "missing"},
		{name: "invalid json", cookie: &http.Cookie{Name: flashCookieName, Value: url.QueryEscape("not-json")}},
		{name: "invalid escape", cookie: &http.Cookie{Name: flashCookieName, Value: "%zz"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/login", nil)

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			w := httptest.NewRecorder()

			if got := consumeFlash(w, req); got != nil {
				t.Fatalf("consumeFlash() = %#v, want nil", got)
			}
		})
	}
}
