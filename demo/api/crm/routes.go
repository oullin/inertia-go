package crm

import (
	"fmt"
	"net/http"
)

type app struct {
	deps    Deps
	service service
}

// RegisterRoutes mounts the CRM HTTP routes onto the provided mux.
func RegisterRoutes(mux *http.ServeMux, deps Deps) error {
	app, err := newApp(deps)

	if err != nil {
		return fmt.Errorf("crm: %w", err)
	}

	auth := func(h http.HandlerFunc) http.Handler {
		return deps.RequireAuth(h)
	}

	mux.Handle("/dashboard", auth(app.dashboardHandler))
	mux.Handle("/contacts", auth(app.contactsHandler))
	mux.Handle("/contacts/create", auth(app.contactsCreateHandler))
	mux.Handle("/contacts/", auth(app.contactByIDHandler))
	mux.Handle("/organizations", auth(app.organizationsHandler))
	mux.Handle("/organizations/", auth(app.organizationByIDHandler))

	return nil
}

func newApp(deps Deps) (app, error) {
	repo, err := newRepository(deps.DB)

	if err != nil {
		return app{}, err
	}

	return app{
		deps:    deps,
		service: newService(repo),
	}, nil
}
