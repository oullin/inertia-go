package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/config"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/middleware"
)

const testSecret = "test-secret-key-for-csrf"

func csrfMiddleware() func(http.Handler) http.Handler {
	return middleware.CSRF(config.CSRFConfig{
		Secret: testSecret,
	})
}

func TestCSRF_SetsCookieOnGET(t *testing.T) {
	handler := csrfMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	cookies := resp.Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "_csrf_token" {
			found = true

			if !c.HttpOnly {
				t.Error("cookie should be HttpOnly")
			}

			if !strings.Contains(c.Value, ".") {
				t.Error("cookie should contain signed token (token.sig)")
			}
		}
	}

	if !found {
		t.Error("missing CSRF cookie")
	}
}

func TestCSRF_SafeMethodsPassThrough(t *testing.T) {
	handler := csrfMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		r := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want %d", method, w.Code, http.StatusOK)
		}
	}
}

func TestCSRF_POSTWithoutToken_403(t *testing.T) {
	handler := csrfMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestCSRF_POSTWithValidToken(t *testing.T) {
	mw := csrfMiddleware()

	// Step 1: GET to obtain the cookie.
	getHandler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	getW := httptest.NewRecorder()
	getHandler.ServeHTTP(getW, getReq)

	cookies := getW.Result().Cookies()

	var tokenCookie *http.Cookie

	for _, c := range cookies {
		if c.Name == "_csrf_token" {
			tokenCookie = c

			break
		}
	}

	if tokenCookie == nil {
		t.Fatal("no CSRF cookie set")
	}

	// Extract raw token (before the ".").
	rawToken := strings.SplitN(tokenCookie.Value, ".", 2)[0]

	// Step 2: POST with cookie + header.
	postHandler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)
	postReq.AddCookie(tokenCookie)
	postReq.Header.Set("X-CSRF-TOKEN", rawToken)

	postW := httptest.NewRecorder()
	postHandler.ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", postW.Code, http.StatusOK)
	}
}

func TestCSRF_POSTWithWrongToken_403(t *testing.T) {
	mw := csrfMiddleware()

	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	getW := httptest.NewRecorder()
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(getW, getReq)

	cookies := getW.Result().Cookies()

	var tokenCookie *http.Cookie

	for _, c := range cookies {
		if c.Name == "_csrf_token" {
			tokenCookie = c

			break
		}
	}

	if tokenCookie == nil {
		t.Fatal("no CSRF cookie set")
	}

	postReq := httptest.NewRequest(http.MethodPost, "/submit", nil)
	postReq.AddCookie(tokenCookie)
	postReq.Header.Set("X-CSRF-TOKEN", "wrong-token")

	postW := httptest.NewRecorder()
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(postW, postReq)

	if postW.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", postW.Code, http.StatusForbidden)
	}
}

func TestCSRF_MutationMethods(t *testing.T) {
	mw := csrfMiddleware()

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		r := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()

		mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("%s: status = %d, want %d", method, w.Code, http.StatusForbidden)
		}
	}
}

func TestCSRFFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
secret: "file-secret"
cookie_name: "_my_csrf"
header_name: "X-MY-CSRF"
secure: false
same_site: "strict"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mw, err := middleware.CSRFFromFile(path)

	if err != nil {
		t.Fatal(err)
	}

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}

	// Check that it used the custom cookie name.
	found := false

	for _, c := range w.Result().Cookies() {
		if c.Name == "_my_csrf" {
			found = true
		}
	}

	if !found {
		t.Error("expected custom cookie name _my_csrf")
	}
}

func TestCSRFFromFile_FileNotFound(t *testing.T) {
	_, err := middleware.CSRFFromFile("/nonexistent/csrf.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCSRFFromFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	if err := os.WriteFile(path, []byte("secret: [\ninvalid"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := middleware.CSRFFromFile(path)

	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestCSRF_SameSiteStrict(t *testing.T) {
	mw := middleware.CSRF(config.CSRFConfig{
		Secret:   testSecret,
		SameSite: "strict",
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCSRF_SameSiteNone(t *testing.T) {
	mw := middleware.CSRF(config.CSRFConfig{
		Secret:   testSecret,
		SameSite: "none",
		Secure:   true,
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCSRFFromFile_AllEnvOverrides(t *testing.T) {
	t.Setenv("INERTIA_CSRF_SECRET", "env-secret")
	t.Setenv("INERTIA_CSRF_COOKIE_NAME", "_env_csrf")
	t.Setenv("INERTIA_CSRF_HEADER_NAME", "X-ENV-CSRF")
	t.Setenv("INERTIA_CSRF_SECURE", "true")
	t.Setenv("INERTIA_CSRF_SAME_SITE", "strict")

	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
secret: "file-secret"
cookie_name: "_file_csrf"
header_name: "X-FILE-CSRF"
secure: false
same_site: "lax"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mw, err := middleware.CSRFFromFile(path)

	if err != nil {
		t.Fatal(err)
	}

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	// Verify env-overridden cookie name.
	found := false

	for _, c := range w.Result().Cookies() {
		if c.Name == "_env_csrf" {
			found = true
		}
	}

	if !found {
		t.Error("expected env-overridden cookie name _env_csrf")
	}
}

func TestCSRF_POSTWithMalformedCookie_RegeneratesToken(t *testing.T) {
	mw := csrfMiddleware()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Send a POST with a cookie that has no "." separator (malformed).
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "_csrf_token", Value: "no-dot-separator"})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	// Should succeed on GET (new token generated and cookie set).
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	// Should have set a new cookie.
	cookies := w.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "_csrf_token" && strings.Contains(c.Value, ".") {
			found = true
		}
	}

	if !found {
		t.Error("expected new CSRF cookie to be set after malformed cookie")
	}
}

func TestCSRF_POSTWithTamperedCookie_RegeneratesToken(t *testing.T) {
	mw := csrfMiddleware()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Send a GET with a cookie that has a valid format but wrong signature.
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "_csrf_token", Value: "fake-token.bad-signature"})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	// Should have set a new cookie (regenerated).
	cookies := w.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "_csrf_token" && c.Value != "fake-token.bad-signature" {
			found = true
		}
	}

	if !found {
		t.Error("expected new CSRF cookie after tampered cookie")
	}
}

func TestCSRF_StoresTokenInContext(t *testing.T) {
	mw := csrfMiddleware()

	var ctxToken string

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxToken = httpx.CSRFTokenFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if ctxToken == "" {
		t.Error("expected CSRF token to be stored in context")
	}
}

func TestCSRFFromFile_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_CSRF_COOKIE_NAME", "_env_csrf")

	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
secret: "file-secret"
cookie_name: "_file_csrf"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mw, err := middleware.CSRFFromFile(path)

	if err != nil {
		t.Fatal(err)
	}

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	found := false

	for _, c := range w.Result().Cookies() {
		if c.Name == "_env_csrf" {
			found = true
		}
	}

	if !found {
		t.Error("expected env-overridden cookie name _env_csrf")
	}
}
