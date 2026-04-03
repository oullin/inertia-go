package crm

import (
	"fmt"
	"net/http"
	"net/mail"
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
	}
}

func newContactFormFromContact(contact database.Contact) contactForm {
	form := contactForm{
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
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
	} else if len(f.FirstName) > 255 {
		errors["first_name"] = "First name must not exceed 255 characters."
	}

	if f.LastName == "" {
		errors["last_name"] = "Last name is required."
	} else if len(f.LastName) > 255 {
		errors["last_name"] = "Last name must not exceed 255 characters."
	}

	if f.Email == "" {
		errors["email"] = "A valid email address is required."
	} else if len(f.Email) > 255 {
		errors["email"] = "Email must not exceed 255 characters."
	} else if _, err := mail.ParseAddress(f.Email); err != nil {
		errors["email"] = "A valid email address is required."
	}

	if len(f.Phone) > 255 {
		errors["phone"] = "Phone must not exceed 255 characters."
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
