package wayfinder

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateTypeScript(t *testing.T) {
	reg := New()
	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "type RouteResult = { url: string; method: string }") {
		t.Error("expected RouteResult type declaration")
	}

	if !strings.Contains(output, "export function login(): RouteResult {") {
		t.Error("expected flat login function")
	}

	if !strings.Contains(output, "return { url: '/login', method: 'get' }") {
		t.Error("expected login return statement")
	}

	if !strings.Contains(output, "export function contactsShow(params: { contact: string | number }): RouteResult {") {
		t.Error("expected flat contactsShow function with typed params")
	}

	if !strings.Contains(output, "return { url: `/contacts/${params.contact}`, method: 'get' }") {
		t.Error("expected contactsShow template literal return")
	}

	if !strings.Contains(output, "} as const") {
		t.Error("expected 'as const' for nested exports")
	}
}

func TestGenerateJavaScript(t *testing.T) {
	reg := New()
	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: false})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if strings.Contains(output, "RouteResult") {
		t.Error("JS output should not contain TypeScript types")
	}

	if !strings.Contains(output, "export function login() {") {
		t.Error("expected JS login function without types")
	}

	if !strings.Contains(output, "(params) =>") {
		t.Error("expected JS nested function without types")
	}
}

func TestGenerateFlatOnly(t *testing.T) {
	reg := New()
	reg.Add("contacts.index", "GET", "/contacts")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{FlatOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export function contactsIndex()") {
		t.Error("expected flat export")
	}

	if strings.Contains(output, "export const contacts") {
		t.Error("should not contain nested exports in FlatOnly mode")
	}
}

func TestGenerateNestedOnly(t *testing.T) {
	reg := New()
	reg.Add("contacts.index", "GET", "/contacts")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if strings.Contains(output, "export function contactsIndex()") {
		t.Error("should not contain flat exports in NestedOnly mode")
	}

	if !strings.Contains(output, "export const contacts") {
		t.Error("expected nested export")
	}
}

func TestGenerateCustomHeader(t *testing.T) {
	reg := New()
	reg.Add("login", "GET", "/login")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{
		Header: "// Custom header",
	})

	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(buf.String(), "// Custom header\n") {
		t.Error("expected custom header")
	}
}

func TestGenerateNestedGrouping(t *testing.T) {
	reg := New()
	reg.Add("contacts.index", "GET", "/contacts")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")
	reg.Add("contacts.store", "POST", "/contacts")
	reg.Add("organizations.index", "GET", "/organizations")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true, TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export const contacts = {") {
		t.Error("expected contacts group")
	}

	if !strings.Contains(output, "export const organizations = {") {
		t.Error("expected organizations group")
	}

	if !strings.Contains(output, "index: (): RouteResult =>") {
		t.Error("expected index member")
	}

	if !strings.Contains(output, "show: (params: { contact: string | number }): RouteResult =>") {
		t.Error("expected show member with params")
	}
}

func TestGenerateTopLevelRoutes(t *testing.T) {
	reg := New()
	reg.Add("login", "GET", "/login")
	reg.Add("logout", "POST", "/logout")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{NestedOnly: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	// Top-level routes (no dot) go into the "app" group.
	if !strings.Contains(output, "export const app = {") {
		t.Error("expected top-level routes grouped under 'app'")
	}
}

func TestDotToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"contacts.index", "contactsIndex"},
		{"contacts.notes.store", "contactsNotesStore"},
		{"login", "login"},
		{"a.b.c.d", "aBCD"},
		{"use-form", "useForm"},
		{"forms.use-form", "formsUseForm"},
		{"data-loading.deferred-props", "dataLoadingDeferredProps"},
	}

	for _, tt := range tests {
		got := dotToCamel(tt.input)

		if got != tt.expected {
			t.Errorf("dotToCamel(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGenerateHyphenatedRoutes(t *testing.T) {
	reg := New()
	reg.Add("features.forms.use-form", "GET", "/features/forms/use-form")
	reg.Add("use-http", "GET", "/use-http")

	var buf bytes.Buffer

	err := Generate(reg, &buf, GenerateOptions{TypeScript: true})

	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	if !strings.Contains(output, "export function featuresFormsUseForm(): RouteResult {") {
		t.Error("expected hyphenated flat function name to be camelCased")
	}

	if !strings.Contains(output, "export function useHttp(): RouteResult {") {
		t.Error("expected top-level hyphenated flat function name to be camelCased")
	}

	if !strings.Contains(output, "formsUseForm: (): RouteResult =>") {
		t.Error("expected hyphenated nested key to be camelCased")
	}

	if !strings.Contains(output, "useHttp: (): RouteResult =>") {
		t.Error("expected top-level hyphenated route in app group to be camelCased")
	}
}

func TestBuildURLTemplate(t *testing.T) {
	tests := []struct {
		pattern  string
		expected string
	}{
		{"/contacts", "`/contacts`"},
		{"/contacts/{contact}", "`/contacts/${params.contact}`"},
		{"/contacts/{contact}/notes/{note}", "`/contacts/${params.contact}/notes/${params.note}`"},
	}

	for _, tt := range tests {
		got := buildURLTemplate(tt.pattern)

		if got != tt.expected {
			t.Errorf("buildURLTemplate(%q) = %q, want %q", tt.pattern, got, tt.expected)
		}
	}
}
