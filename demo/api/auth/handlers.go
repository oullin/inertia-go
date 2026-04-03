package auth

import (
	"errors"
	"net/http"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
)

func (a App) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Auth/Login", httpx.Props{
			"status": r.URL.Query().Get("status"),
		})
	case http.MethodPost:
		a.loginSubmitHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a App) loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	form := newLoginForm(r)
	errorsByField := form.validate()

	if len(errorsByField) == 0 {
		user, err := a.service.authenticate(form)

		switch {
		case err == nil:
			a.setSession(w, user.ID, form.Remember)
			a.deps.SetFlash(w, flash.Message{
				Kind:    "success",
				Title:   "Signed in",
				Message: "The demo session is now authenticated.",
			})
			a.deps.Redirect(w, r, a.deps.RouteURL("dashboard", nil))

			return
		case errors.Is(err, errInvalidCredentials):
			errorsByField["email"] = "Use test@example.com and password to sign in."
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	ctx := inertia.SetValidationErrors(r.Context(), errorsByField)
	a.deps.RenderWithContext(w, r.WithContext(ctx), "Auth/Login", httpx.Props{})
}

func (a App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	a.clearSession(w)
	a.deps.SetFlash(w, flash.Message{
		Kind:    "info",
		Title:   "Signed out",
		Message: "Your demo session has been cleared.",
	})
	a.deps.Redirect(w, r, a.deps.RouteURL("login", nil))
}
