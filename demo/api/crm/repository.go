package crm

import (
	"database/sql"

	"github.com/oullin/inertia-go/demo/api/internal/database"
)

type databaseRepository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) databaseRepository {
	return databaseRepository{db: db}
}

func (r databaseRepository) ListRecentNotes(limit int) ([]database.Note, error) {
	return database.ListRecentNotes(r.db, limit)
}

func (r databaseRepository) CountContacts() int {
	return database.CountContacts(r.db)
}

func (r databaseRepository) CountOrganizations() int {
	return database.CountOrganizations(r.db)
}

func (r databaseRepository) CountNotes() int {
	return database.CountNotes(r.db)
}

func (r databaseRepository) ListContacts(search string, favoritesOnly bool) ([]database.Contact, error) {
	return database.ListContacts(r.db, search, favoritesOnly)
}

func (r databaseRepository) GetContact(id int64) (*database.Contact, error) {
	return database.GetContact(r.db, id)
}

func (r databaseRepository) CreateContact(contact database.Contact) (int64, error) {
	return database.CreateContact(r.db, contact)
}

func (r databaseRepository) UpdateContact(id int64, contact database.Contact) error {
	return database.UpdateContact(r.db, id, contact)
}

func (r databaseRepository) ToggleContactFavorite(id int64) error {
	return database.ToggleContactFavorite(r.db, id)
}

func (r databaseRepository) ListContactNotes(contactID int64) ([]database.Note, error) {
	return database.ListContactNotes(r.db, contactID)
}

func (r databaseRepository) CreateNote(contactID, userID int64, body string) (int64, error) {
	return database.CreateNote(r.db, contactID, userID, body)
}

func (r databaseRepository) ListOrganizations(search string) ([]database.Organization, error) {
	return database.ListOrganizations(r.db, search)
}

func (r databaseRepository) GetOrganization(id int64) (*database.Organization, error) {
	return database.GetOrganization(r.db, id)
}

func (r databaseRepository) UpdateOrganization(id int64, name string) error {
	return database.UpdateOrganization(r.db, id, name)
}

func (r databaseRepository) ListContactsByOrganization(organizationID int64) ([]database.Contact, error) {
	return database.ListContactsByOrganization(r.db, organizationID)
}
