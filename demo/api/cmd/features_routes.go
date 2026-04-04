package main

import (
	"net/http"

	"github.com/oullin/inertia-go/demo/api/auth"
	"github.com/oullin/inertia-go/demo/api/features"
)

func (rt *runtime) registerFeatureRoutes(mux *http.ServeMux, authApp auth.App) {
	features.RegisterRoutes(rt.routes, mux, features.Container{
		DB:          rt.db,
		RequireAuth: authApp.RequireAuth,
		Render:      rt.renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			rt.inertia.Redirect(w, r, url)
		},
		Location: func(w http.ResponseWriter, r *http.Request, url string) {
			rt.inertia.Location(w, r, url)
		},
		RouteURL: rt.routes.URL,
		SetFlash: rt.flashStore.Set,
	})
}
