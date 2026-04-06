package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/cryptox"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func TestLoginValidationAndLogoutMethodGuard(t *testing.T) {
	t.Parallel()

	_, handler := newAuthTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(""))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = "/login"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodDelete, "/login", nil)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	app, _ := newAuthTestHandler(t)
	req = httptest.NewRequest(http.MethodGet, "/logout", nil)
	w = httptest.NewRecorder()

	app.logoutHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestRequireAuthGuestOnlyAndCurrentUserBranches(t *testing.T) {
	t.Parallel()

	app, _ := newAuthTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	app.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("RequireAuth should redirect when no user is present")
	})).ServeHTTP(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/login" {
		t.Fatalf("RequireAuth status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	encrypted, err := cryptox.Encrypt("1", testCryptoKey)

	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/login", nil)

	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: encrypted})

	w = httptest.NewRecorder()

	app.WithCurrentUser(app.GuestOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("GuestOnly should redirect when a user is present")
	}))).ServeHTTP(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/dashboard" {
		t.Fatalf("GuestOnly status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	if app.CurrentUser(httptest.NewRequest(http.MethodGet, "/", nil)) != nil {
		t.Fatal("CurrentUser() without context should return nil")
	}
}

func TestLoadCurrentUserAndPublicUser(t *testing.T) {
	t.Parallel()

	app, _ := newAuthTestHandler(t)

	if app.loadCurrentUser(httptest.NewRequest(http.MethodGet, "/", nil)) != nil {
		t.Fatal("loadCurrentUser() without cookie should return nil")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "not-encrypted"})

	if app.loadCurrentUser(req) != nil {
		t.Fatal("loadCurrentUser() with forged cookie should return nil")
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	encrypted, err := cryptox.Encrypt("abc", testCryptoKey)

	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: encrypted})

	if app.loadCurrentUser(req) != nil {
		t.Fatal("loadCurrentUser() with non-integer payload should return nil")
	}

	user := app.PublicUser(&database.User{
		ID:    1,
		Name:  "Ada Lovelace",
		Email: "ada@example.com",
	})

	payload, ok := user.(map[string]any)

	if !ok || payload["initials"] != "AL" {
		t.Fatalf("PublicUser() = %#v, want initials AL", user)
	}

	if app.PublicUser(nil) != nil {
		t.Fatal("PublicUser(nil) should return nil")
	}
}

func TestNewAppAndNewServiceErrors(t *testing.T) {
	t.Parallel()

	if _, err := newService(nil); err == nil {
		t.Fatal("newService(nil) error = nil, want error")
	}

	if _, err := NewApp(Container{}); err == nil {
		t.Fatal("NewApp(Container{}) error = nil, want validation error")
	}
}
