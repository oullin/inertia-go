package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/oullin/inertia-go/demo/api/internal/flash"
)

func TestSetFlashAndConsumeFlash(t *testing.T) {
	w := httptest.NewRecorder()
	setFlash(w, flash.Message{Kind: "success", Title: "Done", Message: "It worked"})

	flashCookie := findCookie(t, w, flashCookieName)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.AddCookie(flashCookie)
	w2 := httptest.NewRecorder()

	flash := consumeFlash(w2, req)

	if flash == nil {
		t.Fatal("flash should be present")
	}

	if flash.Kind != "success" {
		t.Fatalf("Kind = %q, want %q", flash.Kind, "success")
	}

	if flash.Title != "Done" {
		t.Fatalf("Title = %q, want %q", flash.Title, "Done")
	}

	if flash.Message != "It worked" {
		t.Fatalf("Message = %q, want %q", flash.Message, "It worked")
	}

	cleared := findCookie(t, w2, flashCookieName)

	if cleared.MaxAge != -1 {
		t.Fatalf("MaxAge = %d, want -1", cleared.MaxAge)
	}
}

func TestFlashMessageJSONUsesLowercaseKeys(t *testing.T) {
	data, err := json.Marshal(flash.Message{
		Kind:    "success",
		Title:   "Done",
		Message: "It worked",
	})

	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	want := `{"kind":"success","title":"Done","message":"It worked"}`

	if string(data) != want {
		t.Fatalf("Marshal() = %s, want %s", data, want)
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
