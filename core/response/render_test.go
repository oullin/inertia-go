package response_test

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/response"
)

type testMarshaler struct{}

type failMarshaler struct{}

type marshalError struct{}

func (m *testMarshaler) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }
func (m *testMarshaler) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }

func (m *failMarshaler) Marshal(v any) ([]byte, error)   { return nil, errMarshal }
func (m *failMarshaler) Unmarshal(b []byte, v any) error { return errMarshal }

var errMarshal = &marshalError{}

func (e *marshalError) Error() string { return "marshal failed" }

// --- WriteJSON ---

func TestWriteJSON(t *testing.T) {
	page := &response.Page{
		Component: "Users/Index",
		Props:     map[string]any{"name": "alice"},
		URL:       "/users",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteJSON(w, page, &testMarshaler{})

	if err != nil {
		t.Fatal(err)
	}

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	if resp.Header.Get(httpx.HeaderInertia) != "true" {
		t.Error("missing X-Inertia header")
	}

	var got response.Page

	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.Component != "Users/Index" {
		t.Errorf("component = %q", got.Component)
	}

	if got.URL != "/users" {
		t.Errorf("url = %q", got.URL)
	}
}

func TestWriteJSON_MarshalError(t *testing.T) {
	page := &response.Page{Component: "Page"}

	w := httptest.NewRecorder()
	err := response.WriteJSON(w, page, &failMarshaler{})

	if err == nil {
		t.Error("expected error from marshal failure")
	}
}

// --- WriteHTML ---

func TestWriteHTML(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(
		`<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`,
	))

	page := &response.Page{
		Component: "Home",
		Props:     map[string]any{"title": "Welcome"},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, tmpl, page, "app", &testMarshaler{}, nil)

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
		t.Error("missing empty container div")
	}

	if !contains(body, `<script data-page="app" type="application/json">`) {
		t.Error("missing page data script element")
	}

	if !contains(body, `"component":"Home"`) {
		t.Error("missing component in script JSON")
	}
}

func TestWriteHTML_MarshalError(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))
	page := &response.Page{Component: "Page"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, tmpl, page, "app", &failMarshaler{}, nil)

	if err == nil {
		t.Error("expected error from marshal failure")
	}
}

func TestWriteHTML_ExtraTemplateData(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(
		`{{ .inertia }}|title={{ .pageTitle }}`,
	))

	page := &response.Page{
		Component: "Home",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	extra := httpx.TemplateData{"pageTitle": "My App"}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, tmpl, page, "app", &testMarshaler{}, extra)

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, "title=My App") {
		t.Errorf("extra template data not rendered, body = %s", body)
	}
}

func TestWriteHTML_ScriptClosingTagEscaped(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))

	// Props containing "</script>" should be escaped to prevent injection.
	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{"html": "</script>"},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, tmpl, page, "app", &testMarshaler{}, nil)

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	// Go's json.Marshal escapes "<" as \u003c, but the explicit
	// strings.ReplaceAll("</" -> "<\/") acts as a safety net for
	// custom marshalers.
	if contains(body, "</script></script>") {
		t.Error("raw </script> found in output; should be escaped")
	}
}

func TestWriteHTML_CustomContainerID(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, tmpl, page, "root", &testMarshaler{}, nil)

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if !contains(body, `<div id="root"></div>`) {
		t.Error("missing container div with custom ID")
	}

	if !contains(body, `<script data-page="root"`) {
		t.Error("missing script element with custom container ID")
	}
}

func TestWriteHTML_OmitsEmptyOptionalFields(t *testing.T) {
	tmpl := template.Must(template.New("root").Parse(`{{ .inertia }}`))

	page := &response.Page{
		Component: "Page",
		Props:     map[string]any{},
		URL:       "/",
		Version:   "v1",
	}

	w := httptest.NewRecorder()
	err := response.WriteHTML(w, tmpl, page, "app", &testMarshaler{}, nil)

	if err != nil {
		t.Fatal(err)
	}

	body := w.Body.String()

	if contains(body, `"encryptHistory"`) {
		t.Error("encryptHistory should be omitted when false")
	}

	if contains(body, `"clearHistory"`) {
		t.Error("clearHistory should be omitted when false")
	}

	if contains(body, `"mergeProps"`) {
		t.Error("mergeProps should be omitted when empty")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
