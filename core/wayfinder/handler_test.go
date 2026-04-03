package wayfinder

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	reg := New()
	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	handler := Handler(reg)

	req := httptest.NewRequest(http.MethodGet, "/__wayfinder", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")

	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var routes []Route

	if err := json.Unmarshal(rec.Body.Bytes(), &routes); err != nil {
		t.Fatal(err)
	}

	if len(routes) != 2 {
		t.Errorf("expected 2 routes, got %d", len(routes))
	}
}
