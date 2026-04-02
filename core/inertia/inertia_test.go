package inertia_test

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/props"
	"github.com/oullin/inertia-go/core/response"
)

// --- Constructor variants ---

// --- Options ---

// Template funcs must be registered before parsing, so
// WithTemplateFuncs works for funcs available at execution time
// (e.g. sub-templates). Here we just verify the option is accepted.

// --- Context helpers ---

// --- Redirect with custom status ---

// --- Render error from prop resolution ---

// --- Middleware method ---

// --- StdJSONMarshaler ---

// --- Helpers ---

type failReader struct{}

type testJSONMarshaler struct {
	marshalCalled bool
}

type testLogger struct{}

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
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")
	r.RequestURI = "/users"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Users/Index", httpx.Props{
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

	if resp.Header.Get(httpx.HeaderInertia) != "true" {
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

	err := i.Render(w, r, "Users/Index", httpx.Props{
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

	if !contains(body, `<div id="app"></div>`) {
		t.Error("HTML response missing app container")
	}

	if !contains(body, `<script data-page="app" type="application/json">`) {
		t.Error("HTML response missing page data script element")
	}
}

func TestRender_NoProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
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
	i.ShareProps(httpx.Props{"version": "1.0"})

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
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Page", httpx.Props{"title": "Hello"})

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
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "user", "alice")
	ctx = inertia.SetValidationErrors(ctx, httpx.ValidationErrors{"email": "required"})
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
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Dashboard", httpx.Props{
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

func TestRender_OnceAndScrollProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/feed", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/feed"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Feed", httpx.Props{
		"release_notes": props.Once("Frozen snapshot"),
		"feed": props.Scroll(map[string]any{
			"data": []map[string]any{{"id": "evt_1"}},
		}, "feedPage", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	var page response.Page

	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}

	if page.OnceProps["release_notes"].Prop != "release_notes" {
		t.Errorf("once prop metadata = %+v", page.OnceProps["release_notes"])
	}

	if page.ScrollProps["feed"].PageName != "feedPage" {
		t.Errorf("scroll prop metadata = %+v", page.ScrollProps["feed"])
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
	r.Header.Set(httpx.HeaderInertia, "true")
	w := httptest.NewRecorder()

	i.Location(w, r, "https://external.com")

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}

	if loc := w.Header().Get(httpx.HeaderLocation); loc != "https://external.com" {
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

func TestNewFromFile(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/app.html"
	os.WriteFile(path, []byte(testTemplate), 0644)

	i, err := inertia.NewFromFile(path, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	if i.Version() != "v1" {
		t.Errorf("Version() = %q", i.Version())
	}
}

func TestNewFromFile_NotFound(t *testing.T) {
	_, err := inertia.NewFromFile("/nonexistent/path.html")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestNewFromReader(t *testing.T) {
	r := strings.NewReader(testTemplate)

	i, err := inertia.NewFromReader(r, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	if i.Version() != "v1" {
		t.Errorf("Version() = %q", i.Version())
	}
}

func TestNewFromReader_Error(t *testing.T) {
	_, err := inertia.NewFromReader(&failReader{})

	if err == nil {
		t.Error("expected error from failing reader")
	}
}

func TestNewFromTemplate(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(testTemplate))

	i, err := inertia.NewFromTemplate(tmpl, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	if i.Version() != "v1" {
		t.Errorf("Version() = %q", i.Version())
	}
}

func TestWithVersionFromFile(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/manifest.json"
	os.WriteFile(path, []byte(`{"app.js":"app.abc123.js"}`), 0644)

	i, err := inertia.New(testTemplate, inertia.WithVersionFromFile(path))

	if err != nil {
		t.Fatal(err)
	}

	if i.Version() == "" {
		t.Error("expected non-empty version from file hash")
	}
}

func TestWithVersionFromFile_NotFound(t *testing.T) {
	_, err := inertia.New(testTemplate, inertia.WithVersionFromFile("/nonexistent"))

	if err == nil {
		t.Error("expected error for missing version file")
	}
}

func TestWithContainerID(t *testing.T) {
	i, err := inertia.New(testTemplate, inertia.WithContainerID("root"))

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, `<div id="root"></div>`) {
		t.Error("container ID not applied")
	}
}

func TestWithEncryptHistory(t *testing.T) {
	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"), inertia.WithEncryptHistory())

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if !page.EncryptHistory {
		t.Error("encryptHistory should be true")
	}
}

func TestWithTemplateFuncs(t *testing.T) {

	_, err := inertia.New(testTemplate,
		inertia.WithVersion("v1"),
		inertia.WithTemplateFuncs(template.FuncMap{
			"upper": strings.ToUpper,
		}),
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestWithJSONMarshaler(t *testing.T) {
	m := &testJSONMarshaler{marshalCalled: false}

	i, err := inertia.New(testTemplate,
		inertia.WithVersion("v1"),
		inertia.WithJSONMarshaler(m),
	)

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	if !m.marshalCalled {
		t.Error("custom marshaler was not used")
	}
}

func TestWithLogger(t *testing.T) {
	l := &testLogger{}

	_, err := inertia.New(testTemplate,
		inertia.WithVersion("v1"),
		inertia.WithLogger(l),
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSetProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProps(r.Context(), httpx.Props{"a": "1", "b": "2"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["a"] != "1" {
		t.Errorf("prop a = %v", page.Props["a"])
	}

	if page.Props["b"] != "2" {
		t.Errorf("prop b = %v", page.Props["b"])
	}
}

func TestSetEncryptHistory(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetEncryptHistory(r.Context())
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if !page.EncryptHistory {
		t.Error("encryptHistory should be true from context")
	}
}

func TestSetClearHistory(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetClearHistory(r.Context())
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if !page.ClearHistory {
		t.Error("clearHistory should be true from context")
	}
}

func TestSetTemplateData(t *testing.T) {
	tmpl := `<!DOCTYPE html><html><body>{{ .customKey }}|{{ .inertia }}</body></html>`

	i, err := inertia.New(tmpl, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"

	ctx := inertia.SetTemplateData(r.Context(), httpx.TemplateData{"customKey": "custom-value"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, "custom-value") {
		t.Errorf("template data not rendered, body = %s", body)
	}
}

func TestSetTemplateDatum(t *testing.T) {
	tmpl := `<!DOCTYPE html><html><body>{{ .singleKey }}|{{ .inertia }}</body></html>`

	i, err := inertia.New(tmpl, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"

	ctx := inertia.SetTemplateDatum(r.Context(), "singleKey", "single-value")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, "single-value") {
		t.Errorf("template datum not rendered, body = %s", body)
	}
}

func TestRedirect_CustomStatus(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	i.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMovedPermanently)
	}
}

func TestRender_PropResolutionError(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	err := i.Render(w, r, "Page", httpx.Props{
		"bad": func() (any, error) { return nil, io.ErrUnexpectedEOF },
	})

	if err == nil {
		t.Error("expected error from failing prop resolution")
	}
}

func TestMiddleware_Method(t *testing.T) {
	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := i.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Vary"); got != httpx.HeaderInertia {
		t.Errorf("Vary = %q, want %q", got, httpx.HeaderInertia)
	}
}

func TestStdJSONMarshaler_Unmarshal(t *testing.T) {
	m := &inertia.StdJSONMarshaler{}

	var result map[string]string

	err := m.Unmarshal([]byte(`{"key":"value"}`), &result)

	if err != nil {
		t.Fatal(err)
	}

	if result["key"] != "value" {
		t.Errorf("key = %q, want %q", result["key"], "value")
	}
}

func TestStdJSONMarshaler_Unmarshal_Error(t *testing.T) {
	m := &inertia.StdJSONMarshaler{}

	var result map[string]string

	err := m.Unmarshal([]byte(`invalid json`), &result)

	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func (f *failReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (m *testJSONMarshaler) Marshal(v any) ([]byte, error) {
	m.marshalCalled = true

	return json.Marshal(v)
}

func (m *testJSONMarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (l *testLogger) Printf(format string, v ...any) {}

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

// --- Response Structure ---

func TestRender_URLPreservesFullRequestURI(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/users?page=2&sort=name", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/users?page=2&sort=name"
	w := httptest.NewRecorder()

	i.Render(w, r, "Users/Index")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.URL != "/users?page=2&sort=name" {
		t.Errorf("url = %q, want %q", page.URL, "/users?page=2&sort=name")
	}
}

func TestRender_URLWithTrailingSlash(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/users/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/users/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Users/Index")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.URL != "/users/" {
		t.Errorf("url = %q, want %q", page.URL, "/users/")
	}
}

func TestRender_URLWithTrailingSlashAndQueryParams(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/users/?page=1", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/users/?page=1"
	w := httptest.NewRecorder()

	i.Render(w, r, "Users/Index")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.URL != "/users/?page=1" {
		t.Errorf("url = %q, want %q", page.URL, "/users/?page=1")
	}
}

func TestRender_JSONOmitsEmptyOptionalFields(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	body := w.Body.String()

	for _, field := range []string{
		`"encryptHistory"`, `"clearHistory"`,
		`"mergeProps"`, `"deepMergeProps"`,
		`"deferredProps"`, `"scrollProps"`, `"onceProps"`,
	} {
		if contains(body, field) {
			t.Errorf("JSON should omit empty field %s", field)
		}
	}
}

// --- Shared Props ---

func TestSharedProps_ReturnsCopy(t *testing.T) {
	i := newTestInertia(t)
	i.ShareProp("key", "original")

	copy := i.SharedProps()
	copy["key"] = "mutated"

	if i.SharedProps()["key"] != "original" {
		t.Error("mutating SharedProps() return value should not affect internal state")
	}
}

func TestSharedProps_OverriddenByRenderProps(t *testing.T) {
	i := newTestInertia(t)
	i.ShareProp("title", "shared")

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{"title": "render-override"})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["title"] != "render-override" {
		t.Errorf("title = %v, want %q (render props should override shared)", page.Props["title"], "render-override")
	}
}

func TestSharedProps_OverriddenByContextProps(t *testing.T) {
	i := newTestInertia(t)
	i.ShareProp("title", "shared")

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "title", "context-override")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["title"] != "context-override" {
		t.Errorf("title = %v, want %q (context props should override shared)", page.Props["title"], "context-override")
	}
}

func TestRenderProps_OverrideContextProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "title", "context-val")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page", httpx.Props{"title": "render-val"})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["title"] != "render-val" {
		t.Errorf("title = %v, want %q (render props should override context)", page.Props["title"], "render-val")
	}
}

func TestSharedProps_MergedWithAllSources(t *testing.T) {
	i := newTestInertia(t)
	i.ShareProp("shared", "shared-val")

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "ctx", "ctx-val")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page", httpx.Props{"render": "render-val"})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["shared"] != "shared-val" {
		t.Error("shared prop missing")
	}

	if page.Props["ctx"] != "ctx-val" {
		t.Error("context prop missing")
	}

	if page.Props["render"] != "render-val" {
		t.Error("render prop missing")
	}
}

func TestShareProp_ConcurrentSafe(t *testing.T) {
	i := newTestInertia(t)

	done := make(chan struct{})

	for n := 0; n < 10; n++ {
		go func(n int) {
			i.ShareProp("key", n)
			_ = i.SharedProps()

			done <- struct{}{}
		}(n)
	}

	for n := 0; n < 10; n++ {
		<-done
	}
}

func TestShareProp_FuncValueResolvedOnEachRender(t *testing.T) {
	i := newTestInertia(t)

	callCount := 0

	i.ShareProp("counter", func() any {
		callCount++

		return callCount
	})

	for range 2 {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(httpx.HeaderInertia, "true")
		r.RequestURI = "/"
		w := httptest.NewRecorder()
		i.Render(w, r, "Page")
	}

	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (func should be invoked on each render)", callCount)
	}
}

// --- Validation Errors ---

func TestValidationErrors_MultipleFields(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{
		"name":  "required",
		"email": "invalid",
	})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	errors, ok := page.Props["errors"].(map[string]any)

	if !ok {
		t.Fatal("errors prop not found or not a map")
	}

	if errors["name"] != "required" {
		t.Errorf("errors.name = %v", errors["name"])
	}

	if errors["email"] != "invalid" {
		t.Errorf("errors.email = %v", errors["email"])
	}
}

func TestValidationErrors_EmptyNotAdded(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if _, ok := page.Props["errors"]; ok {
		t.Error("empty validation errors should not be added to props")
	}
}

func TestValidationErrors_NilNotAdded(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if _, ok := page.Props["errors"]; ok {
		t.Error("no SetValidationErrors call should mean no 'errors' prop")
	}
}

func TestValidationErrors_OverridesRenderErrors(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{"field": "required"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page", httpx.Props{"errors": "custom"})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	// Validation errors are applied last in mergeProps, so they override.
	errors, ok := page.Props["errors"].(map[string]any)

	if !ok {
		t.Fatal("errors should be the validation errors map, not a string")
	}

	if errors["field"] != "required" {
		t.Errorf("errors.field = %v", errors["field"])
	}
}

func TestValidationErrors_NestedStructure(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{
		"address": map[string]string{"street": "required"},
	})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["errors"] == nil {
		t.Error("nested validation errors should be preserved")
	}
}

// --- Location ---

func TestLocation_InertiaRequest_EmptyBody(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	w := httptest.NewRecorder()

	i.Location(w, r, "https://external.com")

	if w.Body.Len() > 0 {
		t.Errorf("409 response body should be empty, got %q", w.Body.String())
	}
}

func TestLocation_NonInertiaRequest_CustomStatus(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	i.Location(w, r, "/path", http.StatusMovedPermanently)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMovedPermanently)
	}
}

func TestLocation_InertiaRequest_IgnoresCustomStatus(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	w := httptest.NewRecorder()

	i.Location(w, r, "/path", http.StatusMovedPermanently)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d (always 409 for Inertia)", w.Code, http.StatusConflict)
	}
}

// --- Redirect / Back ---

func TestRedirect_SetsLocationHeader(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	i.Redirect(w, r, "/dashboard")

	if loc := w.Header().Get("Location"); loc != "/dashboard" {
		t.Errorf("Location = %q, want %q", loc, "/dashboard")
	}
}

func TestBack_WithCustomStatus303(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Referer", "/previous")
	w := httptest.NewRecorder()

	i.Back(w, r, http.StatusSeeOther)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", w.Code, http.StatusSeeOther)
	}

	if loc := w.Header().Get("Location"); loc != "/previous" {
		t.Errorf("Location = %q, want %q", loc, "/previous")
	}
}

// --- Deferred Props Integration ---

func TestRender_DeferredPropsMultipleGroups(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"stats":    props.Defer("s", "sidebar"),
		"forecast": props.Defer("f", "sidebar"),
		"logs":     props.Defer("l", "footer"),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if len(page.DeferredProps["sidebar"]) != 2 {
		t.Errorf("DeferredProps[sidebar] = %v, want 2 items", page.DeferredProps["sidebar"])
	}

	if len(page.DeferredProps["footer"]) != 1 {
		t.Errorf("DeferredProps[footer] = %v, want 1 item", page.DeferredProps["footer"])
	}
}

func TestRender_DeferredMergePropsInBothFields(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"items": props.Defer([]int{1}, "list").Merge(),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.DeferredProps["list"] == nil {
		t.Error("deferredProps should contain the group")
	}

	if len(page.MergeProps) != 1 || page.MergeProps[0] != "items" {
		t.Errorf("mergeProps = %v, want [items]", page.MergeProps)
	}
}

func TestRender_DeferredOnPartialReload_Resolved(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Page")
	r.Header.Set(httpx.HeaderPartialData, "stats")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"stats": props.Defer(func() any { return "resolved" }),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["stats"] != "resolved" {
		t.Errorf("stats = %v, want %q", page.Props["stats"], "resolved")
	}

	if len(page.DeferredProps) > 0 {
		t.Error("deferredProps should be empty on partial reload")
	}
}

// --- MergeProp Integration ---

func TestRender_MergePropsInJSON(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"items": props.Merge([]string{"a", "b"}),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if len(page.MergeProps) != 1 || page.MergeProps[0] != "items" {
		t.Errorf("mergeProps = %v, want [items]", page.MergeProps)
	}
}

func TestRender_DeepMergePropsInJSON(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"data": props.DeepMerge(map[string]int{"a": 1}),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if len(page.DeepMergeProps) != 1 || page.DeepMergeProps[0] != "data" {
		t.Errorf("deepMergeProps = %v, want [data]", page.DeepMergeProps)
	}
}

func TestRender_MergeAndDeepMergeTogether(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"shallow": props.Merge([]int{1}),
		"deep":    props.DeepMerge(map[string]int{"a": 1}),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if len(page.MergeProps) != 1 {
		t.Errorf("mergeProps = %v, want 1 entry", page.MergeProps)
	}

	if len(page.DeepMergeProps) != 1 {
		t.Errorf("deepMergeProps = %v, want 1 entry", page.DeepMergeProps)
	}
}

// --- ScrollProp Integration ---

func TestRender_ScrollPropsAllFieldsInJSON(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"feed": props.Scroll([]int{1}, "feedPage", 2, 1, 3),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	scroll := page.ScrollProps["feed"]

	if scroll.PageName != "feedPage" {
		t.Errorf("pageName = %q", scroll.PageName)
	}

	// JSON numbers decode as float64.
	if scroll.CurrentPage != float64(2) {
		t.Errorf("currentPage = %v", scroll.CurrentPage)
	}

	if scroll.PreviousPage != float64(1) {
		t.Errorf("previousPage = %v", scroll.PreviousPage)
	}

	if scroll.NextPage != float64(3) {
		t.Errorf("nextPage = %v", scroll.NextPage)
	}
}

func TestRender_ScrollPropsWithReset(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"feed": props.Scroll([]int{1}, "p", 1, nil, 2).Reset(),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if !page.ScrollProps["feed"].Reset {
		t.Error("scroll reset should be true")
	}
}

func TestRender_ScrollPropsWithMerge(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"feed": props.Scroll([]int{1}, "p", 1, nil, 2).Merge(),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if len(page.MergeProps) != 1 || page.MergeProps[0] != "feed" {
		t.Errorf("mergeProps = %v, want [feed]", page.MergeProps)
	}
}

// --- OnceProp Integration ---

func TestRender_OncePropsMetadataInJSON(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"notes": props.Once("snapshot"),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.OnceProps["notes"].Prop != "notes" {
		t.Errorf("onceProps[notes].prop = %q", page.OnceProps["notes"].Prop)
	}

	if page.Props["notes"] != "snapshot" {
		t.Errorf("notes = %v", page.Props["notes"])
	}
}

func TestRender_OnceExcluded_AbsentFromProps(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderExceptOnceProps, "notes")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page", httpx.Props{
		"notes": props.Once("snapshot"),
	})

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if _, ok := page.Props["notes"]; ok {
		t.Error("OnceProp should be absent when excluded by except-once header")
	}
}

// --- History ---

func TestRender_EncryptHistoryFromBothOptionAndContext(t *testing.T) {
	i, _ := inertia.New(testTemplate, inertia.WithVersion("v1"), inertia.WithEncryptHistory())

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetEncryptHistory(r.Context())
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if !page.EncryptHistory {
		t.Error("encryptHistory should be true (from option OR context)")
	}
}

func TestRender_EncryptHistoryDefaultFalseOmitted(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	body := w.Body.String()

	if contains(body, `"encryptHistory"`) {
		t.Error("encryptHistory should be omitted when false")
	}
}

func TestRender_ClearHistoryOmittedWhenFalse(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"
	w := httptest.NewRecorder()

	i.Render(w, r, "Page")

	body := w.Body.String()

	if contains(body, `"clearHistory"`) {
		t.Error("clearHistory should be omitted when false")
	}
}

// --- Context Helpers ---

func TestSetProp_MultipleCalls_Accumulate(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "a", "1")
	ctx = inertia.SetProp(ctx, "b", "2")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["a"] != "1" {
		t.Errorf("a = %v", page.Props["a"])
	}

	if page.Props["b"] != "2" {
		t.Errorf("b = %v", page.Props["b"])
	}
}

func TestSetProp_SameKeyOverwritten(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "key", "first")
	ctx = inertia.SetProp(ctx, "key", "second")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	if page.Props["key"] != "second" {
		t.Errorf("key = %v, want %q", page.Props["key"], "second")
	}
}

func TestSetProps_MergesWithExisting(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetProp(r.Context(), "a", "1")
	ctx = inertia.SetProps(ctx, httpx.Props{"b": "2", "c": "3"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	var page response.Page

	json.NewDecoder(w.Body).Decode(&page)

	for _, key := range []string{"a", "b", "c"} {
		if page.Props[key] == nil {
			t.Errorf("prop %q should be present", key)
		}
	}
}

func TestSetTemplateData_DoesNotAffectJSON(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.RequestURI = "/"

	ctx := inertia.SetTemplateData(r.Context(), httpx.TemplateData{"extra": "val"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if contains(body, "extra") {
		t.Error("template data should not appear in XHR JSON response")
	}
}

// --- Head management ---

func TestWithHead_RendersDefaults(t *testing.T) {
	tmpl := `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`
	i, err := inertia.New(tmpl,
		inertia.WithVersion("v1"),
		inertia.WithHead(httpx.Head{
			Title: "Default Title",
			Meta: []httpx.MetaTag{
				{Name: "description", Content: "Default desc"},
			},
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, "<title>Default Title</title>") {
		t.Errorf("missing default title, body = %s", body)
	}

	if !contains(body, `name="description"`) {
		t.Errorf("missing default description meta, body = %s", body)
	}
}

func TestWithHeadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/seo.yml"

	os.WriteFile(path, []byte(`
title: "YAML Title"
lang: "en"
meta:
  - name: "description"
    content: "From YAML"
`), 0644)

	tmpl := `<!DOCTYPE html><html lang="{{ .inertiaLang }}"><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`
	i, err := inertia.New(tmpl,
		inertia.WithVersion("v1"),
		inertia.WithHeadFromFile(path),
	)

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, "<title>YAML Title</title>") {
		t.Errorf("missing YAML title, body = %s", body)
	}

	if !contains(body, `lang="en"`) {
		t.Errorf("missing lang attribute, body = %s", body)
	}

	if !contains(body, "From YAML") {
		t.Errorf("missing YAML description, body = %s", body)
	}
}

func TestWithHeadFromFile_FileNotFound(t *testing.T) {
	tmpl := `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`
	_, err := inertia.New(tmpl,
		inertia.WithVersion("v1"),
		inertia.WithHeadFromFile("/nonexistent/seo.yml"),
	)

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWithHeadFromFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/seo.yml"

	os.WriteFile(path, []byte("title: [\ninvalid"), 0644)

	tmpl := `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`
	_, err := inertia.New(tmpl,
		inertia.WithVersion("v1"),
		inertia.WithHeadFromFile(path),
	)

	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestSetHead_OverridesDefault(t *testing.T) {
	tmpl := `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`
	i, err := inertia.New(tmpl,
		inertia.WithVersion("v1"),
		inertia.WithHead(httpx.Head{
			Title: "Default",
			Meta: []httpx.MetaTag{
				{Name: "description", Content: "Default desc"},
				{Name: "robots", Content: "index, follow"},
			},
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetTitle(r.Context(), "Override Title")
	ctx = inertia.SetMeta(ctx, httpx.MetaTag{Name: "description", Content: "Override desc"})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, "<title>Override Title</title>") {
		t.Errorf("title not overridden, body = %s", body)
	}

	if !contains(body, "Override desc") {
		t.Errorf("description not overridden, body = %s", body)
	}

	// Robots should still be present from defaults.
	if !contains(body, "index, follow") {
		t.Errorf("default robots missing, body = %s", body)
	}
}

func TestSetTitle(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetTitle(r.Context(), "My Title")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	if !contains(w.Body.String(), "<title>My Title</title>") {
		t.Errorf("missing title, body = %s", w.Body.String())
	}
}

func TestSetLang(t *testing.T) {
	tmpl := `<!DOCTYPE html><html lang="{{ .inertiaLang }}"><body>{{ .inertia }}</body></html>`
	i, _ := inertia.New(tmpl, inertia.WithVersion("v1"))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetLang(r.Context(), "fr")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	if !contains(w.Body.String(), `lang="fr"`) {
		t.Errorf("missing lang, body = %s", w.Body.String())
	}
}

func TestSetMeta(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetMeta(r.Context(),
		httpx.MetaTag{Property: "og:title", Content: "OG Test"},
		httpx.MetaTag{Name: "twitter:card", Content: "summary"},
	)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, `property="og:title"`) {
		t.Error("missing og:title")
	}

	if !contains(body, `name="twitter:card"`) {
		t.Error("missing twitter:card")
	}
}

func TestSetLinks(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetLinks(r.Context(),
		httpx.LinkTag{Rel: "canonical", Href: "https://example.com/page"},
	)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	if !contains(w.Body.String(), `rel="canonical"`) {
		t.Errorf("missing canonical link, body = %s", w.Body.String())
	}
}

func TestCSRFTokenInHead(t *testing.T) {
	i := newTestInertia(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetCSRFToken(r.Context(), "test-csrf-token")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, `name="csrf-token"`) {
		t.Errorf("missing csrf-token meta, body = %s", body)
	}

	if !contains(body, "test-csrf-token") {
		t.Errorf("missing csrf token value, body = %s", body)
	}
}

func TestLocaleHeadMerge(t *testing.T) {
	tmpl := `<!DOCTYPE html><html lang="{{ .inertiaLang }}" dir="{{ .inertiaDir }}"><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`
	i, _ := inertia.New(tmpl,
		inertia.WithVersion("v1"),
		inertia.WithHead(httpx.Head{
			Title: "Global Default",
			Meta: []httpx.MetaTag{
				{Name: "robots", Content: "index, follow"},
			},
		}),
	)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := inertia.SetLocale(r.Context(), &httpx.Locale{
		Code:      "ar",
		Name:      "Arabic",
		Direction: "rtl",
		Head: httpx.Head{
			Title: "Arabic Title",
			Meta: []httpx.MetaTag{
				{Property: "og:locale", Content: "ar_SA"},
			},
		},
	})
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if !contains(body, `lang="ar"`) {
		t.Errorf("missing lang=ar, body = %s", body)
	}

	if !contains(body, `dir="rtl"`) {
		t.Errorf("missing dir=rtl, body = %s", body)
	}

	if !contains(body, "<title>Arabic Title</title>") {
		t.Errorf("locale title not applied, body = %s", body)
	}

	// Global robots default should still be present.
	if !contains(body, "index, follow") {
		t.Errorf("global robots default missing, body = %s", body)
	}

	if !contains(body, "ar_SA") {
		t.Errorf("locale og:locale missing, body = %s", body)
	}
}

func TestHead_NotInJSON(t *testing.T) {
	i, _ := inertia.New(testTemplate,
		inertia.WithVersion("v1"),
		inertia.WithHead(httpx.Head{Title: "Should Not Appear"}),
	)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderVersion, "v1")

	ctx := inertia.SetTitle(r.Context(), "Also Not Appear")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	i.Render(w, r, "Page")

	body := w.Body.String()

	if contains(body, "Should Not Appear") || contains(body, "Also Not Appear") {
		t.Errorf("head should not appear in JSON response, body = %s", body)
	}
}
