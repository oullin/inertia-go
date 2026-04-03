package features

import (
	"net/http"
	"time"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) globalEventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Events/GlobalEvents", httpx.Props{})
	case http.MethodPost:
		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Event", Message: "Global event action completed."})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.events.global-events", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) visitCallbacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Events/VisitCallbacks", httpx.Props{})
	case http.MethodPost:
		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Callback", Message: "Visit callback action completed."})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.events.visit-callbacks", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) progressHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Events/Progress", httpx.Props{})
}

func (a app) progressSlowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	a.deps.Render(w, r, "Features/Events/Progress", httpx.Props{})
}
