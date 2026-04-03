package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/props"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func registerCRMRoutes(mux *http.ServeMux) {
	mux.Handle("/dashboard", requireDemoAuth(http.HandlerFunc(crmDashboardHandler)))
	mux.Handle("/contacts", requireDemoAuth(http.HandlerFunc(contactsHandler)))
	mux.Handle("/contacts/create", requireDemoAuth(http.HandlerFunc(contactsCreateHandler)))
	mux.Handle("/contacts/", requireDemoAuth(http.HandlerFunc(contactByIDHandler)))
	mux.Handle("/organizations", requireDemoAuth(http.HandlerFunc(organizationsHandler)))
	mux.Handle("/organizations/", requireDemoAuth(http.HandlerFunc(organizationByIDHandler)))
}

func crmDashboardHandler(w http.ResponseWriter, r *http.Request) {
	activity, _ := database.ListRecentNotes(db, 5)

	renderPage(w, r, "Crm/Dashboard", httpx.Props{
		"recentActivity": recentActivityProps(activity),
		"totalContacts": props.Defer(func() any {
			time.Sleep(150 * time.Millisecond)

			return database.CountContacts(db)
		}, "stats"),
		"totalOrganizations": props.Defer(func() any {
			time.Sleep(150 * time.Millisecond)

			return database.CountOrganizations(db)
		}, "stats"),
		"recentNotesCount": props.Defer(func() any {
			time.Sleep(150 * time.Millisecond)

			return database.CountNotes(db)
		}, "stats"),
	})
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listContactsHandler(w, r)
	case http.MethodPost:
		storeContactHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func contactsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	renderPage(w, r, "Contacts/Create", httpx.Props{
		"form":          emptyContactProps(),
		"organizations": organizationOptions(),
	})
}

func contactByIDHandler(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/contacts/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)

		return
	}

	contactID, err := strconv.ParseInt(parts[0], 10, 64)

	if err != nil {
		http.NotFound(w, r)

		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			showContactHandler(w, r, contactID)
		case http.MethodPost, http.MethodPut:
			updateContactHandler(w, r, contactID)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	switch parts[1] {
	case "edit":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		editContactHandler(w, r, contactID)
	case "favorite":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		toggleFavoriteHandler(w, r, contactID)
	case "notes":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		storeNoteHandler(w, r, contactID)
	default:
		http.NotFound(w, r)
	}
}

func organizationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	orgs, err := database.ListOrganizations(db, search)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	renderPage(w, r, "Organizations/Index", httpx.Props{
		"filters": map[string]any{
			"search": search,
		},
		"organizations": organizationsProps(orgs),
	})
}

func organizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/organizations/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)

		return
	}

	organizationID, err := strconv.ParseInt(parts[0], 10, 64)

	if err != nil {
		http.NotFound(w, r)

		return
	}

	switch r.Method {
	case http.MethodGet:
		showOrganizationHandler(w, r, organizationID)
	case http.MethodPost, http.MethodPut:
		updateOrganizationHandler(w, r, organizationID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func listContactsHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	favoriteOnly := r.URL.Query().Get("favorite") == "true"
	contacts, err := database.ListContacts(db, search, favoriteOnly)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	renderPage(w, r, "Contacts/Index", httpx.Props{
		"filters": map[string]any{
			"search":   search,
			"favorite": favoriteOnly,
		},
		"contacts": contactsProps(contacts),
	})
}

func showContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	contact, err := database.GetContact(db, contactID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if contact == nil {
		http.NotFound(w, r)

		return
	}

	notes, err := database.ListContactNotes(db, contactID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	renderPage(w, r, "Contacts/Show", httpx.Props{
		"contact": contactProp(*contact),
		"notes":   notesProps(notes),
	})
}

func editContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	contact, err := database.GetContact(db, contactID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if contact == nil {
		http.NotFound(w, r)

		return
	}

	renderPage(w, r, "Contacts/Edit", httpx.Props{
		"contact":       contactProp(*contact),
		"form":          contactFormProps(*contact),
		"organizations": organizationOptions(),
	})
}

func storeContactHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	contact, errors := contactFromRequest(r)

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		renderPageWithContext(w, r.WithContext(ctx), "Contacts/Create", httpx.Props{
			"form":          contact,
			"organizations": organizationOptions(),
		})

		return
	}

	id, err := database.CreateContact(db, database.Contact{
		OrganizationID: parseOrganizationID(contact["organization_id"]),
		FirstName:      toString(contact["first_name"]),
		LastName:       toString(contact["last_name"]),
		Email:          toString(contact["email"]),
		Phone:          toString(contact["phone"]),
		Address:        toString(contact["address"]),
		City:           toString(contact["city"]),
		Region:         toString(contact["region"]),
		Country:        toString(contact["country"]),
		PostalCode:     toString(contact["postal_code"]),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Contact created",
		Message: "The CRM record is ready for follow-up.",
	})

	i.Redirect(w, r, routeURL("contacts.show", map[string]string{"contact": strconv.FormatInt(id, 10)}))
}

func updateContactHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	contact, errors := contactFromRequest(r)

	if len(errors) > 0 {
		existing, _ := database.GetContact(db, contactID)
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		renderPageWithContext(w, r.WithContext(ctx), "Contacts/Edit", httpx.Props{
			"contact":       contactPropValue(existing),
			"form":          contact,
			"organizations": organizationOptions(),
		})

		return
	}

	err := database.UpdateContact(db, contactID, database.Contact{
		OrganizationID: parseOrganizationID(contact["organization_id"]),
		FirstName:      toString(contact["first_name"]),
		LastName:       toString(contact["last_name"]),
		Email:          toString(contact["email"]),
		Phone:          toString(contact["phone"]),
		Address:        toString(contact["address"]),
		City:           toString(contact["city"]),
		Region:         toString(contact["region"]),
		Country:        toString(contact["country"]),
		PostalCode:     toString(contact["postal_code"]),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Contact updated",
		Message: "The CRM record now reflects the latest details.",
	})
	i.Redirect(w, r, routeURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func toggleFavoriteHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := database.ToggleContactFavorite(db, contactID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Favorite updated",
		Message: "The contact pin state changed successfully.",
	})
	i.Redirect(w, r, routeURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func storeNoteHandler(w http.ResponseWriter, r *http.Request, contactID int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	body := strings.TrimSpace(r.FormValue("body"))

	if body == "" {
		ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{
			"body": "Add note content before saving.",
		})
		showContactHandler(w, r.WithContext(ctx), contactID)

		return
	}

	user := currentUser(r)

	if user == nil {
		i.Redirect(w, r, routeURL("login", nil))

		return
	}

	if _, err := database.CreateNote(db, contactID, user.ID, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Note added",
		Message: "Recent activity has been updated.",
	})
	i.Redirect(w, r, routeURL("contacts.show", map[string]string{"contact": strconv.FormatInt(contactID, 10)}))
}

func showOrganizationHandler(w http.ResponseWriter, r *http.Request, organizationID int64) {
	org, err := database.GetOrganization(db, organizationID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if org == nil {
		http.NotFound(w, r)

		return
	}

	contacts, err := database.ListContactsByOrganization(db, organizationID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	renderPage(w, r, "Organizations/Show", httpx.Props{
		"organization": organizationProp(*org),
		"contacts":     contactsProps(contacts),
	})
}

func updateOrganizationHandler(w http.ResponseWriter, r *http.Request, organizationID int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	name := strings.TrimSpace(r.FormValue("name"))

	if name == "" {
		ctx := inertia.SetValidationErrors(r.Context(), httpx.ValidationErrors{
			"name": "Organization name is required.",
		})
		showOrganizationHandler(w, r.WithContext(ctx), organizationID)

		return
	}

	if err := database.UpdateOrganization(db, organizationID, name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Organization updated",
		Message: "The company record was saved.",
	})
	i.Redirect(w, r, routeURL("organizations.show", map[string]string{"organization": strconv.FormatInt(organizationID, 10)}))
}

func contactsProps(contacts []database.Contact) []map[string]any {
	items := make([]map[string]any, 0, len(contacts))

	for _, contact := range contacts {
		items = append(items, contactProp(contact))
	}

	return items
}

func contactProp(contact database.Contact) map[string]any {
	item := map[string]any{
		"id":          contact.ID,
		"first_name":  contact.FirstName,
		"last_name":   contact.LastName,
		"email":       contact.Email,
		"phone":       contact.Phone,
		"address":     contact.Address,
		"city":        contact.City,
		"region":      contact.Region,
		"country":     contact.Country,
		"postal_code": contact.PostalCode,
		"is_favorite": contact.IsFavorite,
		"created_at":  contact.CreatedAt.Format(time.RFC3339),
		"updated_at":  contact.UpdatedAt.Format(time.RFC3339),
	}

	if contact.OrganizationID != nil {
		item["organization"] = map[string]any{
			"id":   *contact.OrganizationID,
			"name": contact.OrganizationName,
		}
	} else {
		item["organization"] = nil
	}

	return item
}

func contactPropValue(contact *database.Contact) map[string]any {
	if contact == nil {
		return nil
	}

	return contactProp(*contact)
}

func contactFormProps(contact database.Contact) map[string]any {
	values := emptyContactProps()
	values["first_name"] = contact.FirstName
	values["last_name"] = contact.LastName
	values["email"] = contact.Email
	values["phone"] = contact.Phone
	values["address"] = contact.Address
	values["city"] = contact.City
	values["region"] = contact.Region
	values["country"] = contact.Country
	values["postal_code"] = contact.PostalCode

	if contact.OrganizationID != nil {
		values["organization_id"] = fmt.Sprintf("%d", *contact.OrganizationID)
	}

	return values
}

func emptyContactProps() map[string]any {
	return map[string]any{
		"organization_id": "",
		"first_name":      "",
		"last_name":       "",
		"email":           "",
		"phone":           "",
		"address":         "",
		"city":            "",
		"region":          "",
		"country":         "",
		"postal_code":     "",
	}
}

func organizationsProps(orgs []database.Organization) []map[string]any {
	items := make([]map[string]any, 0, len(orgs))

	for _, org := range orgs {
		items = append(items, organizationProp(org))
	}

	return items
}

func organizationProp(org database.Organization) map[string]any {
	return map[string]any{
		"id":             org.ID,
		"name":           org.Name,
		"contacts_count": org.ContactsCount,
	}
}

func organizationOptions() []map[string]any {
	orgs, _ := database.ListOrganizations(db, "")
	options := []map[string]any{{"value": "", "label": "No organization"}}

	for _, org := range orgs {
		options = append(options, map[string]any{
			"value": fmt.Sprintf("%d", org.ID),
			"label": org.Name,
		})
	}

	return options
}

func notesProps(notes []database.Note) []map[string]any {
	items := make([]map[string]any, 0, len(notes))

	for _, note := range notes {
		items = append(items, map[string]any{
			"id":         note.ID,
			"body":       note.Body,
			"created_at": note.CreatedAt.Format(time.RFC3339),
			"user": map[string]any{
				"id":   note.UserID,
				"name": note.UserName,
			},
			"contact": map[string]any{
				"id":   note.ContactID,
				"name": note.ContactName,
			},
		})
	}

	return items
}

func recentActivityProps(notes []database.Note) []map[string]any {
	return notesProps(notes)
}

func contactFromRequest(r *http.Request) (map[string]any, httpx.ValidationErrors) {
	values := map[string]any{
		"organization_id": strings.TrimSpace(r.FormValue("organization_id")),
		"first_name":      strings.TrimSpace(r.FormValue("first_name")),
		"last_name":       strings.TrimSpace(r.FormValue("last_name")),
		"email":           strings.TrimSpace(r.FormValue("email")),
		"phone":           strings.TrimSpace(r.FormValue("phone")),
		"address":         strings.TrimSpace(r.FormValue("address")),
		"city":            strings.TrimSpace(r.FormValue("city")),
		"region":          strings.TrimSpace(r.FormValue("region")),
		"country":         strings.TrimSpace(r.FormValue("country")),
		"postal_code":     strings.TrimSpace(r.FormValue("postal_code")),
	}

	errors := httpx.ValidationErrors{}

	if values["first_name"] == "" {
		errors["first_name"] = "First name is required."
	}

	if values["last_name"] == "" {
		errors["last_name"] = "Last name is required."
	}

	if values["email"] == "" || !strings.Contains(toString(values["email"]), "@") {
		errors["email"] = "A valid email address is required."
	}

	return values, errors
}

func parseOrganizationID(v any) *int64 {
	raw := strings.TrimSpace(toString(v))

	if raw == "" {
		return nil
	}

	id, err := strconv.ParseInt(raw, 10, 64)

	if err != nil {
		return nil
	}

	return &id
}

func toString(v any) string {
	s, _ := v.(string)

	return s
}
