package auth

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
)

type loginForm struct {
	Email    string
	Password string
	Remember bool
}

func newLoginForm(r *http.Request) loginForm {
	return loginForm{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password"),
		Remember: r.FormValue("remember") == "on" || r.FormValue("remember") == "true" || r.FormValue("remember") == "1",
	}
}

func (f loginForm) validate() httpx.ValidationErrors {
	errors := httpx.ValidationErrors{}

	if f.Email == "" {
		errors["email"] = "Enter the demo email address."
	}

	if f.Password == "" {
		errors["password"] = "Enter the demo password."
	}

	return errors
}
