package crm

import (
	"errors"
	"strings"

	"github.com/oullin/inertia-go/demo/api/internal/database"
)

type service struct {
	repo databaseRepository
}

var errUnauthorized = errors.New("crm: current user required")
var errEmptyNoteBody = errors.New("crm: note body required")

func newService(repo databaseRepository) service {
	return service{repo: repo}
}

func (s service) recentActivity(limit int) ([]database.Note, error) {
	return s.repo.ListRecentNotes(limit)
}

func (s service) countContacts() (int, error) {
	return s.repo.CountContacts()
}

func (s service) countOrganizations() (int, error) {
	return s.repo.CountOrganizations()
}

func (s service) countNotes() (int, error) {
	return s.repo.CountNotes()
}

func (s service) listContacts(search string, favoriteOnly bool) ([]database.Contact, error) {
	return s.repo.ListContacts(strings.TrimSpace(search), favoriteOnly)
}

func (s service) listContactsPaginated(search string, favoriteOnly bool, cursor *string, direction string) (database.CursorPage[database.Contact], error) {
	return s.repo.ListContactsPaginated(strings.TrimSpace(search), favoriteOnly, cursor, direction, 15)
}

func (s service) getContact(id int64) (*database.Contact, error) {
	return s.repo.GetContact(id)
}

func (s service) listContactNotes(contactID int64) ([]database.Note, error) {
	return s.repo.ListContactNotes(contactID)
}

func (s service) createContact(form contactForm) (int64, error) {
	return s.repo.CreateContact(form.record())
}

func (s service) updateContact(contactID int64, form contactForm) error {
	return s.repo.UpdateContact(contactID, form.record())
}

func (s service) deleteContact(contactID int64) error {
	return s.repo.DeleteContact(contactID)
}

func (s service) toggleFavorite(contactID int64) error {
	return s.repo.ToggleContactFavorite(contactID)
}

func (s service) createNote(contactID int64, user *database.User, body string) error {
	if user == nil {
		return errUnauthorized
	}

	body = strings.TrimSpace(body)

	if body == "" {
		return errEmptyNoteBody
	}

	_, err := s.repo.CreateNote(contactID, user.ID, body)

	return err
}

func (s service) listOrganizations(search string) ([]database.Organization, error) {
	return s.repo.ListOrganizations(strings.TrimSpace(search))
}

func (s service) listOrganizationsPaginated(search string, page int) (database.OffsetPage[database.Organization], error) {
	return s.repo.ListOrganizationsPaginated(strings.TrimSpace(search), page, 20)
}

func (s service) getOrganization(id int64) (*database.Organization, error) {
	return s.repo.GetOrganization(id)
}

func (s service) updateOrganization(organizationID int64, form organizationForm) error {
	return s.repo.UpdateOrganization(organizationID, form.Name)
}

func (s service) listContactsByOrganization(organizationID int64) ([]database.Contact, error) {
	return s.repo.ListContactsByOrganization(organizationID)
}

func (s service) listContactsByOrgPaginated(organizationID int64, cursor *string, direction string) (database.CursorPage[database.Contact], error) {
	return s.repo.ListContactsByOrgPaginated(organizationID, cursor, direction, 15)
}
