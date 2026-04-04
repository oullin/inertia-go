package main

import (
	"net/http"

	"github.com/oullin/inertia-go/demo/api/auth"
	"github.com/oullin/inertia-go/demo/api/crm"
)

func (rt *runtime) registerCRMRoutes(mux *http.ServeMux, authApp auth.App) error {
	return crm.RegisterRoutes(rt.routes, mux, crm.Container{
		DB:          rt.db,
		RequireAuth: authApp.RequireAuth,
		Render:      rt.renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			rt.inertia.Redirect(w, r, url)
		},
		RouteURL:    rt.routes.URL,
		SetFlash:    rt.flashStore.Set,
		CurrentUser: authApp.CurrentUser,
	})
}
