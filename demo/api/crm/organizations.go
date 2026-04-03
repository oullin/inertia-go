package crm

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
)

func (a app) organizationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	pageNum := 1

	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			pageNum = n
		}
	}

	page, err := a.service.listOrganizationsPaginated(search, pageNum)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	a.deps.Render(w, r, "Organizations/Index", httpx.Props{
		"filters": map[string]any{
			"search": search,
		},
		"organizations": offsetOrganizationsProps(page),
	})
}

func (a app) organizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	organizationID, _, ok := routeIDAndAction(r.URL.Path, "/organizations/")

	if !ok {
		http.NotFound(w, r)

		return
	}

	switch r.Method {
	case http.MethodGet:
		a.showOrganizationHandler(w, r, organizationID)
	case http.MethodPost, http.MethodPut:
		a.updateOrganizationHandler(w, r, organizationID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) showOrganizationHandler(w http.ResponseWriter, r *http.Request, organizationID int64) {
	org, err := a.service.getOrganization(organizationID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if org == nil {
		http.NotFound(w, r)

		return
	}

	var cursor *string

	if c := r.URL.Query().Get("cursor"); c != "" {
		cursor = &c
	}

	contactsPage, err := a.service.listContactsByOrgPaginated(organizationID, cursor)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	a.deps.Render(w, r, "Organizations/Show", httpx.Props{
		"organization": organizationProp(*org),
		"contacts":     cursorContactsProps(contactsPage),
	})
}

func (a app) updateOrganizationHandler(w http.ResponseWriter, r *http.Request, organizationID int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	form := newOrganizationForm(r)
	errors := form.validate()

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		a.showOrganizationHandler(w, r.WithContext(ctx), organizationID)

		return
	}

	if err := a.service.updateOrganization(organizationID, form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	a.deps.SetFlash(w, flash.Message{
		Kind:    "success",
		Title:   "Organization updated",
		Message: "The company record was saved.",
	})
	a.deps.Redirect(w, r, a.deps.RouteURL("organizations.show", map[string]string{"organization": strconv.FormatInt(organizationID, 10)}))
}
