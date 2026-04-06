package crm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func TestRouteIDAndAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path   string
		prefix string
		id     int64
		action string
		ok     bool
	}{
		{path: "/contacts/12", prefix: "/contacts/", id: 12, action: "", ok: true},
		{path: "/contacts/12/edit", prefix: "/contacts/", id: 12, action: "edit", ok: true},
		{path: "/organizations/2", prefix: "/organizations/", id: 2, action: "", ok: true},
		{path: "/contacts", prefix: "/contacts/", ok: false},
		{path: "/contacts/abc", prefix: "/contacts/", ok: false},
		{path: "/contacts/12/edit/extra", prefix: "/contacts/", ok: false},
	}

	for _, tt := range tests {
		id, action, ok := routeIDAndAction(tt.path, tt.prefix)

		if id != tt.id || action != tt.action || ok != tt.ok {
			t.Fatalf("routeIDAndAction(%q, %q) = (%d, %q, %v), want (%d, %q, %v)", tt.path, tt.prefix, id, action, ok, tt.id, tt.action, tt.ok)
		}
	}
}

func TestDashboardAndListHandlers(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)
	req := h.request(http.MethodGet, "/dashboard")
	w := httptest.NewRecorder()

	h.app.dashboardHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Crm/Dashboard")

	page.AssertHasProp(t, "recentActivity")

	req = h.partialRequest("/dashboard", "Crm/Dashboard", "totalContacts,totalOrganizations,recentNotesCount")
	w = httptest.NewRecorder()

	h.app.dashboardHandler(w, req)

	page = h.page(t, w)

	page.AssertHasProp(t, "totalContacts")
	page.AssertHasProp(t, "totalOrganizations")
	page.AssertHasProp(t, "recentNotesCount")

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	req = h.partialRequest("/dashboard", "Crm/Dashboard", "totalContacts,totalOrganizations,recentNotesCount")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	h.app.dashboardHandler(w, req)

	page = h.page(t, w)

	if page.Props["totalContacts"] != nil || page.Props["totalOrganizations"] != nil || page.Props["recentNotesCount"] != nil {
		t.Fatalf("cancelled deferred props = %#v, want nil values", page.Props)
	}

	req = h.request(http.MethodGet, "/contacts?search=Northwind&favorite=true")
	w = httptest.NewRecorder()

	h.app.contactsHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Contacts/Index")

	page.AssertHasProp(t, "filters")
	page.AssertHasProp(t, "contacts")

	req = h.request(http.MethodGet, "/contacts/create?organization_id=1")
	w = httptest.NewRecorder()

	h.app.contactsCreateHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Contacts/Create")

	page.AssertHasProp(t, "organizations")

	req = h.request(http.MethodGet, "/organizations?page=2&search=North")
	w = httptest.NewRecorder()

	h.app.organizationsHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Organizations/Index")

	page.AssertHasProp(t, "organizations")
}

func TestContactByIDDispatchAndMutations(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)

	req := h.request(http.MethodGet, "/contacts/1")
	w := httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Contacts/Show")

	req = h.partialRequest("/contacts/1", "Contacts/Show", "notes")
	w = httptest.NewRecorder()

	h.app.showContactHandler(w, req, 1)

	page = h.page(t, w)

	page.AssertHasProp(t, "notes")

	req = h.request(http.MethodGet, "/contacts/1/edit")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Contacts/Edit")

	req = h.request(http.MethodPost, "/contacts/1/favorite")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/contacts/1" {
		t.Fatalf("favorite status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = h.request(http.MethodDelete, "/contacts/1")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/contacts" {
		t.Fatalf("delete status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = h.request(http.MethodGet, "/contacts/unknown")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodGet, "/contacts/999/favorite")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestStoreAndUpdateContactBranches(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)

	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(url.Values{
		"organization_id": {"bad"},
		"first_name":      {""},
		"last_name":       {""},
		"email":           {"bad"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w := httptest.NewRecorder()

	h.app.storeContactHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Contacts/Create")

	page.AssertHasProp(t, "errors")

	req = httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Mina"},
		"last_name":       {"Cole"},
		"email":           {"mina@example.com"},
		"phone":           {"555-0110"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.storeContactHandler(w, req)

	if w.Code != http.StatusFound || strings.TrimSpace(w.Header().Get("Location")) == "" {
		t.Fatalf("store status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts/1", strings.NewReader(url.Values{
		"first_name": {" "},
		"last_name":  {" "},
		"email":      {"bad"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateContactHandler(w, req, 1)

	page = h.page(t, w)

	page.AssertComponent(t, "Contacts/Edit")

	req = httptest.NewRequest(http.MethodPost, "/contacts/1", strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Updated"},
		"last_name":       {"Name"},
		"email":           {"updated@example.com"},
		"phone":           {"555-0111"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateContactHandler(w, req, 1)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/contacts/1" {
		t.Fatalf("update status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts/999", strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Missing"},
		"last_name":       {"Contact"},
		"email":           {"missing@example.com"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateContactHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("missing update status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts/999", strings.NewReader(url.Values{
		"first_name": {""},
		"last_name":  {""},
		"email":      {"bad"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateContactHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("missing invalid update status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader("bad"))

	req.Header.Set("Content-Type", "multipart/form-data; boundary=bad")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.storeContactHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("store bad request status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts/1", strings.NewReader("bad"))

	req.Header.Set("Content-Type", "multipart/form-data; boundary=bad")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateContactHandler(w, req, 1)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("update bad request status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestStoreNoteAndOrganizationBranches(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)

	req := httptest.NewRequest(http.MethodPost, "/contacts/1/notes", strings.NewReader(url.Values{
		"body": {""},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w := httptest.NewRecorder()

	h.app.storeNoteHandler(w, req, 1)

	page := h.page(t, w)

	page.AssertComponent(t, "Contacts/Show")

	h.user = nil
	req = httptest.NewRequest(http.MethodPost, "/contacts/1/notes", strings.NewReader(url.Values{
		"body": {"Need follow-up"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.storeNoteHandler(w, req, 1)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/login" {
		t.Fatalf("note redirect status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	user, err := database.FindUserByEmail(h.db, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	h.user = user
	req = httptest.NewRequest(http.MethodPost, "/contacts/1/notes", strings.NewReader(url.Values{
		"body": {"Need follow-up"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.storeNoteHandler(w, req, 1)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/contacts/1" {
		t.Fatalf("note status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = h.request(http.MethodGet, "/organizations/1")
	w = httptest.NewRecorder()

	h.app.organizationByIDHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Organizations/Show")

	req = httptest.NewRequest(http.MethodPost, "/organizations/1", strings.NewReader(url.Values{
		"name": {""},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateOrganizationHandler(w, req, 1)

	page = h.page(t, w)

	page.AssertComponent(t, "Organizations/Show")

	req = httptest.NewRequest(http.MethodPost, "/organizations/1", strings.NewReader(url.Values{
		"name": {"Updated Org"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateOrganizationHandler(w, req, 1)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/organizations/1" {
		t.Fatalf("organization update status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}
}

func TestCRMHandlers_MethodGuardsNotFoundAndDatabaseErrors(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)

	req := h.request(http.MethodPut, "/contacts")
	w := httptest.NewRecorder()

	h.app.contactsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("contactsHandler status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = h.request(http.MethodPost, "/contacts/create")
	w = httptest.NewRecorder()

	h.app.contactsCreateHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("contactsCreateHandler status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = h.request(http.MethodGet, "/contacts/999")
	w = httptest.NewRecorder()

	h.app.showContactHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("showContactHandler status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodGet, "/contacts/999/edit")
	w = httptest.NewRecorder()

	h.app.editContactHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("editContactHandler status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodDelete, "/contacts/999")
	w = httptest.NewRecorder()

	h.app.deleteContactHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("deleteContactHandler status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodPost, "/contacts/999/favorite")
	w = httptest.NewRecorder()

	h.app.toggleFavoriteHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("toggleFavoriteHandler status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodDelete, "/organizations")
	w = httptest.NewRecorder()

	h.app.organizationsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("organizationsHandler status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = h.request(http.MethodGet, "/organizations/bad")
	w = httptest.NewRecorder()

	h.app.organizationByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("organizationByIDHandler invalid path status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodGet, "/organizations/999")
	w = httptest.NewRecorder()

	h.app.showOrganizationHandler(w, req, 999)

	if w.Code != http.StatusNotFound {
		t.Fatalf("showOrganizationHandler status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodPost, "/contacts/1/notes")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("contactByIDHandler note status = %d, want %d", w.Code, http.StatusOK)
	}

	req = h.request(http.MethodGet, "/contacts/1/notes")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("contactByIDHandler wrong note method status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = h.request(http.MethodGet, "/contacts/1/missing")
	w = httptest.NewRecorder()

	h.app.contactByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("contactByIDHandler unknown action status = %d, want %d", w.Code, http.StatusNotFound)
	}

	req = h.request(http.MethodDelete, "/organizations/1")
	w = httptest.NewRecorder()

	h.app.organizationByIDHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("organizationByIDHandler status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	h.db.Close()

	req = h.request(http.MethodGet, "/dashboard")
	w = httptest.NewRecorder()

	h.app.dashboardHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("dashboardHandler status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/contacts")
	w = httptest.NewRecorder()

	h.app.listContactsHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("listContactsHandler status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/contacts/create")
	w = httptest.NewRecorder()

	h.app.contactsCreateHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("contactsCreateHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/contacts/1")
	w = httptest.NewRecorder()

	h.app.showContactHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("showContactHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/contacts/1/edit")
	w = httptest.NewRecorder()

	h.app.editContactHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("editContactHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Closed"},
		"last_name":       {"DB"},
		"email":           {"closed@example.com"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.storeContactHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("storeContactHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts/1", strings.NewReader(url.Values{
		"organization_id": {"1"},
		"first_name":      {"Closed"},
		"last_name":       {"DB"},
		"email":           {"closed@example.com"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateContactHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("updateContactHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodPost, "/contacts/1/favorite")
	w = httptest.NewRecorder()

	h.app.toggleFavoriteHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("toggleFavoriteHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = httptest.NewRequest(http.MethodPost, "/contacts/1/notes", strings.NewReader(url.Values{
		"body": {"Need note"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.storeNoteHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("storeNoteHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/organizations")
	w = httptest.NewRecorder()

	h.app.organizationsHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("organizationsHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/organizations/1")
	w = httptest.NewRecorder()

	h.app.showOrganizationHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("showOrganizationHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = httptest.NewRequest(http.MethodPost, "/organizations/1", strings.NewReader(url.Values{
		"name": {"Closed DB"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.updateOrganizationHandler(w, req, 1)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("updateOrganizationHandler db error status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
