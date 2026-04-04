package crm

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/props"
	"github.com/oullin/inertia-go/demo/api/internal/httputil"
)

func (a app) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	activity, err := a.service.recentActivity(10)

	if err != nil {
		slog.Error("recent activity", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	a.deps.Render(w, r, "Crm/Dashboard", httpx.Props{
		"recentActivity": recentActivityProps(activity),
		"totalContacts": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 150*time.Millisecond) != nil {
				return nil
			}

			n, err := a.service.countContacts()

			if err != nil {
				slog.Error("count contacts", "error", err)
			}

			return n
		}, "stats"),
		"totalOrganizations": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 150*time.Millisecond) != nil {
				return nil
			}

			n, err := a.service.countOrganizations()

			if err != nil {
				slog.Error("count organizations", "error", err)
			}

			return n
		}, "stats"),
		"recentNotesCount": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 150*time.Millisecond) != nil {
				return nil
			}

			n, err := a.service.countNotes()

			if err != nil {
				slog.Error("count notes", "error", err)
			}

			return n
		}, "stats"),
	})
}
