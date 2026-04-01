package inertia_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ihttp "github.com/oullin/inertia-go/core/http"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/props"
	"github.com/oullin/inertia-go/core/response"
)

const testTemplate = `<!DOCTYPE html>
<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}</body>
</html>`

func newTestInertia(t *testing.T) *inertia.Inertia {
	t.Helper()
	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	return i
}

func TestNew(t *testing.T) {
	i := newTestInertia(t)

	if i.Version() != "v1" {
		t.Errorf("Version() = %q, want %q", i.Version(), "v1")
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := inertia.New("{{ .invalid }", inertia.WithVersion("v1"))

	if err == nil {
		t.Error("expected error for invalid template")
	}
}

func TestRender_JSONResponse(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.Header.Set(ihttp.HeaderInertia, "true")
	r.Header.Set(ihttp.HeaderVersion, "v1")
	r.RequestURI = "/users"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Users/Index", ihttp.Props{
		"users": []string{"alice", "bob"},
	})

	if err != nil {
		t.Fatal(err)
	}

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}

	if resp.Header.Get(ihttp.HeaderInertia) != "true" {
		t.Error("missing X-Inertia header")
	}

	var page response.Page

	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}

	if page.Component != "Users/Index" {
		t.Errorf("component = %q", page.Component)
	}

	if page.URL != "/users" {
		t.Errorf("url = %q", page.URL)
	}

	if page.Version != "v1" {
		t.Errorf("version = %q", page.Version)
	}
}

func TestRender_HTMLResponse(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.RequestURI = "/users"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Users/Index", ihttp.Props{
		"users": []string{"alice"},
	})

	if err != nil {
		t.Fatal(err)
	}

	resp := w.Result()

	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}

	body := w.Body.String()

	if !contains(body, `id="app"`) {
		t.Error("HTML response missing app container")
	}

	if !contains(body, `data-page=`) {
		t.Error("HTML response missing data-page attribute")
	}
}

func TestRender_NoProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(ihttp.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Dashboard")

	if err != nil {
		t.Fatal(err)
	}

	var page response.Page

	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}

	if page.Component != "Dashboard" {
		t.Errorf("component = %q", page.Component)
	}
}

func TestSharedProps(t *testing.T) {
	i := newTestInertia(t)
	i.ShareProp("app_name", "TestApp")
	i.ShareProps(ihttp.Props{"version": "1.0"})

	shared := i.SharedProps()

	if shared["app_name"] != "TestApp" {
		t.Errorf("app_name = %v", shared["app_name"])
	}

	if shared["version"] != "1.0" {
		t.Errorf("version = %v", shared["version"])
	}
}

func TestSharedProps_MergedInRender(t *testing.T) {
	i := newTestInertia(t)
	i.ShareProp("app_name", "TestApp")

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(ihttp.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Page", ihttp.Props{"title": "Hello"})

	if err != nil {
		t.Fatal(err)
	}

	var page response.Page

	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}

	if page.Props["app_name"] != "TestApp" {
		t.Errorf("shared prop app_name = %v", page.Props["app_name"])
	}

	if page.Props["title"] != "Hello" {
		t.Errorf("prop title = %v", page.Props["title"])
	}
}

func TestContextProps_MergedInRender(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(ihttp.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "user", "alice")
	ctx = inertia.SetValidationErrors(ctx, ihttp.ValidationErrors{"email": "required"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	err := i.Render(w, r, "Page")

	if err != nil {
		t.Fatal(err)
	}

	var page response.Page

	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}

	if page.Props["user"] != "alice" {
		t.Errorf("context prop user = %v", page.Props["user"])
	}

	if page.Props["errors"] == nil {
		t.Error("validation errors not included in props")
	}
}

func TestRender_DeferredProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(ihttp.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Dashboard", ihttp.Props{
		"title": "Dashboard",
		"stats": props.Defer(func() any { return "expensive" }, "sidebar"),
	})

	if err != nil {
		t.Fatal(err)
	}

	var page response.Page

	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}

	if _, ok := page.Props["stats"]; ok {
		t.Error("deferred prop should not be in initial response")
	}

	if page.DeferredProps["sidebar"] == nil {
		t.Error("expected deferred props metadata")
	}
}

func TestRedirect(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	i.Redirect(w, r, "/dashboard")

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusFound)
	}
}

func TestBack(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Referer", "/previous")
	w := httptest.NewRecorder()

	i.Back(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if loc := w.Header().Get("Location"); loc != "/previous" {
		t.Errorf("Location = %q, want %q", loc, "/previous")
	}
}

func TestBack_FallbackToRoot(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	i.Back(w, r)

	if loc := w.Header().Get("Location"); loc != "/" {
		t.Errorf("Location = %q, want %q", loc, "/")
	}
}

func TestLocation_InertiaRequest(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(ihttp.HeaderInertia, "true")
	w := httptest.NewRecorder()

	i.Location(w, r, "https://external.com")

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}

	if loc := w.Header().Get(ihttp.HeaderLocation); loc != "https://external.com" {
		t.Errorf("X-Inertia-Location = %q", loc)
	}
}

func TestLocation_NonInertiaRequest(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	i.Location(w, r, "/dashboard")

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusFound)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
