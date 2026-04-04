package crm

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/props"
)

func (a app) contactsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.listContactsHandler(w, r)
	case http.MethodPost:
		a.storeContactHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) contactsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	orgs, err := a.service.listOrganizations("")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	form := emptyContactForm()

	if orgID := r.URL.Query().Get("organization_id"); orgID != "" {
		form.OrganizationID = orgID
	}

	a.deps.Render(w, r, "Contacts/Create", httpx.Props{
		"form":          contactFormProps(form),
		"organizations": organizationOptions(orgs),
	})
}

func (a app) contactByIDHandler(w http.ResponseWriter, r *http.Request) {
	contactID, action, ok := routeIDAndAction(r.URL.Path, "/contacts/")

	if !ok {
		http.NotFound(w, r)

		return
	}

	if action == "" {
		switch r.Method {
		case http.MethodGet:
			a.showContactHandler(w, r, contactID)
		case http.MethodPost, http.MethodPut:
			a.updateContactHandler(w, r, contactID)
		case http.MethodDelete:
			a.deleteContactHandler(w, r, contactID)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	switch action {
	case "edit":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		a.editContactHandler(w, r, contactID)
	case "favorite":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		a.toggleFavoriteHandler(w, r, contactID)
	case "notes":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		a.storeNoteHandler(w, r, contactID)
	default:
		http.NotFound(w, r)
	}
}

func (a app) listContactsHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	favoriteOnly := r.URL.Query().Get("favorite") == "true"

	var cursor *string

	if c := r.URL.Query().Get("cursor"); c != "" {
		cursor = &c
	}

	direction := r.URL.Query().Get("direction")

	page, err := a.service.listContactsPaginated(search, favoriteOnly, cursor, direction)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	a.deps.Render(w, r, "Contacts/Index", httpx.Props{
		"filters": map[string]any{
			"search":   search,
			"favorite": favoriteOnly,
		},
		"contacts": cursorContactsProps(page),
	})
}

func (a app) showContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	contact, err := a.service.getContact(contactID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if contact == nil {
		http.NotFound(w, r)

		return
	}

	a.deps.Render(w, r, "Contacts/Show", httpx.Props{
		"contact": contactProp(*contact),
		"notes": props.Defer(func() any {
			notes, err := a.service.listContactNotes(contactID)

			if err != nil {
				slog.Error("list contact notes", "error", err)
			}

			return notesProps(notes)
		}),
	})
}

func (a app) deleteContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := a.service.deleteContact(contactID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := a.deps.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Contact deleted",
		Message: "The contact has been removed.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("contacts.index", nil))
}

func (a app) editContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	contact, err := a.service.getContact(contactID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if contact == nil {
		http.NotFound(w, r)

		return
	}

	orgs, err := a.service.listOrganizations("")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	a.deps.Render(w, r, "Contacts/Edit", httpx.Props{
		"contact":       contactProp(*contact),
		"form":          contactFormProps(newContactFormFromContact(*contact)),
		"organizations": organizationOptions(orgs),
	})
}

func (a app) storeContactHandler(w http.ResponseWriter, r *http.Request) {
	if err := httpx.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	form := newContactForm(r)
	errors := form.validate()

	if len(errors) > 0 {
		orgs, err := a.service.listOrganizations("")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		ctx := inertia.SetValidationErrors(r.Context(), errors)
		a.deps.Render(w, r.WithContext(ctx), "Contacts/Create", httpx.Props{
			"form":          contactFormProps(form),
			"organizations": organizationOptions(orgs),
		})

		return
	}

	id, err := a.service.createContact(form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := a.deps.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Contact created",
		Message: "The CRM record is ready for follow-up.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(id, 10)}))
}

func (a app) updateContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := httpx.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	form := newContactForm(r)
	errors := form.validate()

	if len(errors) > 0 {
		existing, err := a.service.getContact(contactID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		orgs, err := a.service.listOrganizations("")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		ctx := inertia.SetValidationErrors(r.Context(), errors)
		a.deps.Render(w, r.WithContext(ctx), "Contacts/Edit", httpx.Props{
			"contact":       contactPropValue(existing),
			"form":          contactFormProps(form),
			"organizations": organizationOptions(orgs),
		})

		return
	}

	if err := a.service.updateContact(contactID, form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := a.deps.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Contact updated",
		Message: "The CRM record now reflects the latest details.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func (a app) toggleFavoriteHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := a.service.toggleFavorite(contactID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := a.deps.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Favorite updated",
		Message: "The contact pin state changed successfully.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func (a app) storeNoteHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := httpx.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	body := strings.TrimSpace(r.FormValue("body"))

	if body == "" {
		ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{
			"body": "Add note content before saving.",
		})
		a.showContactHandler(w, r.WithContext(ctx), contactID)

		return
	}

	err := a.service.createNote(contactID, a.deps.CurrentUser(r), body)

	if errors.Is(err, errUnauthorized) {
		a.deps.Redirect(w, r, a.deps.RouteURL("login", nil))

		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := a.deps.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Note added",
		Message: "Recent activity has been updated.",
	}); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func routeIDAndAction(path, prefix string) (int64, string, bool) {
	trimmed := strings.TrimPrefix(path, prefix)
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		return 0, "", false
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)

	if err != nil {
		return 0, "", false
	}

	if len(parts) == 1 {
		return id, "", true
	}

	return id, parts[1], true
}
