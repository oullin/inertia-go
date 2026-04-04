package features

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) httpExceptionsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Errors/HttpExceptions", httpx.Props{})
}

func (a app) httpExceptionsTriggerHandler(w http.ResponseWriter, r *http.Request) {
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

func (a app) networkErrorsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Errors/NetworkErrors", httpx.Props{})
}
