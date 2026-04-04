package errors

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/wayfinder"
)

// Container contains the host application integrations required by the errors package.
type Container struct {
	RequireAuth func(http.Handler) http.Handler
	Render      func(http.ResponseWriter, *http.Request, string, httpx.Props)
}

// RegisterRoutes mounts the error showcase HTTP routes onto the provided mux.
func RegisterRoutes(routes *wayfinder.Registry, mux *http.ServeMux, deps Container) {
	auth := func(h http.HandlerFunc) http.Handler {
		return deps.RequireAuth(h)
	}

	routes.Handle("features.errors.http-error", auth(httpErrorHandler(deps)), mux)
	mux.Handle("/features/errors/http-error/{code}", auth(httpErrorTriggerHandler()))
	routes.Handle("features.errors.network-errors", auth(networkErrorsHandler(deps)), mux)
}

func httpErrorHandler(deps Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deps.Render(w, r, "Features/Errors/HttpError", httpx.Props{})
	}
}

func httpErrorTriggerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		switch code {
		case "403":
			http.Error(w, "Forbidden", http.StatusForbidden)
		case "404":
			http.Error(w, "Not Found", http.StatusNotFound)
		case "500":
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		case "unhandled":
			http.Error(w, "I'm a teapot", http.StatusTeapot)
		default:
			http.NotFound(w, r)
		}
	}
}

func networkErrorsHandler(deps Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deps.Render(w, r, "Features/Errors/NetworkErrors", httpx.Props{})
	}
}
