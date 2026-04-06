package crm

import (
	"testing"

	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func TestRepositoryMethods(t *testing.T) {
	t.Parallel()

	h := newCRMHarness(t)
	repo, err := newRepository(h.db)

	if err != nil {
		t.Fatalf("newRepository() error = %v", err)
	}

	if _, err := newRepository(nil); err == nil {
		t.Fatal("newRepository(nil) error = nil, want error")
	}

	if _, err := repo.CountContacts(); err != nil {
		t.Fatalf("CountContacts() error = %v", err)
	}

	if _, err := repo.CountOrganizations(); err != nil {
		t.Fatalf("CountOrganizations() error = %v", err)
	}

	if _, err := repo.CountNotes(); err != nil {
		t.Fatalf("CountNotes() error = %v", err)
	}

	if _, err := repo.ListRecentNotes(3); err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}

	if _, err := repo.ListContactsPaginated("", false, nil, "next", 5); err != nil {
		t.Fatalf("ListContactsPaginated() error = %v", err)
	}

	contact, err := repo.GetContact(1)

	if err != nil {
		t.Fatalf("GetContact() error = %v", err)
	}

	if contact == nil {
		t.Fatal("GetContact(1) returned nil")
	}

	if _, err := repo.ListContactNotes(1); err != nil {
		t.Fatalf("ListContactNotes() error = %v", err)
	}

	if _, err := repo.ListOrganizations(""); err != nil {
		t.Fatalf("ListOrganizations() error = %v", err)
	}

	if _, err := repo.ListOrganizationsPaginated("", 1, 5); err != nil {
		t.Fatalf("ListOrganizationsPaginated() error = %v", err)
	}

	org, err := repo.GetOrganization(1)

	if err != nil {
		t.Fatalf("GetOrganization() error = %v", err)
	}

	if org == nil {
		t.Fatal("GetOrganization(1) returned nil")
	}

	if _, err := repo.ListContactsByOrganization(1); err != nil {
		t.Fatalf("ListContactsByOrganization() error = %v", err)
	}

	if _, err := repo.ListContactsByOrgPaginated(1, nil, "next", 5); err != nil {
		t.Fatalf("ListContactsByOrgPaginated() error = %v", err)
	}

	id, err := repo.CreateContact(database.Contact{
		FirstName: "Repo",
		LastName:  "Created",
		Email:     "repo.created@example.com",
		Phone:     "555-0199",
	})

	if err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}

	if err := repo.UpdateContact(id, database.Contact{
		FirstName: "Repo",
		LastName:  "Updated",
		Email:     "repo.updated@example.com",
		Phone:     "555-0200",
	}); err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}

	if err := repo.ToggleContactFavorite(id); err != nil {
		t.Fatalf("ToggleContactFavorite() error = %v", err)
	}

	if _, err := repo.CreateNote(id, h.user.ID, "Repository note"); err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}

	if err := repo.DeleteContact(id); err != nil {
		t.Fatalf("DeleteContact() error = %v", err)
	}

	if err := repo.UpdateOrganization(1, "Repository Updated Org"); err != nil {
		t.Fatalf("UpdateOrganization() error = %v", err)
	}
}
