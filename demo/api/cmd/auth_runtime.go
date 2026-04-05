package main

import (
	"net/http"
	"os"

	"github.com/oullin/inertia-go/demo/api/auth"
)

func (rt *runtime) newAuth() auth.App {
	redirectFn := func(w http.ResponseWriter, r *http.Request, url string) {
		rt.inertia.Redirect(w, r, url)
	}

	return auth.NewApp(auth.Container{
		DB:        rt.db,
		CryptoKey: rt.cryptoKey,
		Render:    rt.renderPage,
		Redirect:  redirectFn,

		RouteURL:     rt.routes.URL,
		SetFlash:     rt.flashStore.Set,
		SecureCookie: os.Getenv("APP_SECURE_COOKIES") == "true",
	})
}
