package main

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/crm"
)

func registerCRMRoutes(mux *http.ServeMux) {
	crm.RegisterRoutes(mux, crm.Deps{
		DB:          db,
		RequireAuth: requireDemoAuth,
		Render:      renderPage,
		RenderWithContext: func(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
			renderPageWithContext(w, r, component, pageProps)
		},
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Redirect(w, r, url)
		},
		RouteURL: routeURL,
		SetFlash: func(w http.ResponseWriter, flash crm.Flash) {
			setFlash(w, flashPayload{
				Kind:    flash.Kind,
				Title:   flash.Title,
				Message: flash.Message,
			})
		},
		CurrentUser: currentUser,
	})
}
