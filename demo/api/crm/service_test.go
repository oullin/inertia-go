package crm

import (
	"errors"
	"testing"

	"github.com/oullin/inertia-go/demo/api/internal/database"
)

type fakeRepository struct {
	createdContact     database.Contact
	updatedContactID   int64
	updatedContact     database.Contact
	toggledFavoriteID  int64
	createdNoteContact int64
	createdNoteUser    int64
	createdNoteBody    string
	updatedOrgID       int64
	updatedOrgName     string
}

func (f *fakeRepository) ListRecentNotes(limit int) ([]database.Note, error) {
	return nil, nil
}

func (f *fakeRepository) CountContacts() int {
	return 0
}

func (f *fakeRepository) CountOrganizations() int {
	return 0
}

func (f *fakeRepository) CountNotes() int {
	return 0
}

func (f *fakeRepository) ListContacts(search string, favoritesOnly bool) ([]database.Contact, error) {
	return nil, nil
}

func (f *fakeRepository) GetContact(id int64) (*database.Contact, error) {
	return nil, nil
}

func (f *fakeRepository) CreateContact(contact database.Contact) (int64, error) {
	f.createdContact = contact

	return 42, nil
}

func (f *fakeRepository) UpdateContact(id int64, contact database.Contact) error {
	f.updatedContactID = id
	f.updatedContact = contact

	return nil
}

func (f *fakeRepository) ToggleContactFavorite(id int64) error {
	f.toggledFavoriteID = id

	return nil
}

func (f *fakeRepository) ListContactNotes(contactID int64) ([]database.Note, error) {
	return nil, nil
}

func (f *fakeRepository) CreateNote(contactID, userID int64, body string) (int64, error) {
	f.createdNoteContact = contactID
	f.createdNoteUser = userID
	f.createdNoteBody = body

	return 1, nil
}

func (f *fakeRepository) ListOrganizations(search string) ([]database.Organization, error) {
	return nil, nil
}

func (f *fakeRepository) GetOrganization(id int64) (*database.Organization, error) {
	return nil, nil
}

func (f *fakeRepository) UpdateOrganization(id int64, name string) error {
	f.updatedOrgID = id
	f.updatedOrgName = name

	return nil
}

func (f *fakeRepository) ListContactsByOrganization(organizationID int64) ([]database.Contact, error) {
	return nil, nil
}

func TestServiceCreateContact(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	svc := newService(repo)
	id, err := svc.createContact(contactForm{
		OrganizationID: "8",
		FirstName:      "Mina",
		LastName:       "Cole",
		Email:          "mina@example.test",
	})

	if err != nil {
		t.Fatalf("createContact() error = %v", err)
	}

	if id != 42 {
		t.Fatalf("id = %d, want 42", id)
	}

	if repo.createdContact.FirstName != "Mina" {
		t.Fatalf("created contact first_name = %q", repo.createdContact.FirstName)
	}

	if repo.createdContact.OrganizationID == nil || *repo.createdContact.OrganizationID != 8 {
		t.Fatalf("organization_id = %v, want 8", repo.createdContact.OrganizationID)
	}
}

func TestServiceUpdateContact(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	svc := newService(repo)

	if err := svc.updateContact(9, contactForm{FirstName: "Mina"}); err != nil {
		t.Fatalf("updateContact() error = %v", err)
	}

	if repo.updatedContactID != 9 {
		t.Fatalf("updatedContactID = %d, want 9", repo.updatedContactID)
	}

	if repo.updatedContact.FirstName != "Mina" {
		t.Fatalf("updated first_name = %q", repo.updatedContact.FirstName)
	}
}

func TestServiceToggleFavorite(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	svc := newService(repo)

	if err := svc.toggleFavorite(5); err != nil {
		t.Fatalf("toggleFavorite() error = %v", err)
	}

	if repo.toggledFavoriteID != 5 {
		t.Fatalf("toggledFavoriteID = %d, want 5", repo.toggledFavoriteID)
	}
}

func TestServiceCreateNoteRequiresUser(t *testing.T) {
	t.Parallel()

	svc := newService(&fakeRepository{})
	err := svc.createNote(3, nil, "hello")

	if !errors.Is(err, errUnauthorized) {
		t.Fatalf("createNote() error = %v, want errUnauthorized", err)
	}
}

func TestServiceCreateNote(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	svc := newService(repo)
	user := &database.User{ID: 7}

	if err := svc.createNote(3, user, "  hello  "); err != nil {
		t.Fatalf("createNote() error = %v", err)
	}

	if repo.createdNoteContact != 3 || repo.createdNoteUser != 7 || repo.createdNoteBody != "hello" {
		t.Fatalf("note call = (%d, %d, %q)", repo.createdNoteContact, repo.createdNoteUser, repo.createdNoteBody)
	}
}

func TestServiceUpdateOrganization(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	svc := newService(repo)

	if err := svc.updateOrganization(11, organizationForm{Name: "Acme"}); err != nil {
		t.Fatalf("updateOrganization() error = %v", err)
	}

	if repo.updatedOrgID != 11 || repo.updatedOrgName != "Acme" {
		t.Fatalf("updateOrganization call = (%d, %q)", repo.updatedOrgID, repo.updatedOrgName)
	}
}
