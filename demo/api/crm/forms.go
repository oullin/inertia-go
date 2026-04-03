package crm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

type contactForm struct {
	OrganizationID string
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	Address        string
	City           string
	Region         string
	Country        string
	PostalCode     string
}

type organizationForm struct {
	Name string
}

func newContactForm(r *http.Request) contactForm {
	return contactForm{
		OrganizationID: strings.TrimSpace(r.FormValue("organization_id")),
		FirstName:      strings.TrimSpace(r.FormValue("first_name")),
		LastName:       strings.TrimSpace(r.FormValue("last_name")),
		Email:          strings.TrimSpace(r.FormValue("email")),
		Phone:          strings.TrimSpace(r.FormValue("phone")),
		Address:        strings.TrimSpace(r.FormValue("address")),
		City:           strings.TrimSpace(r.FormValue("city")),
		Region:         strings.TrimSpace(r.FormValue("region")),
		Country:        strings.TrimSpace(r.FormValue("country")),
		PostalCode:     strings.TrimSpace(r.FormValue("postal_code")),
	}
}

func newContactFormFromContact(contact database.Contact) contactForm {
	form := contactForm{
		FirstName:  contact.FirstName,
		LastName:   contact.LastName,
		Email:      contact.Email,
		Phone:      contact.Phone,
		Address:    contact.Address,
		City:       contact.City,
		Region:     contact.Region,
		Country:    contact.Country,
		PostalCode: contact.PostalCode,
	}

	if contact.OrganizationID != nil {
		form.OrganizationID = fmt.Sprintf("%d", *contact.OrganizationID)
	}

	return form
}

func emptyContactForm() contactForm {
	return contactForm{}
}

func (f contactForm) validate() httpx.ValidationErrors {
	errors := httpx.ValidationErrors{}

	if f.FirstName == "" {
		errors["first_name"] = "First name is required."
	}

	if f.LastName == "" {
		errors["last_name"] = "Last name is required."
	}

	if f.Email == "" || !strings.Contains(f.Email, "@") {
		errors["email"] = "A valid email address is required."
	}

	return errors
}

func (f contactForm) record() database.Contact {
	return database.Contact{
		OrganizationID: parseOrganizationID(f.OrganizationID),
		FirstName:      f.FirstName,
		LastName:       f.LastName,
		Email:          f.Email,
		Phone:          f.Phone,
		Address:        f.Address,
		City:           f.City,
		Region:         f.Region,
		Country:        f.Country,
		PostalCode:     f.PostalCode,
	}
}

func newOrganizationForm(r *http.Request) organizationForm {
	return organizationForm{
		Name: strings.TrimSpace(r.FormValue("name")),
	}
}

func (f organizationForm) validate() httpx.ValidationErrors {
	errors := httpx.ValidationErrors{}

	if f.Name == "" {
		errors["name"] = "Organization name is required."
	}

	return errors
}

func parseOrganizationID(raw string) *int64 {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return nil
	}

	id, err := strconv.ParseInt(raw, 10, 64)

	if err != nil {
		return nil
	}

	return &id
}
