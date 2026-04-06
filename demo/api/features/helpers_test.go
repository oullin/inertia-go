package features

import (
	"bytes"
	"database/sql"
	"mime/multipart"
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

type featureHarness struct {
	t        *testing.T
	db       *sql.DB
	inertia  *inertia.Inertia
	registry *wayfinder.Registry
	app      app
	flashes  []flash.Message
}

func newFeatureHarness(t *testing.T) *featureHarness {
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

	h := &featureHarness{
		t:        t,
		db:       db,
		inertia:  i,
		registry: buildFeatureRegistry(),
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
	locationFn := func(w http.ResponseWriter, r *http.Request, url string) {
		i.Location(w, r, url)
	}
	setFlashFn := func(w http.ResponseWriter, msg flash.Message) error {
		h.flashes = append(h.flashes, msg)

		return nil
	}

	h.app = newApp(Container{
		DB:          db,
		RequireAuth: requireAuthFn,
		Render:      renderFn,

		Redirect: redirectFn,

		Location: locationFn,

		RouteURL: h.registry.URL,
		SetFlash: setFlashFn,
	})

	return h
}

func buildFeatureRegistry() *wayfinder.Registry {
	routes := wayfinder.New()

	routes.Group("features.forms", "/features/forms", func(g *wayfinder.Group) {
		g.Add("use-form", "GET", "/use-form")
		g.Add("form-component", "GET", "/form-component")
		g.Add("file-uploads", "GET", "/file-uploads")
		g.Add("validation", "GET", "/validation")
		g.Add("precognition", "GET", "/precognition")
		g.Add("optimistic-updates", "GET", "/optimistic-updates")
		g.Add("use-form-context", "GET", "/use-form-context")
		g.Add("dotted-keys", "GET", "/dotted-keys")
		g.Add("wayfinder", "GET", "/wayfinder")
	})
	routes.Group("features.navigation", "/features/navigation", func(g *wayfinder.Group) {
		g.Add("links", "GET", "/links")
		g.Add("preserve-state", "GET", "/preserve-state")
		g.Add("preserve-scroll", "GET", "/preserve-scroll")
		g.Add("view-transitions", "GET", "/view-transitions")
		g.Add("history-management", "GET", "/history-management")
		g.Add("async-requests", "GET", "/async-requests")
		g.Add("async-slow", "GET", "/async-slow")
		g.Add("manual-visits", "GET", "/manual-visits")
		g.Add("redirects", "GET", "/redirects")
		g.Add("scroll-management", "GET", "/scroll-management")
		g.Add("instant-visits", "GET", "/instant-visits")
		g.Add("instant-visit-target", "GET", "/instant-visit-target")
		g.Add("url-fragments", "GET", "/url-fragments")
	})
	routes.Group("features.data-loading", "/features/data-loading", func(g *wayfinder.Group) {
		g.Add("deferred-props", "GET", "/deferred-props")
		g.Add("partial-reloads", "GET", "/partial-reloads")
		g.Add("infinite-scroll", "GET", "/infinite-scroll")
		g.Add("when-visible", "GET", "/when-visible")
		g.Add("polling", "GET", "/polling")
		g.Add("prop-merging", "GET", "/prop-merging")
		g.Add("optional-props", "GET", "/optional-props")
		g.Add("once-props", "GET", "/once-props/{page}")
	})
	routes.Group("features.prefetching", "/features/prefetching", func(g *wayfinder.Group) {
		g.Add("link-prefetch", "GET", "/link-prefetch")
		g.Add("stale-while-revalidate", "GET", "/stale-while-revalidate")
		g.Add("manual-prefetch", "GET", "/manual-prefetch")
		g.Add("cache-management", "GET", "/cache-management")
	})
	routes.Group("features.state", "/features/state", func(g *wayfinder.Group) {
		g.Add("remember", "GET", "/remember")
		g.Add("flash-data", "GET", "/flash-data")
		g.Add("shared-props", "GET", "/shared-props")
	})
	routes.Group("features.layouts", "/features/layouts", func(g *wayfinder.Group) {
		g.Add("persistent-layouts", "GET", "/persistent-layouts")
		g.Add("persistent-layouts-page-2", "GET", "/persistent-layouts/page-2")
		g.Add("nested-layouts", "GET", "/nested-layouts")
		g.Add("head", "GET", "/head")
		g.Add("layout-props", "GET", "/layout-props")
	})
	routes.Group("features.events", "/features/events", func(g *wayfinder.Group) {
		g.Add("global-events", "GET", "/global-events")
		g.Add("visit-callbacks", "GET", "/visit-callbacks")
		g.Add("progress", "GET", "/progress")
		g.Add("progress-slow", "GET", "/progress/slow")
	})
	routes.Group("features.http", "/features/http", func(g *wayfinder.Group) {
		g.Add("use-http", "GET", "/use-http")
	})

	return routes
}

func (h *featureHarness) request(method, target string, body []byte) *http.Request {
	h.t.Helper()

	req := httptest.NewRequest(method, target, bytes.NewReader(body))

	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()

	return req
}

func (h *featureHarness) partialRequest(method, target, component, only string) *http.Request {
	h.t.Helper()

	req := h.request(method, target, nil)

	req.Header.Set(httpx.HeaderPartialComponent, component)
	req.Header.Set(httpx.HeaderPartialData, only)

	return req
}

func (h *featureHarness) page(t *testing.T, w *httptest.ResponseRecorder) assert.AssertableInertia {
	t.Helper()

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	return assert.AssertFromBytes(t, w.Body.Bytes())
}

func (h *featureHarness) lastFlash(t *testing.T) flash.Message {
	t.Helper()

	if len(h.flashes) == 0 {
		t.Fatal("no flash recorded")
	}

	return h.flashes[len(h.flashes)-1]
}

func newMultipartRequest(t *testing.T, target string, files map[string][]string) (*http.Request, string) {
	t.Helper()

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	for field, names := range files {
		for _, name := range names {
			part, err := writer.CreateFormFile(field, name)

			if err != nil {
				t.Fatalf("CreateFormFile(%s, %s) error = %v", field, name, err)
			}

			if _, err := part.Write([]byte("demo")); err != nil {
				t.Fatalf("part.Write() error = %v", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, &body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()

	return req, writer.FormDataContentType()
}
