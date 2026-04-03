package features

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func (a app) useFormHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Forms/UseForm", httpx.Props{})
	case http.MethodPost:
		httpx.ParseForm(r)
		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))

		errors := httpx.ValidationErrors{}

		if name == "" {
			errors["name"] = "Name is required."
		}

		if email == "" {
			errors["email"] = "Email is required."
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)
			a.deps.RenderWithContext(w, r.WithContext(ctx), "Features/Forms/UseForm", httpx.Props{})

			return
		}

		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Form submitted", Message: fmt.Sprintf("Hello, %s!", name)})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.use-form", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) formComponentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Forms/FormComponent", httpx.Props{})
	case http.MethodPost:
		httpx.ParseForm(r)
		name := strings.TrimSpace(r.FormValue("name"))

		errors := httpx.ValidationErrors{}

		if name == "" {
			errors["name"] = "Name is required."
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)
			a.deps.RenderWithContext(w, r.WithContext(ctx), "Features/Forms/FormComponent", httpx.Props{})

			return
		}

		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Form submitted", Message: fmt.Sprintf("Hello, %s!", name)})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.form-component", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) fileUploadsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Forms/FileUploads", httpx.Props{})
	case http.MethodPost:
		httpx.ParseForm(r)
		files := r.MultipartForm.File["files"]
		count := len(files)

		if _, _, err := r.FormFile("photo"); err == nil {
			count++
		}

		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Files uploaded", Message: fmt.Sprintf("%d file(s) received.", count)})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.file-uploads", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) validationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Forms/Validation", httpx.Props{})
	case http.MethodPost:
		httpx.ParseForm(r)
		errors := httpx.ValidationErrors{}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		age := strings.TrimSpace(r.FormValue("age"))

		if name == "" {
			errors["name"] = "Name is required."
		}

		if email == "" {
			errors["email"] = "That doesn't look like a valid email."
		}

		if age != "" {
			if n, err := strconv.Atoi(age); err != nil || n < 18 || n > 120 {
				errors["age"] = "You must be at least 18 years old."
			}
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)
			a.deps.RenderWithContext(w, r.WithContext(ctx), "Features/Forms/Validation", httpx.Props{})

			return
		}

		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Valid", Message: "All fields passed validation."})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.validation", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) validationSecondaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	httpx.ParseForm(r)
	errors := httpx.ValidationErrors{}
	feedback := strings.TrimSpace(r.FormValue("feedback"))

	if feedback == "" {
		errors["feedback"] = "Feedback is required."
	}

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		a.deps.RenderWithContext(w, r.WithContext(ctx), "Features/Forms/Validation", httpx.Props{})

		return
	}

	a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Secondary form", Message: "Feedback submitted."})
	a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.validation", nil))
}

func (a app) precognitionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Forms/Precognition", httpx.Props{})
	case http.MethodPost:
		httpx.ParseForm(r)
		errors := httpx.ValidationErrors{}

		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirmation")

		if len(username) < 3 || len(username) > 20 {
			errors["username"] = "Username must be 3-20 characters."
		}

		if email == "" || !strings.Contains(email, "@") {
			errors["email"] = "A valid email is required."
		}

		if len(password) < 8 {
			errors["password"] = "Password must be at least 8 characters."
		}

		if password != passwordConfirm {
			errors["password_confirmation"] = "Passwords do not match."
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)
			a.deps.RenderWithContext(w, r.WithContext(ctx), "Features/Forms/Precognition", httpx.Props{})

			return
		}

		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Account created", Message: fmt.Sprintf("Welcome, %s!", username)})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.precognition", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) optimisticUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	contacts, _ := database.ListContacts(a.deps.DB, "", false)
	items := make([]map[string]any, 0, 10)

	for i, c := range contacts {
		if i >= 10 {
			break
		}

		items = append(items, map[string]any{
			"id":          c.ID,
			"first_name":  c.FirstName,
			"last_name":   c.LastName,
			"email":       c.Email,
			"is_favorite": c.IsFavorite,
		})
	}

	a.deps.Render(w, r, "Features/Forms/OptimisticUpdates", httpx.Props{"contacts": items})
}

func (a app) optimisticToggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/features/forms/optimistic-toggle/")
	id, err := strconv.ParseInt(strings.Trim(path, "/"), 10, 64)

	if err != nil {
		http.NotFound(w, r)

		return
	}

	database.ToggleContactFavorite(a.deps.DB, id)
	a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Favorite updated", Message: "Toggle applied."})
	a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.optimistic-updates", nil))
}

func (a app) formContextHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Forms/UseFormContext", httpx.Props{})
}

func (a app) dottedKeysHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Forms/DottedKeys", httpx.Props{})
	case http.MethodPost:
		httpx.ParseForm(r)
		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Dotted keys", Message: "Nested form data received."})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.forms.dotted-keys", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) wayfinderHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Forms/Wayfinder", httpx.Props{})
}
