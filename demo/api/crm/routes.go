package crm

import (
	"fmt"
	"net/http"

	"github.com/oullin/inertia-go/core/wayfinder"
)

type app struct {
	deps Container
	repo *databaseRepository
}

// RegisterRoutes mounts the CRM HTTP routes onto the provided mux.
func RegisterRoutes(routes *wayfinder.Registry, mux *http.ServeMux, deps Container) error {
	app, err := newApp(deps)

	if err != nil {
		return fmt.Errorf("crm: %w", err)
	}

	auth := func(h http.HandlerFunc) http.Handler {
		return deps.RequireAuth(h)
	}

	routes.Handle("dashboard", auth(app.dashboardHandler), mux)
	routes.Handle("contacts.index", auth(app.contactsHandler), mux)
	routes.Handle("contacts.create", auth(app.contactsCreateHandler), mux)
	mux.Handle("/contacts/", auth(app.contactByIDHandler))
	routes.Handle("organizations.index", auth(app.organizationsHandler), mux)
	mux.Handle("/organizations/", auth(app.organizationByIDHandler))

	return nil
}

func newApp(deps Container) (app, error) {
	repo, err := newRepository(deps.DB)

	if err != nil {
		return app{}, err
	}

	return app{
		deps: deps,
		repo: repo,
	}, nil
}
