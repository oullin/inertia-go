package main

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/auth"
	"github.com/oullin/inertia-go/demo/api/internal/flash"
)

func newAuthApp() auth.App {
	return auth.New(auth.Deps{
		DB:     db,
		Render: renderPage,
		RenderWithContext: func(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
			renderPageWithContext(w, r, component, pageProps)
		},
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Redirect(w, r, url)
		},
		RouteURL: routeURL,
		SetFlash: func(w http.ResponseWriter, message auth.Flash) {
			setFlash(w, flash.Message{
				Kind:    message.Kind,
				Title:   message.Title,
				Message: message.Message,
			})
		},
	})
}
