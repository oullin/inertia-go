package features

type app struct {
	deps Container
}

func newApp(deps Container) app {
	return app{deps: deps}
}
