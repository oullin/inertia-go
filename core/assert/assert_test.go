package assert_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oullin/inertia-go/core/assert"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
)

// --- Failure paths ---

// mockTB captures assertion failures without failing the real test.
type mockTB struct {
	testing.TB
	failed bool
}

type failReader struct{}

type readError struct{}

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

func (m *mockTB) Helper()                           {}
func (m *mockTB) Errorf(format string, args ...any) { m.failed = true }
func (m *mockTB) Fatalf(format string, args ...any) { m.failed = true }

func newAssertable() assert.AssertableInertia {
	body := []byte(`{"component":"Page","props":{"name":"alice"},"url":"/","version":"v1"}`)

	return assert.AssertFromBytes(&mockTB{}, body)
}

func TestAssertComponent_Failure(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertComponent(m, "WrongComponent")

	if !m.failed {
		t.Error("expected assertion failure for wrong component")
	}
}

func TestAssertURL_Failure(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertURL(m, "/wrong")

	if !m.failed {
		t.Error("expected assertion failure for wrong URL")
	}
}

func TestAssertVersion_Failure(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertVersion(m, "wrong")

	if !m.failed {
		t.Error("expected assertion failure for wrong version")
	}
}

func TestAssertHasProp_Failure(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertHasProp(m, "nonexistent")

	if !m.failed {
		t.Error("expected assertion failure for missing prop")
	}
}

func TestAssertPropEquals_Failure(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertPropEquals(m, "name", "wrong-value")

	if !m.failed {
		t.Error("expected assertion failure for wrong prop value")
	}
}

func TestAssertPropEquals_MissingKey(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertPropEquals(m, "missing", "any")

	if !m.failed {
		t.Error("expected assertion failure for missing key")
	}
}

func TestAssertMissingProp_Failure(t *testing.T) {
	m := &mockTB{}
	a := newAssertable()
	a.AssertMissingProp(m, "name")

	if !m.failed {
		t.Error("expected assertion failure when prop exists")
	}
}

func TestAssertFromBytes_InvalidJSON(t *testing.T) {
	m := &mockTB{}
	assert.AssertFromBytes(m, []byte(`not json`))

	if !m.failed {
		t.Error("expected failure for invalid JSON")
	}
}

func TestAssertFromReader_ReadError(t *testing.T) {
	m := &mockTB{}
	assert.AssertFromReader(m, &failReader{})

	if !m.failed {
		t.Error("expected failure for read error")
	}
}

func TestAssertFromHandler_Non200(t *testing.T) {
	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	m := &mockTB{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"
	assert.AssertFromHandler(m, i, handler, r)

	if !m.failed {
		t.Error("expected failure for non-200 response")
	}
}

func (f *failReader) Read(p []byte) (int, error) {
	return 0, &readError{}
}

func (e *readError) Error() string { return "read failed" }
