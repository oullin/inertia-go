package main

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/demo/api/auth"
)

var demoRouteManifest = map[string]string{
	"login":                "/login",
	"logout":               "/logout",
	"dashboard":            "/dashboard",
	"contacts.index":       "/contacts",
	"contacts.create":      "/contacts/create",
	"contacts.store":       "/contacts",
	"contacts.show":        "/contacts/{contact}",
	"contacts.edit":        "/contacts/{contact}/edit",
	"contacts.update":      "/contacts/{contact}",
	"contacts.favorite":    "/contacts/{contact}/favorite",
	"contacts.notes.store": "/contacts/{contact}/notes",
	"organizations.index":  "/organizations",
	"organizations.show":   "/organizations/{organization}",
	"organizations.update": "/organizations/{organization}",
}

func withDemoProps(authApp auth.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := authApp.CurrentUser(r)
		ctx := r.Context()
		ctx = inertia.SetProps(ctx, httpx.Props{
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
			"routes": manifestProps(),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func renderPage(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	renderPageWithContext(w, r, component, pageProps)
}

func renderPageWithContext(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	ctx := r.Context()

	if flash := consumeFlash(w, r); flash != nil {
		ctx = inertia.SetProp(ctx, "flash", flash)
	}

	if err := i.Render(w, r.WithContext(ctx), component, pageProps); err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			http.Error(w, "page not found", http.StatusNotFound)
		default:
			http.Error(w, "demo internal error", http.StatusInternalServerError)
		}
	}
}

func manifestProps() map[string]any {
	props := make(map[string]any, len(demoRouteManifest))

	for name, pattern := range demoRouteManifest {
		props[name] = pattern
	}

	return props
}

func routeURL(name string, params map[string]string) string {
	pattern, ok := demoRouteManifest[name]

	if !ok {
		return "/"
	}

	for key, value := range params {
		pattern = strings.ReplaceAll(pattern, "{"+key+"}", value)
	}

	return pattern
}
