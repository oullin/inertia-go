package crm

import (
	"net/http"
)

type app struct {
	deps    Deps
	service service
}

// RegisterRoutes mounts the CRM HTTP routes onto the provided mux.
func RegisterRoutes(mux *http.ServeMux, deps Deps) {
	app := newApp(deps)

	mux.Handle("/dashboard", deps.RequireAuth(http.HandlerFunc(app.dashboardHandler)))
	mux.Handle("/contacts", deps.RequireAuth(http.HandlerFunc(app.contactsHandler)))
	mux.Handle("/contacts/create", deps.RequireAuth(http.HandlerFunc(app.contactsCreateHandler)))
	mux.Handle("/contacts/", deps.RequireAuth(http.HandlerFunc(app.contactByIDHandler)))
	mux.Handle("/organizations", deps.RequireAuth(http.HandlerFunc(app.organizationsHandler)))
	mux.Handle("/organizations/", deps.RequireAuth(http.HandlerFunc(app.organizationByIDHandler)))
}

func newApp(deps Deps) app {
	return app{
		deps:    deps,
		service: newService(deps.Repository),
	}
}
