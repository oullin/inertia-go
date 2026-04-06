package crm

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/assert"
	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/wayfinder"
	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/seed"
	"github.com/oullin/inertia-go/demo/api/internal/testutil"
)

type crmHarness struct {
	t        *testing.T
	db       *sql.DB
	inertia  *inertia.Inertia
	registry *wayfinder.Registry
	app      app
	user     *database.User
	flashes  []flash.Message
}

func newCRMHarness(t *testing.T) *crmHarness {
	t.Helper()

	i, err := inertia.New(testutil.TestTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatalf("inertia.New() error = %v", err)
	}

	db, err := database.Open(":memory:")

	if err != nil {
		t.Fatalf("database.Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := seed.Run(db); err != nil {
		t.Fatalf("seed.Run() error = %v", err)
	}

	user, err := database.FindUserByEmail(db, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	h := &crmHarness{
		t:        t,
		db:       db,
		inertia:  i,
		registry: buildCRMRegistry(),
		user:     user,
	}

	requireAuthFn := func(next http.Handler) http.Handler { return next }
	renderFn := func(w http.ResponseWriter, r *http.Request, component string, props httpx.Props) {
		if strings.TrimSpace(r.RequestURI) == "" {
			r.RequestURI = r.URL.RequestURI()
		}

		if err := i.Render(w, r, component, props); err != nil {
			t.Fatalf("Render(%s) error = %v", component, err)
		}
	}
	redirectFn := func(w http.ResponseWriter, r *http.Request, url string) {
		i.Redirect(w, r, url)
	}
	setFlashFn := func(w http.ResponseWriter, msg flash.Message) error {
		h.flashes = append(h.flashes, msg)

		return nil
	}
	currentUserFn := func(r *http.Request) *database.User {
		return h.user
	}

	h.app, err = newApp(Container{
		DB:          db,
		RequireAuth: requireAuthFn,
		Render:      renderFn,

		Redirect: redirectFn,

		RouteURL: h.registry.URL,
		SetFlash: setFlashFn,

		CurrentUser: currentUserFn,
	})

	if err != nil {
		t.Fatalf("newApp() error = %v", err)
	}

	return h
}

func buildCRMRegistry() *wayfinder.Registry {
	routes := wayfinder.New()

	routes.Add("login", "GET", "/login")
	routes.Add("dashboard", "GET", "/dashboard")

	routes.Group("contacts", "/contacts", func(g *wayfinder.Group) {
		g.Add("index", "GET", "")
		g.Add("create", "GET", "/create")
		g.Add("show", "GET", "/{contact}")
		g.Add("edit", "GET", "/{contact}/edit")
		g.Add("update", "POST", "/{contact}")
		g.Add("destroy", "DELETE", "/{contact}")
		g.Add("favorite", "POST", "/{contact}/favorite")

		g.Group("notes", "", func(ng *wayfinder.Group) {
			ng.Add("store", "POST", "/{contact}/notes")
		})
	})
	routes.Group("organizations", "/organizations", func(g *wayfinder.Group) {
		g.Add("index", "GET", "")
		g.Add("show", "GET", "/{organization}")
		g.Add("update", "POST", "/{organization}")
	})

	return routes
}

func (h *crmHarness) request(method, target string) *http.Request {
	h.t.Helper()

	req := httptest.NewRequest(method, target, nil)

	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()

	return req
}

func (h *crmHarness) formRequest(method, target, body string) *http.Request {
	h.t.Helper()

	req := httptest.NewRequest(method, target, strings.NewReader(body))

	req.Header.Set(httpx.HeaderInertia, "true")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.RequestURI = req.URL.RequestURI()

	return req
}

func (h *crmHarness) partialRequest(target, component, only string) *http.Request {
	h.t.Helper()

	req := h.request(http.MethodGet, target)

	req.Header.Set(httpx.HeaderPartialComponent, component)
	req.Header.Set(httpx.HeaderPartialData, only)

	return req
}

func (h *crmHarness) page(t *testing.T, w *httptest.ResponseRecorder) assert.AssertableInertia {
	t.Helper()

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	return assert.AssertFromBytes(t, w.Body.Bytes())
}
