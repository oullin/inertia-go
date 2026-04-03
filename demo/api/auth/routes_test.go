package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/assert"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/seed"
)

func TestLoginHandlerRendersPage(t *testing.T) {
	_, handler := newAuthTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set(httpx.HeaderInertia, "true")
	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())
	page.AssertComponent(t, "Auth/Login")
}

func TestLoginHandlerCreatesSession(t *testing.T) {
	_, handler := newAuthTestHandler(t)

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"password"},
		"remember": {"true"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Fatalf("location = %q, want %q", got, "/dashboard")
	}

	cookie := findCookie(t, w, SessionCookieName)

	if cookie.Value != "1" {
		t.Fatalf("session cookie value = %q, want %q", cookie.Value, "1")
	}
}

func TestLoginHandlerRejectsInvalidPassword(t *testing.T) {
	_, handler := newAuthTestHandler(t)

	body := strings.NewReader(url.Values{
		"email":    {"test@example.com"},
		"password": {"wrong-password"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")
	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

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

func TestLogoutHandlerClearsSession(t *testing.T) {
	_, handler := newAuthTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "1"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/login" {
		t.Fatalf("location = %q, want %q", got, "/login")
	}

	if findCookie(t, w, SessionCookieName).MaxAge != -1 {
		t.Fatalf("logout should clear the session cookie")
	}
}

func TestWithCurrentUserLoadsUserFromCookie(t *testing.T) {
	app, _ := newAuthTestHandler(t)

	handler := app.WithCurrentUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.CurrentUser(r)

		if user == nil {
			t.Fatal("current user should be present")
		}

		if user.Email != "test@example.com" {
			t.Fatalf("email = %q, want %q", user.Email, "test@example.com")
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "1"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func newAuthTestHandler(t *testing.T) (App, http.Handler) {
	t.Helper()

	testInertia, err := inertia.New(testTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatal(err)
	}

	testDB, err := database.Open(":memory:")

	if err != nil {
		t.Fatal(err)
	}

	if err := seed.Run(testDB); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testDB.Close()
	})

	app := New(Deps{
		DB: testDB,
		Render: func(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
			t.Helper()

			if err := testInertia.Render(w, r, component, pageProps); err != nil {
				t.Fatalf("render failed: %v", err)
			}
		},
		RenderWithContext: func(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
			t.Helper()

			if err := testInertia.Render(w, r, component, pageProps); err != nil {
				t.Fatalf("render with context failed: %v", err)
			}
		},
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			testInertia.Redirect(w, r, url)
		},
		RouteURL: authTestRouteURL,
		SetFlash: func(http.ResponseWriter, Flash) {},
	})

	mux := http.NewServeMux()
	app.RegisterRoutes(mux)

	return app, app.WithCurrentUser(mux)
}

func authTestRouteURL(name string, params map[string]string) string {
	pattern := map[string]string{
		"login":     "/login",
		"logout":    "/logout",
		"dashboard": "/dashboard",
	}[name]

	if pattern == "" {
		return "/"
	}

	for key, value := range params {
		pattern = strings.ReplaceAll(pattern, "{"+key+"}", value)
	}

	return pattern
}
