package testing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ihttp "github.com/oullin/inertia-go/core/http"
	"github.com/oullin/inertia-go/core/inertia"
	itesting "github.com/oullin/inertia-go/core/testing"
)

const testTemplate = `<!DOCTYPE html>
<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}</body>
</html>`

func TestAssertFromBytes(t *testing.T) {
	body := []byte(`{"component":"Users/Index","props":{"name":"alice"},"url":"/users","version":"v1"}`)

	a := itesting.AssertFromBytes(t, body)
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
		i.Render(w, r, "Dashboard", ihttp.Props{
			"title": "Test Dashboard",
		})
	})

	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	r.RequestURI = "/dashboard"

	a := itesting.AssertFromHandler(t, i, handler, r)
	a.AssertComponent(t, "Dashboard")
	a.AssertURL(t, "/dashboard")
	a.AssertVersion(t, "v1")
	a.AssertPropEquals(t, "title", "Test Dashboard")
}
