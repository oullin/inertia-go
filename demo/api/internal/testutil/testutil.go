package testutil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const TestTemplate = `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`

func FindCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()

	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}

	t.Fatalf("cookie %q not found", name)

	return nil
}

func FindCSRFMetaToken(t *testing.T, body string) string {
	t.Helper()

	const prefix = `name="csrf-token" content="`

	start := strings.Index(body, prefix)

	if start == -1 {
		t.Fatal("csrf meta tag not found")
	}

	start += len(prefix)
	end := strings.Index(body[start:], `"`)

	if end == -1 {
		t.Fatal("csrf meta token not terminated")
	}

	return body[start : start+end]
}
