package main

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/wayfinder"
	"github.com/oullin/inertia-go/demo/api/auth"
)

var routes *wayfinder.Registry

func initRoutes() {
	routes = wayfinder.New()
	routes.Add("login", "GET", "/login")
	routes.Add("logout", "POST", "/logout")
	routes.Add("dashboard", "GET", "/dashboard")
	routes.Group("contacts", "/contacts", func(g *wayfinder.Group) {
		g.Add("index", "GET", "")
		g.Add("create", "GET", "/create")
		g.Add("store", "POST", "")
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

	// Feature showcase routes
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
	routes.Group("features.errors", "/features/errors", func(g *wayfinder.Group) {
		g.Add("http-exceptions", "GET", "/http-exceptions")
		g.Add("network-errors", "GET", "/network-errors")
	})
	routes.Group("features.http", "/features/http", func(g *wayfinder.Group) {
		g.Add("use-http", "GET", "/use-http")
	})
}

func withDemoProps(authApp auth.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := authApp.CurrentUser(r)

		sidebarOpen := true

		if cookie, err := r.Cookie("sidebar_open"); err == nil {
			sidebarOpen = cookie.Value != "false"
		}

		ctx := r.Context()
		ctx = inertia.SetProps(ctx, httpx.Props{
			"sidebarOpen": sidebarOpen,
			"app": map[string]any{
				"name":        "Inertia.js Kitchen Sink",
				"productLine": "Go Demo Port",
				"environment": "Demo",
			},
			"auth": map[string]any{
				"user": authApp.PublicUser(user),
			},
			"workspace": map[string]any{
				"name": "Inertia Go",
				"plan": "Porting",
			},
			"routes": routes.ManifestProps(),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func renderPage(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	renderPageWithContext(w, r, component, pageProps)
}

func renderPageWithContext(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	ctx := r.Context()

	if err := i.Render(w, r.WithContext(ctx), component, pageProps); err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			http.Error(w, "page not found", http.StatusNotFound)
		default:
			http.Error(w, "demo internal error", http.StatusInternalServerError)
		}
	}
}
