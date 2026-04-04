package features

import (
	"log/slog"
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
		if err := httpx.ParseForm(r); err != nil {
			slog.Error("flash: parse form", "error", err)
		}

		kind := r.FormValue("kind")

		var msg flash.Message

		switch kind {
		case "error":
			msg = flash.Message{Kind: "error", Title: "Error", Message: "Something went wrong!"}
		case "warning":
			msg = flash.Message{Kind: "warning", Title: "Warning", Message: "Proceed with caution."}
		default:
			msg = flash.Message{Kind: "success", Title: "Flash sent", Message: "This is a success flash message."}
		}

		if err := a.deps.SetFlash(w, msg); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.state.flash-data", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) flashDataActionHandler(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimPrefix(r.URL.Path, "/features/state/flash-data/")
	action = strings.Trim(action, "/")

	var msg flash.Message

	switch action {
	case "error":
		msg = flash.Message{Kind: "error", Title: "Error", Message: "Something went wrong!"}
	case "warning":
		msg = flash.Message{Kind: "warning", Title: "Warning", Message: "Proceed with caution."}
	default:
		msg = flash.Message{Kind: "info", Title: "Info", Message: "Here's some information."}
	}

	if err := a.deps.SetFlash(w, msg); err != nil {
		slog.Error("flash: set", "error", err)
	}

	a.deps.Redirect(w, r, a.deps.RouteURL("features.state.flash-data", nil))
}

func (a app) sharedPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/State/SharedProps", httpx.Props{})
}
