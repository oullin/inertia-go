package main

import (
	"net/http"
	"os"

	"github.com/oullin/inertia-go/demo/api/auth"
)

func (rt *runtime) newAuthApp() auth.App {
	return auth.New(auth.Deps{
		DB:        rt.db,
		CryptoKey: rt.cryptoKey,
		Render:    rt.renderPage,
		Redirect: func(w http.ResponseWriter, r *http.Request, url string) {
			rt.inertia.Redirect(w, r, url)
		},
		RouteURL:     rt.routes.URL,
		SetFlash:     rt.flashStore.Set,
		SecureCookie: os.Getenv("APP_SECURE_COOKIES") == "true",
	})
}
