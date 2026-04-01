package props_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/props"
)

// No X-Inertia-Partial-Data header — onlySet is empty.

type tryPropValue struct{ val any }

func TestResolve_FullRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"users":  []string{"alice", "bob"},
		"title":  "Users Page",
		"lazy":   func() any { return "computed" },
		"always": props.Always("always-val"),
	}

	result, err := props.Resolve(r, "Users/Index", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["users"] == nil {
		t.Error("expected 'users' prop")
	}

	if result.Props["title"] != "Users Page" {
		t.Errorf("title = %v", result.Props["title"])
	}

	if result.Props["lazy"] != "computed" {
		t.Errorf("lazy = %v", result.Props["lazy"])
	}

	if result.Props["always"] != "always-val" {
		t.Errorf("always = %v", result.Props["always"])
	}
}

func TestResolve_DeferredExcludedOnInitial(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"name":  "test",
		"stats": props.Defer(func() any { return "expensive" }, "sidebar"),
	}

	result, err := props.Resolve(r, "Dashboard", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["stats"]; ok {
		t.Error("deferred prop 'stats' should be excluded on initial load")
	}

	if result.Props["name"] != "test" {
		t.Errorf("name = %v", result.Props["name"])
	}

	sidebarGroup := result.DeferredProps["sidebar"]

	if len(sidebarGroup) != 1 || sidebarGroup[0] != "stats" {
		t.Errorf("DeferredProps[sidebar] = %v, want [stats]", sidebarGroup)
	}
}

func TestResolve_DeferredIncludedOnPartialReload(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Dashboard")
	r.Header.Set(httpx.HeaderPartialData, "stats")

	merged := httpx.Props{
		"name":  "test",
		"stats": props.Defer(func() any { return "expensive" }),
	}

	result, err := props.Resolve(r, "Dashboard", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["stats"] != "expensive" {
		t.Errorf("stats = %v, want %q", result.Props["stats"], "expensive")
	}

	if _, ok := result.Props["name"]; ok {
		t.Error("'name' should be excluded from partial reload requesting only 'stats'")
	}
}

func TestResolve_OptionalExcludedOnInitial(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"name":     "test",
		"optional": props.Optional("opt-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("optional prop should be excluded on initial load")
	}
}

func TestResolve_OptionalIncludedOnPartialReload(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Page")
	r.Header.Set(httpx.HeaderPartialData, "optional")

	merged := httpx.Props{
		"name":     "test",
		"optional": props.Optional("opt-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["optional"] != "opt-val" {
		t.Errorf("optional = %v, want %q", result.Props["optional"], "opt-val")
	}
}

func TestResolve_OnceExcludedWhenInExceptHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderExceptOnceProps, "meta")

	merged := httpx.Props{
		"name": "test",
		"meta": props.Once("metadata"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["meta"]; ok {
		t.Error("once prop 'meta' should be excluded when in X-Inertia-Except-Once-Props")
	}
}

func TestResolve_OnceIncludedOnFirstVisit(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"meta": props.Once("metadata"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["meta"] != "metadata" {
		t.Errorf("meta = %v, want %q", result.Props["meta"], "metadata")
	}
}

func TestResolve_MergeRecorded(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"posts": props.Merge([]string{"a", "b"}),
		"deep":  props.DeepMerge(map[string]int{"x": 1}),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "posts" {
		t.Errorf("MergeProps = %v, want [posts]", result.MergeProps)
	}

	if len(result.DeepMergeProps) != 1 || result.DeepMergeProps[0] != "deep" {
		t.Errorf("DeepMergeProps = %v, want [deep]", result.DeepMergeProps)
	}
}

func TestResolve_PartialExceptTakesPrecedence(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Page")
	r.Header.Set(httpx.HeaderPartialData, "a,b,c")
	r.Header.Set(httpx.HeaderPartialExcept, "b")

	merged := httpx.Props{
		"a": "val-a",
		"b": "val-b",
		"c": "val-c",
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["b"]; ok {
		t.Error("prop 'b' should be excluded by X-Inertia-Partial-Except")
	}

	if result.Props["a"] != "val-a" {
		t.Errorf("a = %v", result.Props["a"])
	}

	if result.Props["c"] != "val-c" {
		t.Errorf("c = %v", result.Props["c"])
	}
}

func TestResolve_AlwaysIncludedInPartialReload(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Page")
	r.Header.Set(httpx.HeaderPartialData, "name")

	merged := httpx.Props{
		"name":   "test",
		"always": props.Always("always-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["always"] != "always-val" {
		t.Error("AlwaysProp should be included even when not in partial data list")
	}
}

func TestResolve_LazyFuncWithError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"failing": func() (any, error) {
			return nil, http.ErrAbortHandler
		},
	}

	_, err := props.Resolve(r, "Page", merged)

	if err == nil {
		t.Error("expected error from lazy func")
	}
}

func TestResolve_OptionalExcludedOnPartialWithoutOnly(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Page")

	merged := httpx.Props{
		"name":     "test",
		"optional": props.Optional("opt-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("optional prop should be excluded on partial reload without explicit request")
	}

	if result.Props["name"] != "test" {
		t.Errorf("name = %v", result.Props["name"])
	}
}

func TestResolve_DeferredMergeOnInitial(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"items": props.Defer(func() any { return []string{"a"} }, "list").Merge(),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["items"]; ok {
		t.Error("deferred prop should be excluded on initial load")
	}

	if len(result.DeferredProps["list"]) != 1 || result.DeferredProps["list"][0] != "items" {
		t.Errorf("DeferredProps[list] = %v", result.DeferredProps["list"])
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "items" {
		t.Errorf("MergeProps = %v, want [items]", result.MergeProps)
	}
}

func TestResolve_DeferredMergeOnPartialReload(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")
	r.Header.Set(httpx.HeaderPartialComponent, "Page")
	r.Header.Set(httpx.HeaderPartialData, "items")

	merged := httpx.Props{
		"items": props.Defer(func() any { return []string{"a"} }, "list").Merge(),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["items"] == nil {
		t.Error("deferred merge prop should be included in partial reload")
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "items" {
		t.Errorf("MergeProps = %v, want [items]", result.MergeProps)
	}
}

func (tp tryPropValue) TryProp() (any, error) { return tp.val, nil }

func TestResolve_TryProper(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(httpx.HeaderInertia, "true")

	merged := httpx.Props{
		"custom": tryPropValue{val: "resolved-via-try"},
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["custom"] != "resolved-via-try" {
		t.Errorf("custom = %v, want %q", result.Props["custom"], "resolved-via-try")
	}
}
