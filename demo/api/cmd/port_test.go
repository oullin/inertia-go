package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/assert"
	"github.com/oullin/inertia-go/core/config"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/middleware"
	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/seed"
)

func TestLoginHandlerRendersPage(t *testing.T) {
	testMux := newPortTestMux(t)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set(httpx.HeaderInertia, "true")
	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())
	page.AssertComponent(t, "Auth/Login")
	page.AssertHasProp(t, "auth")
	page.AssertHasProp(t, "routes")
}

func TestLoginHandlerCreatesSession(t *testing.T) {
	testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"password"},
		"remember": {"true"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Fatalf("location = %q, want %q", got, "/dashboard")
	}

	if findCookie(t, w, demoSessionCookieName).Value != "1" {
		t.Fatalf("expected session cookie for seeded user")
	}
}

func TestLoginHandlerRejectsInvalidPassword(t *testing.T) {
	testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"wrong-password"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.Header.Set(httpx.HeaderInertia, "true")
	req.RequestURI = "/login"
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())
	page.AssertComponent(t, "Auth/Login")

	errors, ok := page.Props["errors"].(map[string]any)

	if !ok {
		t.Fatal("errors prop not found or not a map")
	}

	if errors["email"] != "Use test@example.com and password to sign in." {
		t.Fatalf("errors[email] = %v", errors["email"])
	}
}

func TestDashboardRequiresSession(t *testing.T) {
	testMux := newPortTestMux(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/login" {
		t.Fatalf("location = %q, want %q", got, "/login")
	}
}

func TestDashboardRendersForAuthenticatedUser(t *testing.T) {
	testMux := newPortTestMux(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req.Header.Set(httpx.HeaderInertia, "true")
	req.RequestURI = "/dashboard"
	req.AddCookie(&http.Cookie{Name: demoSessionCookieName, Value: "1"})
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())
	page.AssertComponent(t, "Crm/Dashboard")
	page.AssertHasProp(t, "recentActivity")
	page.AssertHasProp(t, "auth")
}

func TestLegacyDemoRoutesReturnNotFound(t *testing.T) {
	testMux := newPortTestMux(t)

	tests := []struct {
		name string
		path string
	}{
		{name: "root", path: "/"},
		{name: "feature demo", path: "/features/forms/use-form"},
		{name: "dashboard navigation", path: "/dashboard/navigation"},
		{name: "dashboard data", path: "/dashboard/data"},
		{name: "dashboard state", path: "/dashboard/state"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.AddCookie(&http.Cookie{Name: demoSessionCookieName, Value: "1"})
			w := httptest.NewRecorder()

			testMux.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
			}
		})
	}
}

func TestStoreContactCreatesRecord(t *testing.T) {
	testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	before := database.CountContacts(db)

	body := strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Mina"},
		"last_name":       {"Cole"},
		"email":           {"mina@example.test"},
		"phone":           {"+1 555 0107"},
		"address":         {"55 New St"},
		"city":            {"Portland"},
		"region":          {"OR"},
		"country":         {"USA"},
		"postal_code":     {"97201"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/contacts", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	req.AddCookie(&http.Cookie{Name: demoSessionCookieName, Value: "1"})
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if database.CountContacts(db) != before+1 {
		t.Fatalf("contact count = %d, want %d", database.CountContacts(db), before+1)
	}
}

func TestStoreNoteAppendsActivity(t *testing.T) {
	testMux := newPortTestMux(t)
	csrfCookie, rawToken := issuePortCSRFCookie(t, testMux, "/login")

	before := database.CountNotes(db)
	body := strings.NewReader(url.Values{
		"body": {"Need legal review before the Friday call."},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/contacts/1/notes", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	req.AddCookie(&http.Cookie{Name: demoSessionCookieName, Value: "1"})
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if database.CountNotes(db) != before+1 {
		t.Fatalf("notes = %d, want %d", database.CountNotes(db), before+1)
	}
}

func newPortTestMux(t *testing.T) http.Handler {
	t.Helper()

	testInertia, err := inertia.New(testTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatal(err)
	}

	i = testInertia
	setupPortTestDB(t)
	t.Cleanup(func() { i = nil })

	mux := http.NewServeMux()
	registerAuthRoutes(mux)
	registerCRMRoutes(mux)

	cfg := config.DefaultI18n()
	cfg.URLPrefix = false

	return dashboardAppHandler(
		withDemoProps(mux),
		middleware.CSRF(config.CSRFConfig{}, []byte("0123456789abcdef0123456789abcdef")),
		cfg,
	)
}

func setupPortTestDB(t *testing.T) {
	t.Helper()

	testDB, err := database.Open(":memory:")

	if err != nil {
		t.Fatal(err)
	}

	if err := seed.Run(testDB); err != nil {
		t.Fatal(err)
	}

	db = testDB

	t.Cleanup(func() {
		testDB.Close()
		db = nil
	})
}

func issuePortCSRFCookie(t *testing.T, handler http.Handler, path string) (*http.Cookie, string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return findCookie(t, w, "XSRF-TOKEN"), findCSRFMetaToken(t, w.Body.String())
}
