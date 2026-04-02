package i18n_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/oullin/inertia-go/core/i18n"
)

const testConfig = `
default_locale: "en"
url_prefix: true

locales:
  en:
    name: "English"
    direction: "ltr"
    head:
      title: "My App"
      meta:
        - name: "description"
          content: "English description"
  es:
    name: "Español"
    direction: "ltr"
    head:
      title: "Mi App"
      meta:
        - name: "description"
          content: "Descripción en español"
  ar:
    name: "العربية"
    direction: "rtl"
    head:
      title: "تطبيقي"
`

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.yml")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestLoadConfig(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, err := i18n.LoadConfig(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.DefaultLocale != "en" {
		t.Errorf("default locale = %q, want %q", cfg.DefaultLocale, "en")
	}

	if len(cfg.Locales) != 3 {
		t.Errorf("locales count = %d, want 3", len(cfg.Locales))
	}

	en := cfg.Lookup("en")

	if en == nil {
		t.Fatal("en locale not found")
	}

	if en.Code != "en" {
		t.Errorf("en.Code = %q", en.Code)
	}

	if en.Name != "English" {
		t.Errorf("en.Name = %q", en.Name)
	}

	if en.Head.Title != "My App" {
		t.Errorf("en.Head.Title = %q", en.Head.Title)
	}

	ar := cfg.Lookup("ar")

	if ar == nil {
		t.Fatal("ar locale not found")
	}

	if ar.Direction != "rtl" {
		t.Errorf("ar.Direction = %q", ar.Direction)
	}
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_I18N_DEFAULT_LOCALE", "es")

	path := writeConfigFile(t, testConfig)
	cfg, err := i18n.LoadConfig(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.DefaultLocale != "es" {
		t.Errorf("default locale = %q, want %q", cfg.DefaultLocale, "es")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := i18n.LoadConfig("/nonexistent/i18n.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDefault(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	d := cfg.Default()

	if d == nil || d.Code != "en" {
		t.Error("default locale should be en")
	}
}

func TestCodes(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	codes := cfg.Codes()

	if len(codes) != 3 {
		t.Errorf("codes count = %d, want 3", len(codes))
	}
}

func TestMiddleware_DetectsPrefix(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedPath string

	handler := cfg.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/es/dashboard", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if capturedPath != "/dashboard" {
		t.Errorf("path = %q, want %q", capturedPath, "/dashboard")
	}
}

func TestMiddleware_DefaultLocale(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedPath string

	handler := cfg.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	// No prefix detected, path should remain unchanged.
	if capturedPath != "/dashboard" {
		t.Errorf("path = %q, want %q", capturedPath, "/dashboard")
	}
}

func TestMiddleware_StripsPrefixFromPath(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedURI string

	handler := cfg.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURI = r.RequestURI
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/ar/settings?tab=profile", nil)
	r.RequestURI = "/ar/settings?tab=profile"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if capturedURI != "/settings?tab=profile" {
		t.Errorf("requestURI = %q, want %q", capturedURI, "/settings?tab=profile")
	}
}

func TestMiddleware_RootWithPrefix(t *testing.T) {
	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedPath string

	handler := cfg.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/es", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if capturedPath != "/" {
		t.Errorf("path = %q, want %q", capturedPath, "/")
	}
}
