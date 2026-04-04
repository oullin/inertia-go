package features

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/props"
	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/httputil"
)

const contactsPerPage = 15

func (a app) deferredPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/DataLoading/DeferredProps", httpx.Props{
		"quickStat": 42,
		"slowStats": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 800*time.Millisecond) != nil {
				return nil
			}

			totalContacts, err := database.CountContacts(a.deps.DB)

			if err != nil {
				slog.Error("count contacts", "error", err)
			}

			return map[string]any{
				"totalContacts":  totalContacts,
				"totalFavorites": 12,
			}
		}, "slow"),
		"heavyData": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 1500*time.Millisecond) != nil {
				return nil
			}

			items := make([]map[string]any, 20)

			for i := range items {
				items[i] = map[string]any{"id": i + 1, "name": fmt.Sprintf("Item %d", i+1)}
			}

			return items
		}, "heavy"),
	})
}

func (a app) partialReloadsHandler(w http.ResponseWriter, r *http.Request) {
	users := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
		{"id": 3, "name": "Carol"},
	}

	a.deps.Render(w, r, "Features/DataLoading/PartialReloads", httpx.Props{
		"users":        users,
		"stats":        map[string]any{"total": 3, "active": 2},
		"timestamp":    time.Now().Format(time.RFC3339),
		"randomNumber": rand.Intn(1000),
	})
}

func (a app) infiniteScrollHandler(w http.ResponseWriter, r *http.Request) {
	var cursor *string

	if c := r.URL.Query().Get("cursor"); c != "" {
		cursor = &c
	}

	page, err := database.ListContactsPaginated(a.deps.DB, "", false, cursor, "next", contactsPerPage)

	if err != nil {
		slog.Error("list contacts paginated", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	items := make([]map[string]any, 0, len(page.Data))

	for _, c := range page.Data {
		item := map[string]any{
			"id":          c.ID,
			"first_name":  c.FirstName,
			"last_name":   c.LastName,
			"email":       c.Email,
			"is_favorite": c.IsFavorite,
		}

		if c.OrganizationID != nil {
			item["organization"] = map[string]any{"name": c.OrganizationName}
		}

		items = append(items, item)
	}

	a.deps.Render(w, r, "Features/DataLoading/InfiniteScroll", httpx.Props{
		"contacts": map[string]any{
			"data":        items,
			"next_cursor": page.NextCursor,
		},
	})
}

func (a app) whenVisibleHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/DataLoading/WhenVisible", httpx.Props{
		"section1": props.Optional(func() any {
			if httputil.SleepCtx(r.Context(), 500*time.Millisecond) != nil {
				return nil
			}

			return map[string]any{"title": "Section 1", "content": "Loaded on visibility."}
		}),
		"section2": props.Optional(func() any {
			if httputil.SleepCtx(r.Context(), 800*time.Millisecond) != nil {
				return nil
			}

			return map[string]any{"title": "Section 2", "content": "Also loaded lazily."}
		}),
		"section3": props.Optional(func() any {
			if httputil.SleepCtx(r.Context(), 1000*time.Millisecond) != nil {
				return nil
			}

			return map[string]any{"title": "Section 3", "content": "Third lazy section."}
		}),
	})
}

func (a app) pollingHandler(w http.ResponseWriter, r *http.Request) {
	contactCount, err := database.CountContacts(a.deps.DB)

	if err != nil {
		slog.Error("count contacts", "error", err)
	}

	a.deps.Render(w, r, "Features/DataLoading/Polling", httpx.Props{
		"currentTime":  time.Now().Format(time.RFC3339),
		"randomNumber": rand.Intn(1000),
		"contactCount": contactCount,
	})
}

func (a app) propMergingHandler(w http.ResponseWriter, r *http.Request) {
	notifications := []map[string]any{
		{"id": 1, "text": "New contact added", "time": "2m ago"},
		{"id": 2, "text": "Organization updated", "time": "5m ago"},
		{"id": 3, "text": "Note created", "time": "12m ago"},
	}

	activities := []map[string]any{
		{"id": 1, "action": "Logged in", "time": "1m ago"},
		{"id": 2, "action": "Viewed contacts", "time": "3m ago"},
	}

	a.deps.Render(w, r, "Features/DataLoading/PropMerging", httpx.Props{
		"notifications": props.Merge(notifications),
		"activities":    props.Merge(activities),
		"timestamp":     time.Now().Format(time.RFC3339),
	})
}

func (a app) optionalPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/DataLoading/OptionalProps", httpx.Props{
		"regularData": map[string]any{"message": "Always loaded"},
		"optionalData": props.Optional(func() any {
			return map[string]any{"message": "Loaded on demand"}
		}),
		"deferredData": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 500*time.Millisecond) != nil {
				return nil
			}

			return map[string]any{"message": "Loaded asynchronously"}
		}),
	})
}

func (a app) oncePropsHandler(w http.ResponseWriter, r *http.Request) {
	pageNum, _ := strconv.Atoi(r.PathValue("page"))

	if pageNum < 1 {
		pageNum = 1
	}

	a.deps.Render(w, r, "Features/DataLoading/OnceProps", httpx.Props{
		"page":       pageNum,
		"staticData": props.Once(map[string]any{"cached": true, "message": "This won't change on reload"}),
		"freshData":  map[string]any{"timestamp": time.Now().Format(time.RFC3339)},
		"dynamicData": map[string]any{
			"page":      pageNum,
			"generated": time.Now().Format(time.RFC3339),
		},
	})
}
