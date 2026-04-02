package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/props"
	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/seed"
)

type flashPayload struct {
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

const flashCookieName = "beacon_flash"

var db *sql.DB

func registerDashboardRoutes(mux *http.ServeMux) {
	mux.Handle("/", i.Middleware(http.HandlerFunc(rootHandler)))
	mux.Handle("/dashboard", i.Middleware(http.HandlerFunc(overviewHandler)))
	mux.Handle("/dashboard/navigation", i.Middleware(http.HandlerFunc(navigationHandler)))
	mux.Handle("/dashboard/scroll", i.Middleware(http.HandlerFunc(scrollHandler)))
	mux.Handle("/dashboard/redirects", i.Middleware(http.HandlerFunc(redirectsHandler)))
	mux.Handle("/dashboard/redirects/redirect", i.Middleware(http.HandlerFunc(redirectsRedirectHandler)))
	mux.Handle("/dashboard/redirects/location", i.Middleware(http.HandlerFunc(redirectsLocationHandler)))
	mux.Handle("/dashboard/forms", i.Middleware(http.HandlerFunc(formsHandler)))
	mux.Handle("/dashboard/forms/invite", i.Middleware(http.HandlerFunc(inviteHandler)))
	mux.Handle("/dashboard/forms/upload", i.Middleware(http.HandlerFunc(uploadHandler)))
	mux.Handle("/dashboard/forms/escalate", i.Middleware(http.HandlerFunc(escalateHandler)))
	mux.Handle("/dashboard/forms/http-preview", http.HandlerFunc(httpPreviewHandler))
	mux.Handle("/dashboard/data", i.Middleware(http.HandlerFunc(dataHandler)))
	mux.Handle("/dashboard/feed", i.Middleware(http.HandlerFunc(feedHandler)))
	mux.Handle("/dashboard/state", i.Middleware(http.HandlerFunc(stateHandler)))
	mux.Handle("/dashboard/state/error", http.HandlerFunc(stateErrorHandler))
	mux.Handle("/dashboard/seed", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seed.Run(db)
		w.WriteHeader(http.StatusOK)
	}))

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	i.Redirect(w, r, "/dashboard")
}

func overviewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := inertia.SetTitle(r.Context(), "Overview - Progressive Oullin")
	ctx = inertia.SetMeta(ctx,
		httpx.MetaTag{Name: "description", Content: "Dashboard overview with key operational metrics"},
		httpx.MetaTag{Property: "og:title", Content: "Overview - Progressive Oullin"},
	)

	renderDashboardWithContext(w, r.WithContext(ctx), "Dashboard/Overview", httpx.Props{
		"pageTitle":       "Overview",
		"recentInvites":   database.ListInvites(db),
		"recentUploads":   database.ListUploads(db),
		"recentApprovals": database.ListApprovals(db),
		"stats": map[string]any{
			"invites":     database.InviteCount(db),
			"uploads":     database.UploadCount(db),
			"approvals":   database.ApprovalCount(db),
			"escalations": database.GetCounter(db, "priority_escalations"),
		},
	})
}

func navigationHandler(w http.ResponseWriter, r *http.Request) {
	pane := query(r, "pane", "links")

	renderDashboard(w, r, "Dashboard/Navigation", httpx.Props{
		"pageTitle":    "Navigation and visit control",
		"pageSubtitle": "Compare declarative links and manual visits with preserved state and view transitions.",
		"pane":         pane,
		"visitMatrix": []map[string]any{
			{"title": "Links", "summary": "Standard Inertia links with active state and view transitions.", "mode": "link"},
			{"title": "Manual visits", "summary": "router.get / router.visit for controlled query and state handling.", "mode": "manual"},
			{"title": "Redirects", "summary": "POST mutations redirect back through the adapter middleware.", "mode": "redirect"},
			{"title": "Location visits", "summary": "Force a full page navigation when you need a hard handoff.", "mode": "location"},
		},
		"visitTargets": []map[string]any{
			{"label": "Dashboard overview", "href": "/dashboard", "detail": "Standard page transition"},
			{"label": "Data loading lab", "href": "/dashboard/data", "detail": "Route used by prefetch + deferred demos"},
			{"label": "Infinite feed", "href": "/dashboard/feed", "detail": "Long-scroll route with feed pagination"},
		},
	})
}

func scrollHandler(w http.ResponseWriter, r *http.Request) {
	preserve := query(r, "preserve", "keep")

	renderDashboard(w, r, "Dashboard/Scroll", httpx.Props{
		"pageTitle":    "Scroll preservation",
		"pageSubtitle": "Compare preserveScroll behavior and layout persistence across Inertia visits.",
		"preserve":     preserve,
		"scrollLedger": ledgerRows(preserve),
		"layoutNotes": []map[string]any{
			{"label": "Persistent shell", "detail": "The sidebar and workspace header survive page changes via the shared layout."},
			{"label": "View transitions", "detail": "Primary route buttons request browser view transitions where available."},
			{"label": "Preserve scroll", "detail": "Buttons below intentionally toggle preserveScroll behavior on revisit."},
		},
	})
}

func redirectsHandler(w http.ResponseWriter, r *http.Request) {
	renderDashboard(w, r, "Dashboard/Redirects", httpx.Props{
		"pageTitle":    "Redirects and location visits",
		"pageSubtitle": "Exercise POST-redirect-back flows and full-page location handoffs.",
		"redirectNotes": []map[string]any{
			{"label": "POST-redirect-back", "detail": "The adapter intercepts a 302 after a POST and replays it as an Inertia visit, keeping the SPA shell intact."},
			{"label": "Flash messaging", "detail": "The redirect carries a cookie-based flash message that surfaces on the destination page."},
			{"label": "Pane targeting", "detail": "The redirect URL includes query params so the destination component can highlight the relevant section."},
		},
		"locationNotes": []map[string]any{
			{"label": "Full-page handoff", "detail": "i.Location() forces the browser to do a hard navigation, bypassing the Inertia SPA layer entirely."},
			{"label": "Cross-route jump", "detail": "The demo jumps to /dashboard/state, proving the browser fully reloads rather than doing an XHR swap."},
			{"label": "Use cases", "detail": "Location visits are ideal for OAuth callbacks, file downloads, or any flow that must leave the SPA."},
		},
	})
}

func redirectsRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Redirect completed",
		Message: "The visit returned through the middleware as a standard redirect flow.",
	})

	i.Redirect(w, r, "/dashboard/redirects")
}

func redirectsLocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	i.Location(w, r, "/dashboard/state?source=location")
}

func formsHandler(w http.ResponseWriter, r *http.Request) {
	renderDashboard(w, r, "Dashboard/Forms", formsProps())
}

func inviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	role := strings.TrimSpace(r.FormValue("role"))
	errors := httpx.ValidationErrors{}

	if name == "" {
		errors["name"] = "Add a teammate name."
	}

	if email == "" || !strings.Contains(email, "@") {
		errors["email"] = "Use a valid email address."
	}

	if role == "" {
		errors["role"] = "Select a role for the invite."
	}

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		ctx = inertia.SetProp(ctx, "flash", map[string]any{
			"kind":    "error",
			"title":   "Invite incomplete",
			"message": "Fix the highlighted fields and resubmit the workflow.",
		})

		renderDashboardWithContext(w, r.WithContext(ctx), "Dashboard/Forms", formsProps())

		return
	}

	id := fmt.Sprintf("invite_%03d", database.InviteCount(db)+200)
	database.CreateInvite(db, id, name, email, role, "Pending")

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Invite queued",
		Message: fmt.Sprintf("%s was added to the onboarding queue as %s.", name, role),
	})

	i.Redirect(w, r, "/dashboard/forms")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	_ = r.ParseMultipartForm(8 << 20)

	label := strings.TrimSpace(r.FormValue("label"))
	file, header, err := r.FormFile("file")
	errors := httpx.ValidationErrors{}

	if label == "" {
		errors["label"] = "Give this upload a label."
	}

	if err != nil {
		errors["file"] = "Attach a file to upload."
	} else {
		file.Close()
	}

	if len(errors) > 0 {
		ctx := inertia.SetValidationErrors(r.Context(), errors)
		ctx = inertia.SetProp(ctx, "flash", map[string]any{
			"kind":    "error",
			"title":   "Upload blocked",
			"message": "Provide both a label and a file before submitting.",
		})

		renderDashboardWithContext(w, r.WithContext(ctx), "Dashboard/Forms", formsProps())

		return
	}

	id := fmt.Sprintf("upload_%02d", database.UploadCount(db)+50)
	database.CreateUpload(db, id, label, header.Filename, humanSize(header.Size), "Queued")

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Upload accepted",
		Message: fmt.Sprintf("%s is now in the document pipeline.", header.Filename),
	})

	i.Redirect(w, r, "/dashboard/forms")
}

func escalateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	database.IncrementCounter(db, "priority_escalations")

	id := fmt.Sprintf("approval_%02d", database.ApprovalCount(db)+50)
	database.CreateApproval(db, id, "Priority routing promotion", "Synced")

	setFlash(w, flashPayload{
		Kind:    "success",
		Title:   "Priority escalation synced",
		Message: "The queue count now reflects the promoted account.",
	})

	i.Redirect(w, r, "/dashboard/forms")
}

func httpPreviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	seats, _ := strconv.Atoi(defaultString(r.FormValue("seats"), "12"))
	tier := defaultString(strings.TrimSpace(r.FormValue("tier")), "growth")
	rate := 28

	switch tier {
	case "enterprise":
		rate = 42
	case "starter":
		rate = 16
	}

	payload := map[string]any{
		"tier":            tier,
		"seats":           seats,
		"monthlyEstimate": seats * rate,
		"annualEstimate":  seats * rate * 12,
		"recommended":     seats >= 25,
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(payload)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	timeframe := query(r, "timeframe", "30d")
	signalPage := queryInt(r, "signalPage", 1)
	capacityMode := query(r, "capacityMode", "baseline")
	signals := signalChunk(signalPage, timeframe)
	signalProp := any(signals)

	if signalPage > 1 && httpx.IsInertiaRequest(r) {
		signalProp = props.Merge(signals)
	}

	capacityProp := any(capacityModel("baseline"))

	if capacityMode == "burst" && httpx.IsInertiaRequest(r) {
		capacityProp = props.DeepMerge(capacityModel("burst"))
	}

	renderDashboard(w, r, "Dashboard/Data", httpx.Props{
		"pageTitle":    "Data loading and prop orchestration",
		"pageSubtitle": "Shared props, partial reloads, deferred panels, once metadata, polling, prefetching, and load-when-visible blocks.",
		"timeframe":    timeframe,
		"accountSummary": []map[string]any{
			{"label": "Qualified handoffs", "value": summaryValue(timeframe, 32, 44, 57), "delta": "+12%", "tone": "positive"},
			{"label": "Open invoices", "value": summaryValue(timeframe, 18, 26, 41), "delta": "-4%", "tone": "info"},
			{"label": "SLA variance", "value": summaryValue(timeframe, 6, 9, 13) + " hrs", "delta": "stable", "tone": "neutral"},
		},
		"liveMetrics": map[string]any{
			"activeWorkers": 14 + time.Now().Second()%3,
			"queueDepth":    38 + time.Now().Second()%7,
			"callbackLag":   fmt.Sprintf("%dm", 4+time.Now().Second()%6),
			"updatedAt":     time.Now().Format("15:04:05"),
		},
		"releaseNotes": props.Once(map[string]any{
			"tag":     "Loaded once",
			"title":   "Release snapshot",
			"message": "This panel is marked as a once prop and should be reused across partial reloads.",
		}),
		"trendNarrative": props.Defer(func() any {
			time.Sleep(250 * time.Millisecond)

			return map[string]any{
				"headline": fmt.Sprintf("The %s runway stays healthy after routing changes.", timeframeLabel(timeframe)),
				"body":     "Deferred props land after the shell is interactive, keeping the route fast while still showing richer analysis.",
			}
		}, "insights"),
		"trendBreakdown": props.Defer(func() any {
			time.Sleep(250 * time.Millisecond)

			return []map[string]any{
				{"label": "Expansion routed", "value": "84%"},
				{"label": "Collections auto-cleared", "value": "71%"},
				{"label": "Manual touch rate", "value": "22%"},
			}
		}, "insights"),
		"auditTrail": props.Optional(func() any {
			return []map[string]any{
				{"label": "Northwind payout review", "detail": "Triggered by queue threshold breach", "time": "5m ago"},
				{"label": "Atlas contract sync", "detail": "Loaded only once this block became visible", "time": "14m ago"},
				{"label": "Juniper invoice escalation", "detail": "Partial reload preserved the outer layout", "time": "28m ago"},
			}
		}),
		"signalsPage":    signalPage,
		"signalsHasMore": signalPage < 3,
		"signals":        signalProp,
		"capacityMode":   capacityMode,
		"capacityModel":  capacityProp,
		"prefetchTargets": []map[string]any{
			{"label": "Infinite feed", "href": "/dashboard/feed"},
			{"label": "State lab", "href": "/dashboard/state"},
		},
	})
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	team := query(r, "team", "all")
	kind := query(r, "kind", "all")
	page := queryInt(r, "feedPage", 1)
	feedItems, nextPage, totalLabel := feedChunk(page, team, kind)
	feedProp := props.Scroll(feedItems, "feedPage", page, previousPage(page), nextPage)

	if page > 1 && httpx.IsInertiaRequest(r) {
		feedProp = feedProp.Merge()
	}

	if hasResetHeader(r, "activityFeed") {
		feedProp = feedProp.Reset()
	}

	renderDashboard(w, r, "Dashboard/Feed", httpx.Props{
		"pageTitle":    "Infinite activity feed",
		"pageSubtitle": "A merged queue of customer, invoice, and routing events with reset-aware filters.",
		"filters": map[string]any{
			"team": team,
			"kind": kind,
		},
		"feedSummary": map[string]any{
			"title":     totalLabel,
			"subtitle":  "The feed route uses scroll metadata plus merged partial reloads.",
			"page":      page,
			"hasNext":   nextPage != nil,
			"teamLabel": strings.Title(team),
			"kindLabel": strings.Title(kind),
		},
		"activityFeed": feedProp,
	})
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	mode := query(r, "mode", "default")
	ctx := r.Context()

	switch mode {
	case "encrypted":
		ctx = inertia.SetEncryptHistory(ctx)
	case "clear":
		ctx = inertia.SetClearHistory(ctx)
	}

	renderDashboardWithContext(w, r.WithContext(ctx), "Dashboard/State", httpx.Props{
		"pageTitle":    "State, progress, and failure modes",
		"pageSubtitle": "Remember client state, inspect navigation events, and exercise history flags with a long scroll region.",
		"historyMode":  mode,
		"releaseTrack": props.Once(map[string]any{
			"build":  "2026.04.01",
			"status": "Stable",
			"note":   "This snapshot stays sticky while the page reloads around it.",
		}),
		"playbooks": []map[string]any{
			{"title": "Preserve scroll", "detail": "Revisit the page while keeping the long task list anchored."},
			{"title": "Encrypt history", "detail": "Mark the next navigation payload for encrypted history state."},
			{"title": "Clear history", "detail": "Drop previously encrypted state on the next round-trip."},
		},
		"longTasks": stateTasks(),
		"errorLinks": []map[string]any{
			{"label": "Trigger 404", "href": "/dashboard/state/error?code=404"},
			{"label": "Trigger 500", "href": "/dashboard/state/error?code=500"},
		},
	})
}

func stateErrorHandler(w http.ResponseWriter, r *http.Request) {
	switch queryInt(r, "code", 500) {
	case 404:
		http.Error(w, "dashboard demo missing route", http.StatusNotFound)
	default:
		http.Error(w, "dashboard demo internal error", http.StatusInternalServerError)
	}
}

func renderDashboard(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	renderDashboardWithContext(w, r, component, pageProps)
}

func renderDashboardWithContext(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	ctx := r.Context()

	if flash := consumeFlash(w, r); flash != nil {
		ctx = inertia.SetProp(ctx, "flash", flash)
	}

	ctx = inertia.SetProp(ctx, "route", map[string]any{
		"path":  r.URL.Path,
		"query": r.URL.Query(),
	})

	if err := i.Render(w, r.WithContext(ctx), component, pageProps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func setFlash(w http.ResponseWriter, flash flashPayload) {
	data, _ := json.Marshal(flash)

	http.SetCookie(w, &http.Cookie{
		Name:     flashCookieName,
		Value:    url.QueryEscape(string(data)),
		Path:     "/",
		HttpOnly: true,
	})
}

func consumeFlash(w http.ResponseWriter, r *http.Request) map[string]any {
	cookie, err := r.Cookie(flashCookieName)

	if err != nil {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:   flashCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	raw, err := url.QueryUnescape(cookie.Value)

	if err != nil {
		return nil
	}

	var flash map[string]any

	if err := json.Unmarshal([]byte(raw), &flash); err != nil {
		return nil
	}

	return flash
}

func ledgerRows(mode string) []map[string]any {
	rows := make([]map[string]any, 0, 18)

	for idx := 1; idx <= 18; idx++ {
		rows = append(rows, map[string]any{
			"id":      fmt.Sprintf("inv_%03d", idx),
			"account": fmt.Sprintf("Account %02d", idx),
			"owner":   []string{"Maya", "Ava", "Noah"}[idx%3],
			"amount":  fmt.Sprintf("$%d,4%02d", 9+idx, idx),
			"status":  []string{"Queued", "Review", "Closed"}[idx%3],
			"note":    fmt.Sprintf("Scroll demo row %d with preserve=%s.", idx, mode),
		})
	}

	return rows
}

func formsProps() httpx.Props {
	return httpx.Props{
		"pageTitle":     "Forms and mutations",
		"pageSubtitle":  "Run invite, upload, HTTP preview, and optimistic queue workflows against the fake backend.",
		"inviteRoles":   []string{"Analyst", "Operator", "Manager", "Executive"},
		"recentInvites": database.ListInvites(db),
		"uploadedFiles": database.ListUploads(db),
		"approvalSummary": map[string]any{
			"pending":       6,
			"priorityCount": database.GetCounter(db, "priority_escalations"),
			"sla":           "1h 12m",
		},
		"recentApprovals": database.ListApprovals(db),
	}
}

func signalChunk(page int, timeframe string) []map[string]any {
	all := [][]map[string]any{
		{
			{"id": "sig_01", "label": "Expansion queue recovered", "detail": timeframeLabel(timeframe) + " backlog fell after new routing", "tone": "positive"},
			{"id": "sig_02", "label": "Collections retry window", "detail": "Finance approved the second autopay cadence", "tone": "info"},
			{"id": "sig_03", "label": "Onboarding overlap", "detail": "Two accounts need shared ownership for seven days", "tone": "warning"},
		},
		{
			{"id": "sig_04", "label": "Renewal approvals", "detail": "Legal review volume normalized after the latest upload", "tone": "positive"},
			{"id": "sig_05", "label": "Manual invoice touch", "detail": "High-value invoices still need human review", "tone": "warning"},
			{"id": "sig_06", "label": "Partner handoff", "detail": "Channel attribution data loaded from the latest sheet", "tone": "info"},
		},
		{
			{"id": "sig_07", "label": "Cold queue audit", "detail": "The optional audit block surfaces follow-up debt", "tone": "neutral"},
			{"id": "sig_08", "label": "Feed merge ready", "detail": "Signals now mirror the feed append strategy", "tone": "positive"},
			{"id": "sig_09", "label": "History encryption flag", "detail": "State routes can opt into encrypted history snapshots", "tone": "info"},
		},
	}

	if page < 1 || page > len(all) {
		return nil
	}

	return all[page-1]
}

func capacityModel(mode string) map[string]any {
	if mode == "burst" {
		return map[string]any{
			"routing": map[string]any{
				"coverage": "99.1%",
				"overflow": "2 queues expanded",
			},
			"collections": map[string]any{
				"headroom": "18h",
				"eta":      "Recovered",
			},
		}
	}

	return map[string]any{
		"routing": map[string]any{
			"coverage": "96.8%",
			"overflow": "1 queue warm",
		},
		"collections": map[string]any{
			"headroom": "9h",
			"eta":      "Holding",
		},
	}
}

func feedChunk(page int, team, kind string) ([]map[string]any, any, string) {
	base := []map[string]any{
		{"id": "feed_01", "team": "sales", "kind": "handoff", "title": "Northwind moved to finance review", "detail": "The handoff included updated discount context.", "time": "5m ago"},
		{"id": "feed_02", "team": "finance", "kind": "invoice", "title": "Atlas invoice flagged for manual approval", "detail": "Queue rule escalated the amount automatically.", "time": "12m ago"},
		{"id": "feed_03", "team": "success", "kind": "task", "title": "Juniper onboarding task reopened", "detail": "Owner requested another compliance artifact.", "time": "19m ago"},
		{"id": "feed_04", "team": "sales", "kind": "alert", "title": "Cedar expansion marked urgent", "detail": "A pending payment blocks the renewal package.", "time": "26m ago"},
		{"id": "feed_05", "team": "finance", "kind": "invoice", "title": "Helio charge approved", "detail": "A manual override completed the queue review.", "time": "38m ago"},
		{"id": "feed_06", "team": "success", "kind": "handoff", "title": "Beacon EU rollout synced", "detail": "Ops and implementation now share the same milestone.", "time": "51m ago"},
		{"id": "feed_07", "team": "sales", "kind": "task", "title": "Mistral plan update posted", "detail": "The pricing preview landed through useHttp.", "time": "1h ago"},
		{"id": "feed_08", "team": "finance", "kind": "alert", "title": "Quarter-end reconciliation reminder", "detail": "Two CSV uploads still need final classification.", "time": "1h ago"},
		{"id": "feed_09", "team": "success", "kind": "task", "title": "Polar handoff due tomorrow", "detail": "The feed keeps appending without tearing the layout.", "time": "2h ago"},
		{"id": "feed_10", "team": "sales", "kind": "handoff", "title": "Orchid package returned to pricing", "detail": "A preserve-state visit kept local filters intact.", "time": "2h ago"},
		{"id": "feed_11", "team": "finance", "kind": "invoice", "title": "Summit payment received", "detail": "Deferred analytics updated after the shell loaded.", "time": "3h ago"},
		{"id": "feed_12", "team": "success", "kind": "alert", "title": "Bluebird onboarding owner changed", "detail": "The remembered filter state stayed pinned.", "time": "4h ago"},
	}

	filtered := make([]map[string]any, 0, len(base))

	for _, item := range base {
		if team != "all" && item["team"] != team {
			continue
		}

		if kind != "all" && item["kind"] != kind {
			continue
		}

		filtered = append(filtered, item)
	}

	pageSize := 4
	start := (page - 1) * pageSize

	if start >= len(filtered) {
		return nil, nil, fmt.Sprintf("%d matching events", len(filtered))
	}

	end := start + pageSize

	if end > len(filtered) {
		end = len(filtered)
	}

	var nextPage any

	if end < len(filtered) {
		nextPage = page + 1
	}

	return filtered[start:end], nextPage, fmt.Sprintf("%d matching events", len(filtered))
}

func stateTasks() []map[string]any {
	tasks := make([]map[string]any, 0, 24)

	for idx := 1; idx <= 24; idx++ {
		tasks = append(tasks, map[string]any{
			"id":      fmt.Sprintf("task_%02d", idx),
			"title":   fmt.Sprintf("Playbook checkpoint %02d", idx),
			"owner":   []string{"Ava", "Maya", "Noah", "Priya"}[idx%4],
			"status":  []string{"Pinned", "Queued", "Watching"}[idx%3],
			"summary": "Used to demonstrate remembered filters and preserved scroll regions.",
		})
	}

	return tasks
}

func summaryValue(timeframe string, a, b, c int) string {
	switch timeframe {
	case "7d":
		return strconv.Itoa(a)
	case "90d":
		return strconv.Itoa(c)
	default:
		return strconv.Itoa(b)
	}
}

func timeframeLabel(timeframe string) string {
	switch timeframe {
	case "7d":
		return "7 day"
	case "90d":
		return "90 day"
	default:
		return "30 day"
	}
}

func previousPage(page int) any {
	if page <= 1 {
		return nil
	}

	return page - 1
}

func hasResetHeader(r *http.Request, key string) bool {
	raw := r.Header.Get(httpx.HeaderReset)

	if raw == "" {
		return false
	}

	for _, item := range strings.Split(raw, ",") {
		if strings.TrimSpace(item) == key {
			return true
		}
	}

	return false
}

func query(r *http.Request, key, fallback string) string {
	val := strings.TrimSpace(r.URL.Query().Get(key))

	if val == "" {
		return fallback
	}

	return val
}

func queryInt(r *http.Request, key string, fallback int) int {
	val := query(r, key, "")

	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)

	if err != nil || parsed < 1 {
		return fallback
	}

	return parsed
}

func defaultString(val, fallback string) string {
	if strings.TrimSpace(val) == "" {
		return fallback
	}

	return strings.TrimSpace(val)
}

func humanSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < 1024*1024 {
		return fmt.Sprintf("%.0f KB", float64(bytes)/1024)
	}

	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
