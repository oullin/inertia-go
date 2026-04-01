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
