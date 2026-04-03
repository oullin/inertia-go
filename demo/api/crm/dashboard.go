package crm

import (
	"net/http"
	"time"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/props"
)

func (a app) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	activity, _ := a.service.recentActivity(10)

	a.deps.Render(w, r, "Crm/Dashboard", httpx.Props{
		"recentActivity": recentActivityProps(activity),
		"totalContacts": props.Defer(func() any {
			time.Sleep(150 * time.Millisecond)

			return a.service.countContacts()
		}, "stats"),
		"totalOrganizations": props.Defer(func() any {
			time.Sleep(150 * time.Millisecond)

			return a.service.countOrganizations()
		}, "stats"),
		"recentNotesCount": props.Defer(func() any {
			time.Sleep(150 * time.Millisecond)

			return a.service.countNotes()
		}, "stats"),
	})
}
