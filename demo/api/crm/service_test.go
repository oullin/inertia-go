package crm

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/seed"
)

type serviceTestHarness struct {
	db  *sql.DB
	svc service
}

func newServiceTestHarness(t *testing.T) serviceTestHarness {
	t.Helper()

	db, err := database.Open(":memory:")

	if err != nil {
		t.Fatalf("database.Open() error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := seed.Run(db); err != nil {
		t.Fatalf("seed.Run() error = %v", err)
	}

	return serviceTestHarness{
		db:  db,
		svc: newService(newRepository(db)),
	}
}

func TestServiceCreateContact(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	id, err := h.svc.createContact(contactForm{
		OrganizationID: "2",
		FirstName:      " Mina ",
		LastName:       " Cole ",
		Email:          " mina@example.test ",
		Phone:          " +65 1234 ",
	})

	if err != nil {
		t.Fatalf("createContact() error = %v", err)
	}

	contact, err := database.GetContact(h.db, id)

	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}

	if contact == nil {
		t.Fatal("GetContact() = nil, want contact")
	}

	if contact.OrganizationID == nil || *contact.OrganizationID != 2 {
		t.Fatalf("organization_id = %v, want 2", contact.OrganizationID)
	}

	if contact.FirstName != "Mina" || contact.LastName != "Cole" || contact.Email != "mina@example.test" {
		t.Fatalf("contact = %#v", *contact)
	}

	if contact.Phone != "+65 1234" {
		t.Fatalf("phone = %q, want %q", contact.Phone, "+65 1234")
	}
}

func TestServiceUpdateContact(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	err := h.svc.updateContact(1, contactForm{
		OrganizationID: "3",
		FirstName:      " Mina ",
		LastName:       " Cole ",
		Email:          " mina.updated@example.test ",
		Phone:          " +65 9999 ",
	})

	if err != nil {
		t.Fatalf("updateContact() error = %v", err)
	}

	contact, err := database.GetContact(h.db, 1)

	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}

	if contact == nil {
		t.Fatal("GetContact() = nil, want contact")
	}

	if contact.OrganizationID == nil || *contact.OrganizationID != 3 {
		t.Fatalf("organization_id = %v, want 3", contact.OrganizationID)
	}

	if contact.FirstName != "Mina" || contact.LastName != "Cole" || contact.Email != "mina.updated@example.test" {
		t.Fatalf("updated contact = %#v", *contact)
	}

	if contact.Phone != "+65 9999" {
		t.Fatalf("phone = %q, want %q", contact.Phone, "+65 9999")
	}
}

func TestServiceToggleFavorite(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	before, err := database.GetContact(h.db, 2)

	if err != nil {
		t.Fatalf("GetContact() before toggle error = %v", err)
	}

	if before == nil {
		t.Fatal("GetContact() before toggle = nil, want contact")
	}

	if err := h.svc.toggleFavorite(2); err != nil {
		t.Fatalf("toggleFavorite() error = %v", err)
	}

	after, err := database.GetContact(h.db, 2)

	if err != nil {
		t.Fatalf("GetContact() after toggle error = %v", err)
	}

	if after == nil {
		t.Fatal("GetContact() after toggle = nil, want contact")
	}

	if after.IsFavorite == before.IsFavorite {
		t.Fatalf("favorite = %t before toggle and %t after toggle", before.IsFavorite, after.IsFavorite)
	}
}

func TestServiceCreateNoteRequiresUser(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	err := h.svc.createNote(3, nil, "hello")

	if !errors.Is(err, errUnauthorized) {
		t.Fatalf("createNote() error = %v, want errUnauthorized", err)
	}
}

func TestServiceCreateNoteRejectsEmptyBody(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	user, err := database.FindUserByID(h.db, 1)

	if err != nil {
		t.Fatalf("FindUserByID() error = %v", err)
	}

	err = h.svc.createNote(3, user, "   ")

	if !errors.Is(err, errEmptyNoteBody) {
		t.Fatalf("createNote() error = %v, want errEmptyNoteBody", err)
	}
}

func TestServiceCreateNote(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	user, err := database.FindUserByID(h.db, 1)

	if err != nil {
		t.Fatalf("FindUserByID() error = %v", err)
	}

	if user == nil {
		t.Fatal("FindUserByID() = nil, want user")
	}

	before, err := database.ListContactNotes(h.db, 3)

	if err != nil {
		t.Fatalf("ListContactNotes() before create error = %v", err)
	}

	if err := h.svc.createNote(3, user, "  hello  "); err != nil {
		t.Fatalf("createNote() error = %v", err)
	}

	after, err := database.ListContactNotes(h.db, 3)

	if err != nil {
		t.Fatalf("ListContactNotes() after create error = %v", err)
	}

	if len(after) != len(before)+1 {
		t.Fatalf("len(ListContactNotes()) = %d, want %d", len(after), len(before)+1)
	}

	found := false

	for _, note := range after {
		if note.ContactID == 3 && note.UserID == 1 && note.Body == "hello" {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("created note not found in %#v", after)
	}
}

func TestServiceUpdateOrganization(t *testing.T) {
	t.Parallel()

	h := newServiceTestHarness(t)

	if err := h.svc.updateOrganization(1, organizationForm{Name: "  Acme HQ  "}); err != nil {
		t.Fatalf("updateOrganization() error = %v", err)
	}

	org, err := database.GetOrganization(h.db, 1)

	if err != nil {
		t.Fatalf("GetOrganization() error = %v", err)
	}

	if org == nil {
		t.Fatal("GetOrganization() = nil, want organization")
	}

	if org.Name != "Acme HQ" {
		t.Fatalf("organization name = %q, want %q", org.Name, "Acme HQ")
	}
}
