package features

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) persistentLayoutsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/PersistentLayouts", httpx.Props{})
}

func (a app) persistentLayoutsPage2Handler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/PersistentLayoutsPageTwo", httpx.Props{})
}

func (a app) nestedLayoutsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/NestedLayouts", httpx.Props{
		"title": "Nested Layouts",
		"breadcrumbs": []map[string]string{
			{"title": "Features"},
			{"title": "Layouts"},
			{"title": "Nested Layouts"},
		},
	})
}

func (a app) headHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/Head", httpx.Props{})
}

func (a app) layoutPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Layouts/LayoutProps", httpx.Props{})
}
