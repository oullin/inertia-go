package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oullin/inertia-go/core/config"
)

func TestDefaultCSRF(t *testing.T) {
	cfg := config.DefaultCSRF()

	if cfg.Secret != "" {
		t.Errorf("Secret = %q, want empty", cfg.Secret)
	}

	if cfg.CookieName != "_csrf_token" {
		t.Errorf("CookieName = %q, want %q", cfg.CookieName, "_csrf_token")
	}

	if cfg.HeaderName != "X-CSRF-TOKEN" {
		t.Errorf("HeaderName = %q, want %q", cfg.HeaderName, "X-CSRF-TOKEN")
	}

	if cfg.Secure {
		t.Error("Secure should be false by default")
	}

	if cfg.SameSite != "lax" {
		t.Errorf("SameSite = %q, want %q", cfg.SameSite, "lax")
	}
}

func TestLoadCSRF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
secret: "my-secret"
secure: true
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadCSRF(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Secret != "my-secret" {
		t.Errorf("Secret = %q, want %q", cfg.Secret, "my-secret")
	}

	if !cfg.Secure {
		t.Error("Secure should be true (from file)")
	}

	// Defaults should be applied for fields not in the file.
	if cfg.CookieName != "_csrf_token" {
		t.Errorf("CookieName = %q, want %q (default)", cfg.CookieName, "_csrf_token")
	}

	if cfg.HeaderName != "X-CSRF-TOKEN" {
		t.Errorf("HeaderName = %q, want %q (default)", cfg.HeaderName, "X-CSRF-TOKEN")
	}

	if cfg.SameSite != "lax" {
		t.Errorf("SameSite = %q, want %q (default)", cfg.SameSite, "lax")
	}
}

func TestLoadCSRF_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_CSRF_SECRET", "env-secret")

	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
secret: "file-secret"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadCSRF(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Secret != "env-secret" {
		t.Errorf("Secret = %q, want %q (env override)", cfg.Secret, "env-secret")
	}
}

func TestLoadCSRF_FileNotFound(t *testing.T) {
	_, err := config.LoadCSRF("/nonexistent/csrf.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}
