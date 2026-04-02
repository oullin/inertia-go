package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/oullin/inertia-go/core/assert"
	"github.com/oullin/inertia-go/core/config"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/middleware"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

const testTemplate = `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`

func TestDashboardOverviewRoute(t *testing.T) {
	testInertia := newDashboardTestInertia(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req.RequestURI = "/dashboard"

	page := assert.AssertFromHandler(t, testInertia, overviewHandler, req)

	page.AssertComponent(t, "Dashboard/Overview")
	page.AssertURL(t, "/dashboard")
	page.AssertHasProp(t, "pageTitle")
	page.AssertHasProp(t, "stats")
	page.AssertHasProp(t, "app")
}

func TestDashboardInviteValidationErrors(t *testing.T) {
	testInertia := newDashboardTestInertia(t)
	body := strings.NewReader(url.Values{
		"name":  {""},
		"email": {"bad-email"},
		"role":  {""},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/invite", body)
	req.RequestURI = "/dashboard/forms/invite"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	page := assert.AssertFromHandler(t, testInertia, inviteHandler, req)

	page.AssertComponent(t, "Dashboard/Forms")

	errors, ok := page.Props["errors"].(map[string]any)

	if !ok {
		t.Fatalf("errors prop missing or invalid: %#v", page.Props["errors"])
	}

	if errors["name"] != "Add a teammate name." {
		t.Fatalf("unexpected name error: %#v", errors["name"])
	}

	if errors["email"] != "Use a valid email address." {
		t.Fatalf("unexpected email error: %#v", errors["email"])
	}

	if errors["role"] != "Select a role for the invite." {
		t.Fatalf("unexpected role error: %#v", errors["role"])
	}
}

func TestDashboardInviteSuccessRedirectsWithFlashCookie(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)
	body := strings.NewReader(url.Values{
		"name":  {"Nina Patel"},
		"email": {"nina@northstarhq.test"},
		"role":  {"Operator"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/invite", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderInertia, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusFound)
	}

	if got := resp.Header.Get("Location"); got != "/dashboard/forms" {
		t.Fatalf("location = %q, want %q", got, "/dashboard/forms")
	}

	cookies := resp.Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected flash cookie to be set")
	}
}

func TestDashboardInvitePrecognitionSkipsMutation(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	before := database.InviteCount(db)

	body := strings.NewReader(url.Values{
		"name":  {"Nina Patel"},
		"email": {"nina@northstarhq.test"},
		"role":  {"Operator"},
	}.Encode())

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/invite", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(httpx.HeaderPrecognition, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)

	w := httptest.NewRecorder()
	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	if got := w.Header().Get(httpx.HeaderPrecognition); got != "true" {
		t.Fatalf("Precognition header = %q, want %q", got, "true")
	}

	if after := database.InviteCount(db); after != before {
		t.Fatalf("invite count = %d, want %d", after, before)
	}
}

func TestDashboardFeedIncludesScrollMetadata(t *testing.T) {
	testInertia := newDashboardTestInertia(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard/feed", nil)
	req.RequestURI = "/dashboard/feed"

	page := assert.AssertFromHandler(t, testInertia, feedHandler, req)

	page.AssertComponent(t, "Dashboard/Feed")

	scroll, ok := page.ScrollProps["activityFeed"]

	if !ok {
		t.Fatal("missing scrollProps.activityFeed metadata")
	}

	if scroll.PageName != "feedPage" {
		t.Fatalf("pageName = %q, want %q", scroll.PageName, "feedPage")
	}

	if scroll.CurrentPage != float64(1) {
		t.Fatalf("currentPage = %#v, want 1", scroll.CurrentPage)
	}

	if scroll.NextPage != float64(2) {
		t.Fatalf("nextPage = %#v, want 2", scroll.NextPage)
	}

	if scroll.Reset {
		t.Fatal("reset should be false on first feed page")
	}
}

func TestSeedEndpointPopulatesState(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	req := httptest.NewRequest(http.MethodPost, "/dashboard/seed", nil)
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if count := database.InviteCount(db); count != 20 {
		t.Fatalf("invites = %d, want 20", count)
	}

	if count := database.UploadCount(db); count != 15 {
		t.Fatalf("uploads = %d, want 15", count)
	}

	if count := database.ApprovalCount(db); count != 15 {
		t.Fatalf("approvals = %d, want 15", count)
	}

	if count := database.GetCounter(db, "priority_escalations"); count != 18 {
		t.Fatalf("priority_escalations = %d, want 18", count)
	}
}

func newDashboardTestMux(t *testing.T) http.Handler {
	t.Helper()

	newDashboardTestInertia(t)

	mux := http.NewServeMux()
	registerDashboardRoutes(mux)

	cfg := config.DefaultI18n()
	cfg.URLPrefix = false

	return dashboardAppHandler(
		mux,
		middleware.CSRF(config.CSRFConfig{}, []byte("0123456789abcdef0123456789abcdef")),
		cfg,
	)
}

func issueDashboardCSRFCookie(t *testing.T, handler http.Handler) (*http.Cookie, string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return findCookie(t, w, "XSRF-TOKEN"), findCSRFMetaToken(t, w.Body.String())
}

func newDashboardTestInertia(t *testing.T) *inertia.Inertia {
	t.Helper()
	setupTestDB(t)

	testInertia, err := inertia.New(testTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatal(err)
	}

	testInertia.ShareProps(httpx.Props{
		"app": map[string]any{
			"name":        "Beacon Ops Console",
			"productLine": "Revenue Operations",
			"environment": "Sandbox",
		},
		"auth": map[string]any{
			"user": map[string]any{
				"name":     "Maya Tan",
				"title":    "Ops Director",
				"initials": "MT",
			},
		},
		"workspace": map[string]any{
			"name": "Northstar HQ",
			"plan": "Growth",
		},
	})

	i = testInertia

	t.Cleanup(func() {
		i = nil
	})

	return testInertia
}

func setupTestDB(t *testing.T) {
	t.Helper()

	testDB, err := database.Open(":memory:")

	if err != nil {
		t.Fatal(err)
	}

	database.CreateInviteAt(testDB, "invite_104", "Aria Lim", "aria@northstarhq.test", "Operator", "Accepted", timeAgo(9))
	database.CreateInviteAt(testDB, "invite_103", "Noah Chen", "noah@northstarhq.test", "Manager", "Pending", timeAgo(41))

	database.CreateUploadAt(testDB, "upload_01", "March collections export", "collections-march.csv", "84 KB", "Processed", timeAgo(12))
	database.CreateUploadAt(testDB, "upload_00", "Sales handoff packet", "handoff-q2.pdf", "1.2 MB", "Ready", timeAgo(120))

	database.CreateApprovalAt(testDB, "approval_02", "Northwind priority routing", "Approved", timeAgo(16))
	database.CreateApprovalAt(testDB, "approval_01", "Atlas expansion guardrail", "Queued", timeAgo(47))

	database.SetCounter(testDB, "priority_escalations", 4)

	db = testDB

	t.Cleanup(func() {
		testDB.Close()
		db = nil
	})
}

func timeAgo(minutes int) time.Time {
	return time.Now().Add(-time.Duration(minutes) * time.Minute)
}

func findCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()

	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}

	t.Fatalf("cookie %q not found", name)

	return nil
}

func findCSRFMetaToken(t *testing.T, body string) string {
	t.Helper()

	const prefix = `name="csrf-token" content="`
	start := strings.Index(body, prefix)

	if start == -1 {
		t.Fatal("csrf meta tag not found")
	}

	start += len(prefix)
	end := strings.Index(body[start:], `"`)

	if end == -1 {
		t.Fatal("csrf meta token not terminated")
	}

	return body[start : start+end]
}
