package features

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) persistentLayoutsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Layouts/PersistentLayouts", httpx.Props{})
}

func (a app) persistentLayoutsPage2Handler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Layouts/PersistentLayoutsPageTwo", httpx.Props{})
}

func (a app) nestedLayoutsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Layouts/NestedLayouts", httpx.Props{})
}

func (a app) headHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Layouts/Head", httpx.Props{})
}

func (a app) layoutPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Layouts/LayoutProps", httpx.Props{})
}
