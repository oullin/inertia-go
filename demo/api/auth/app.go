package auth

// App bundles the auth HTTP handlers, guards, and session helpers.
type App struct {
	deps    Container
	service service
}

// New builds the auth runtime with the provided host integrations.
func New(deps Container) App {
	return App{
		deps:    deps,
		service: newService(deps.DB),
	}
}
