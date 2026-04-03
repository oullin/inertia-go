package features

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) rememberHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/State/Remember", httpx.Props{})
}

func (a app) flashDataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/State/FlashData", httpx.Props{})
	case http.MethodPost:
		a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Flash sent", Message: "This is a success flash message."})
		a.deps.Redirect(w, r, a.deps.RouteURL("features.state.flash-data", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) flashDataActionHandler(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimPrefix(r.URL.Path, "/features/state/flash-data/")
	action = strings.Trim(action, "/")

	switch action {
	case "error":
		a.deps.SetFlash(w, flash.Message{Kind: "error", Title: "Error", Message: "Something went wrong!"})
	case "warning":
		a.deps.SetFlash(w, flash.Message{Kind: "warning", Title: "Warning", Message: "Proceed with caution."})
	default:
		a.deps.SetFlash(w, flash.Message{Kind: "info", Title: "Info", Message: "Here's some information."})
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("features.state.flash-data", nil))
}

func (a app) sharedPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/State/SharedProps", httpx.Props{})
}
