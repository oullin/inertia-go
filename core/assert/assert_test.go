package assert_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/assert"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
)

const testTemplate = `<!DOCTYPE html>
<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}</body>
</html>`

func TestAssertFromBytes(t *testing.T) {
	body := []byte(`{"component":"Users/Index","props":{"name":"alice"},"url":"/users","version":"v1"}`)

	a := assert.AssertFromBytes(t, body)
	a.AssertComponent(t, "Users/Index")
	a.AssertURL(t, "/users")
	a.AssertVersion(t, "v1")
	a.AssertHasProp(t, "name")
	a.AssertPropEquals(t, "name", "alice")
	a.AssertMissingProp(t, "nonexistent")
}

func TestAssertFromHandler(t *testing.T) {
	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "Dashboard", httpx.Props{
			"title": "Test Dashboard",
		})
	})

	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	r.RequestURI = "/dashboard"

	a := assert.AssertFromHandler(t, i, handler, r)
	a.AssertComponent(t, "Dashboard")
	a.AssertURL(t, "/dashboard")
	a.AssertVersion(t, "v1")
	a.AssertPropEquals(t, "title", "Test Dashboard")
}
