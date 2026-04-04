package main

import (
	"net/http"
	"os"

	"github.com/oullin/inertia-go/demo/api/auth"
)

func newAuthApp() auth.App {
	return auth.New(auth.Deps{
		DB:     db,
		Render: renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			i.Redirect(w, r, url)
		},
		RouteURL:     routes.URL,
		SetFlash:     flashStore.Set,
		SecureCookie: os.Getenv("APP_SECURE_COOKIES") == "true",
	})
}
