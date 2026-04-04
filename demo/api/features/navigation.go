package features

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/internal/httputil"
)

func (a app) linksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Navigation/Links", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Action", Message: "Link action processed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.links", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) preserveStateHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Navigation/PreserveState", httpx.Props{
		"serverCounter": 0,
		"timestamp":     time.Now().Format(time.RFC3339),
	})
}

func (a app) preserveScrollHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Navigation/PreserveScroll", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) viewTransitionsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Navigation/ViewTransitions", httpx.Props{})
}

func (a app) historyManagementHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Navigation/HistoryManagement", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "History", Message: "Action recorded in history."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.history-management", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) asyncRequestsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Navigation/AsyncRequests", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) asyncSlowHandler(w http.ResponseWriter, r *http.Request) {
	if httputil.SleepCtx(r.Context(), 2*time.Second) != nil {
		return
	}

	a.container.Render(w, r, "Features/Navigation/AsyncRequests", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) manualVisitsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Navigation/ManualVisits", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) redirectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Navigation/Redirects", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Redirected", Message: "Standard redirect completed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.redirects", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) redirectsActionHandler(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")

	switch action {
	case "back":
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Back", Message: "Redirected back."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.redirects", nil))
	case "to-route":
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Named route", Message: "Redirected to named route."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.redirects", nil))
	case "external":
		a.container.Location(w, r, "https://inertiajs.com")
	default:
		http.NotFound(w, r)
	}
}

func (a app) scrollManagementHandler(w http.ResponseWriter, r *http.Request) {
	items := make([]map[string]any, 50)

	for i := range items {
		items[i] = map[string]any{"id": i + 1, "title": fmt.Sprintf("Item #%d", i+1)}
	}

	a.container.Render(w, r, "Features/Navigation/ScrollManagement", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
		"items":     items,
	})
}

func (a app) instantVisitsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Navigation/InstantVisits", httpx.Props{
		"sourceTimestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) instantVisitTargetHandler(w http.ResponseWriter, r *http.Request) {
	items := make([]map[string]any, 50)

	for i := range items {
		items[i] = map[string]any{"id": i + 1, "title": fmt.Sprintf("Item #%d", i+1)}
	}

	a.container.Render(w, r, "Features/Navigation/InstantVisitTarget", httpx.Props{
		"greeting":        "Welcome to the target page!",
		"serverTimestamp": time.Now().Format(time.RFC3339),
		"items":           items,
	})
}

func (a app) urlFragmentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.container.Render(w, r, "Features/Navigation/UrlFragments", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.container.SetFlash(w, flash.Message{Kind: "success", Title: "Fragment", Message: "Fragment action processed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.url-fragments", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) urlFragmentsActionHandler(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")

	switch action {
	case "redirect-with-hash":
		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.url-fragments", nil)+"#section-2")
	case "preserve-fragment":
		a.container.Redirect(w, r, a.container.RouteURL("features.navigation.url-fragments", nil))
	default:
		http.NotFound(w, r)
	}
}
