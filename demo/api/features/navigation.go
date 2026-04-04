package features

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) linksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Navigation/Links", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Action", Message: "Link action processed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.links", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) preserveStateHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Navigation/PreserveState", httpx.Props{
		"serverCounter": 0,
		"timestamp":     time.Now().Format(time.RFC3339),
	})
}

func (a app) preserveScrollHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Navigation/PreserveScroll", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) viewTransitionsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Navigation/ViewTransitions", httpx.Props{})
}

func (a app) historyManagementHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Navigation/HistoryManagement", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "History", Message: "Action recorded in history."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.history-management", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) asyncRequestsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Navigation/AsyncRequests", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) asyncSlowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	a.deps.Render(w, r, "Features/Navigation/AsyncRequests", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) manualVisitsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Navigation/ManualVisits", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) redirectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Navigation/Redirects", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Redirected", Message: "Standard redirect completed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.redirects", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) redirectsActionHandler(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimPrefix(r.URL.Path, "/features/navigation/redirects/")
	action = strings.Trim(action, "/")

	switch action {
	case "back":
		if err := a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Back", Message: "Redirected back."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.redirects", nil))
	case "to-route":
		if err := a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Named route", Message: "Redirected to named route."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.redirects", nil))
	case "external":
		a.deps.Location(w, r, "https://inertiajs.com")
	default:
		http.NotFound(w, r)
	}
}

func (a app) scrollManagementHandler(w http.ResponseWriter, r *http.Request) {
	items := make([]map[string]any, 50)

	for i := range items {
		items[i] = map[string]any{"id": i + 1, "title": fmt.Sprintf("Item #%d", i+1)}
	}

	a.deps.Render(w, r, "Features/Navigation/ScrollManagement", httpx.Props{
		"timestamp": time.Now().Format(time.RFC3339),
		"items":     items,
	})
}

func (a app) instantVisitsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Navigation/InstantVisits", httpx.Props{
		"sourceTimestamp": time.Now().Format(time.RFC3339),
	})
}

func (a app) instantVisitTargetHandler(w http.ResponseWriter, r *http.Request) {
	items := make([]map[string]any, 50)

	for i := range items {
		items[i] = map[string]any{"id": i + 1, "title": fmt.Sprintf("Item #%d", i+1)}
	}

	a.deps.Render(w, r, "Features/Navigation/InstantVisitTarget", httpx.Props{
		"greeting":        "Welcome to the target page!",
		"serverTimestamp": time.Now().Format(time.RFC3339),
		"items":           items,
	})
}

func (a app) urlFragmentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.deps.Render(w, r, "Features/Navigation/UrlFragments", httpx.Props{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	case http.MethodPost:
		if err := a.deps.SetFlash(w, flash.Message{Kind: "success", Title: "Fragment", Message: "Fragment action processed."}); err != nil {
			slog.Error("flash: set", "error", err)
		}

		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.url-fragments", nil))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a app) urlFragmentsActionHandler(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimPrefix(r.URL.Path, "/features/navigation/url-fragments/")
	action = strings.Trim(action, "/")

	switch action {
	case "redirect-with-hash":
		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.url-fragments", nil)+"#section-2")
	case "preserve-fragment":
		a.deps.Redirect(w, r, a.deps.RouteURL("features.navigation.url-fragments", nil))
	default:
		http.NotFound(w, r)
	}
}
