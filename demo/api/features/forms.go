package features

import (
	"fmt"
	"log/slog"
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
		a.container.Render(w, r, "Features/Forms/UseForm", httpx.Props{})
	case http.MethodPost:
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))

		errors := httpx.ValidationErrors{}

		if strings.TrimSpace(name) == "" {
			errors["name"] = "Name is required."
		}

		if strings.TrimSpace(email) == "" {
			errors["email"] = "Email is required."
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)

			a.container.Render(w, r.WithContext(ctx), "Features/Forms/UseForm", httpx.Props{})

			return
		}

		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Form submitted", Message: fmt.Sprintf("Hello, %s!", name)}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.forms.use-form", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) formComponentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Forms/FormComponent", httpx.Props{})
	case http.MethodPost:
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		name := strings.TrimSpace(r.FormValue("name"))

		errors := httpx.ValidationErrors{}

		if strings.TrimSpace(name) == "" {
			errors["name"] = "Name is required."
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)

			a.container.Render(w, r.WithContext(ctx), "Features/Forms/FormComponent", httpx.Props{})

			return
		}

		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Form submitted", Message: fmt.Sprintf("Hello, %s!", name)}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.forms.form-component", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) fileUploadsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Forms/FileUploads", httpx.Props{})
	case http.MethodPost:
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		if r.MultipartForm == nil {
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		files := r.MultipartForm.File["files"]
		count := len(files)

		if _, _, err := r.FormFile("photo"); err == nil {
			count++
		}

		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Files uploaded", Message: fmt.Sprintf("%d file(s) received.", count)}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.forms.file-uploads", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) validationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Forms/Validation", httpx.Props{})
	case http.MethodPost:
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		errors := httpx.ValidationErrors{}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		age := strings.TrimSpace(r.FormValue("age"))

		if strings.TrimSpace(name) == "" {
			errors["name"] = "Name is required."
		}

		if strings.TrimSpace(email) == "" {
			errors["email"] = "That doesn't look like a valid email."
		}

		if strings.TrimSpace(age) != "" {
			if n, err := strconv.Atoi(age); err != nil || n < 18 || n > 120 {
				errors["age"] = "You must be at least 18 years old."
			}
		}

		if len(errors) > 0 {
			ctx := inertia.SetValidationErrors(r.Context(), errors)

			a.container.Render(w, r.WithContext(ctx), "Features/Forms/Validation", httpx.Props{})

			return
		}

		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Valid", Message: "All fields passed validation."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.forms.validation", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) validationSecondaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	if err := httpx.ParseForm(r); err != nil {
		slog.Error("parse form", "error", err)

		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	errors := httpx.ValidationErrors{}
	feedback := strings.TrimSpace(r.FormValue("feedback"))

	if strings.TrimSpace(feedback) == "" {
		errors["feedback"] = "Feedback is required."
	}

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)

		a.container.Render(w, r.WithContext(ctx), "Features/Forms/Validation", httpx.Props{})

		return
	}

	if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Secondary form", Message: "Feedback submitted."}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.container.Redirect(w, r, a.container.RouteURL("features.forms.validation", nil))
}

func (a app) precognitionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Forms/Precognition", httpx.Props{})
	case http.MethodPost:
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		errors := httpx.ValidationErrors{}

		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirmation")

		if len(username) < 3 || len(username) > 20 {
			errors["username"] = "Username must be 3-20 characters."
		}

		if strings.TrimSpace(email) == "" || !strings.Contains(email, "@") {
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

			a.container.Render(w, r.WithContext(ctx), "Features/Forms/Precognition", httpx.Props{})

			return
		}

		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Account created", Message: fmt.Sprintf("Welcome, %s!", username)}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.forms.precognition", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) optimisticUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	contacts, err := database.ListContacts(a.container.DB, "", false)

	if err != nil {
		slog.Error("list contacts", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

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

	a.container.Render(w, r, "Features/Forms/OptimisticUpdates", httpx.Props{"contacts": items})
}

func (a app) optimisticToggleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		http.NotFound(w, r)

		return
	}

	err = database.ToggleContactFavorite(a.container.DB, id)

	if err != nil {
		if err := a.container.SetFlash(w, flash.Message{Kind: "error", Title: "Error", Message: "Unable to toggle favorite."}); err != nil {
			slog.Error("flash: set", "error", err)
		}
	} else {
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Favorite updated", Message: "Toggle applied."}); err != nil {
			slog.Error("flash: set", "error", err)
		}
	}

	a.container.Redirect(w, r, a.container.RouteURL("features.forms.optimistic-updates", nil))
}

func (a app) formContextHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Forms/UseFormContext", httpx.Props{})
}

func (a app) dottedKeysHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Forms/DottedKeys", httpx.Props{})
	case http.MethodPost:
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("parse form", "error", err)

			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}

		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Dotted keys", Message: "Nested form data received."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.forms.dotted-keys", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) wayfinderHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Forms/Wayfinder", httpx.Props{})
}
