package auth

// App bundles the auth HTTP handlers, guards, and session helpers.
type App struct {
	container Container
	service   service
}

// New builds the auth runtime with the provided host integrations.
func New(container Container) App {
	return App{
		container: container,
		service:   newService(container.DB),
	}
}
