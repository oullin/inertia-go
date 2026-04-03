package main

import (
	"net/http"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/demo/api/auth"
	"github.com/oullin/inertia-go/demo/api/crm"
)

func registerCRMRoutes(mux *http.ServeMux, authApp auth.App) {
	crm.RegisterRoutes(mux, crm.Deps{
		DB:          db,
		RequireAuth: authApp.RequireAuth,
		Render:      renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Redirect(w, r, url)
		},
		RouteURL: routes.URL,
		SetFlash: func(w http.ResponseWriter, msg flash.Message) {
			flashStore.Set(w, msg)
		},
		CurrentUser: authApp.CurrentUser,
	})
}
