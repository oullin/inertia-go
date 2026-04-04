package main

import (
	"net/http"

	"github.com/oullin/inertia-go/demo/api/auth"
	"github.com/oullin/inertia-go/demo/api/features"
)

func registerFeatureRoutes(mux *http.ServeMux, authApp auth.App) {
	features.RegisterRoutes(mux, features.Deps{
		DB:          db,
		RequireAuth: authApp.RequireAuth,
		Render:      renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Redirect(w, r, url)
		},
		Location: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Location(w, r, url)
		},
		RouteURL: routes.URL,
		SetFlash: flashStore.Set,
	})
}
