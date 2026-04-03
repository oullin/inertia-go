package main

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func registerAuthRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/login", guestOnly(http.HandlerFunc(loginHandler)))
	mux.Handle("/logout", requireDemoAuth(http.HandlerFunc(logoutHandler)))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if currentUser(r) != nil {
		i.Redirect(w, r, routeURL("dashboard", nil))

		return
	}

	i.Redirect(w, r, routeURL("login", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderPage(w, r, "Auth/Login", httpx.Props{
			"status": r.URL.Query().Get("status"),
		})
	case http.MethodPost:
		loginSubmitHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on" || r.FormValue("remember") == "true" || r.FormValue("remember") == "1"
	errors := httpx.ValidationErrors{}

	if email == "" {
		errors["email"] = "Enter the demo email address."
	}

	if password == "" {
		errors["password"] = "Enter the demo password."
	}

	user, err := database.FindUserByEmail(db, email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if user == nil || user.Password != password {
		errors["email"] = "Use test@example.com and password to sign in."
	}

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		renderPageWithContext(w, r.WithContext(ctx), "Auth/Login", httpx.Props{})

		return
	}

	setDemoSession(w, user.ID, remember)

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Signed in",
		Message: "The demo session is now authenticated.",
	})

	i.Redirect(w, r, routeURL("dashboard", nil))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	clearDemoSession(w)
	setFlash(w, flashPayload{
		Kind:    "info",
		Title:   "Signed out",
		Message: "Your demo session has been cleared.",
	})
	i.Redirect(w, r, routeURL("login", nil))
}
