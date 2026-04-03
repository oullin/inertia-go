package features

type app struct {
	deps Deps
}

func newApp(deps Deps) app {
	return app{deps: deps}
}
