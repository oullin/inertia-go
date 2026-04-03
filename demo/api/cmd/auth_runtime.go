package main

import (
	"net/http"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/demo/api/auth"
)

func newAuthApp() auth.App {
	return auth.New(auth.Deps{
		DB:     db,
		Render: renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Redirect(w, r, url)
		},
		RouteURL: routes.URL,
		SetFlash: func(w http.ResponseWriter, msg flash.Message) {
			flashStore.Set(w, msg)
		},
	})
}
