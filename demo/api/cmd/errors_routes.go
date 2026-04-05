package main

import (
	"net/http"

	"github.com/oullin/inertia-go/demo/api/auth"
	apierrors "github.com/oullin/inertia-go/demo/api/errors"
)

func (rt *runtime) registerErrorRoutes(mux *http.ServeMux, authApp auth.App) error {
	return apierrors.RegisterRoutes(rt.routes, mux, apierrors.Container{
		RequireAuth: authApp.RequireAuth,
		Render:      rt.renderPage,
	})
}
