package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

// ---------------------------------------------------------------------------
// Method-not-allowed guards on direct handler calls
// ---------------------------------------------------------------------------

func TestInviteHandlerGetNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/forms/invite", nil)
	w := httptest.NewRecorder()

	inviteHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestUploadHandlerGetNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/forms/upload", nil)
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestEscalateHandlerGetNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/forms/escalate", nil)
	w := httptest.NewRecorder()

	escalateHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestRedirectsRedirectHandlerGetDirect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/redirects/redirect", nil)
	w := httptest.NewRecorder()

	redirectsRedirectHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestRedirectsLocationHandlerGetDirect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/redirects/location", nil)
	w := httptest.NewRecorder()

	redirectsLocationHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

// ---------------------------------------------------------------------------
// A. Pure utility functions
// ---------------------------------------------------------------------------

func TestQuery(t *testing.T) {
	tests := []struct {
		name, url, key, fallback, want string
	}{
		{"present", "/x?k=val", "k", "fb", "val"},
		{"empty", "/x?k=", "k", "fb", "fb"},
		{"absent", "/x", "k", "fb", "fb"},
		{"whitespace", "/x?k=%20%20", "k", "fb", "fb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)

			if got := query(r, tt.key, tt.fallback); got != tt.want {
				t.Errorf("query(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestQueryInt(t *testing.T) {
	tests := []struct {
		name, url string
		fallback  int
		want      int
	}{
		{"valid", "/x?n=3", 1, 3},
		{"non-numeric", "/x?n=abc", 1, 1},
		{"zero", "/x?n=0", 1, 1},
		{"negative", "/x?n=-2", 1, 1},
		{"empty", "/x?n=", 1, 1},
		{"absent", "/x", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)

			if got := queryInt(r, "n", tt.fallback); got != tt.want {
				t.Errorf("queryInt = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDefaultString(t *testing.T) {
	tests := []struct {
		name, val, fallback, want string
	}{
		{"non-empty", "hello", "fb", "hello"},
		{"trimmed", "  hello  ", "fb", "hello"},
		{"empty", "", "fb", "fb"},
		{"whitespace", "   ", "fb", "fb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultString(tt.val, tt.fallback); got != tt.want {
				t.Errorf("defaultString(%q) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestHumanSize(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"bytes", 512, "512 B"},
		{"one KB", 1024, "1 KB"},
		{"kilobytes", 2048, "2 KB"},
		{"one MB", 1048576, "1.0 MB"},
		{"megabytes", 1572864, "1.5 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanSize(tt.bytes); got != tt.want {
				t.Errorf("humanSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestSummaryValue(t *testing.T) {
	tests := []struct {
		name, timeframe, want string
	}{
		{"7d", "7d", "10"},
		{"30d", "30d", "20"},
		{"90d", "90d", "30"},
		{"unknown", "1y", "20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := summaryValue(tt.timeframe, 10, 20, 30); got != tt.want {
				t.Errorf("summaryValue(%q) = %q, want %q", tt.timeframe, got, tt.want)
			}
		})
	}
}

func TestTimeframeLabel(t *testing.T) {
	tests := []struct {
		name, timeframe, want string
	}{
		{"7d", "7d", "7 day"},
		{"90d", "90d", "90 day"},
		{"30d", "30d", "30 day"},
		{"unknown", "1y", "30 day"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeframeLabel(tt.timeframe); got != tt.want {
				t.Errorf("timeframeLabel(%q) = %q, want %q", tt.timeframe, got, tt.want)
			}
		})
	}
}

func TestPreviousPage(t *testing.T) {
	tests := []struct {
		name string
		page int
		want any
	}{
		{"page 1", 1, nil},
		{"page 0", 0, nil},
		{"page 2", 2, 1},
		{"page 5", 5, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := previousPage(tt.page); got != tt.want {
				t.Errorf("previousPage(%d) = %v, want %v", tt.page, got, tt.want)
			}
		})
	}
}

func TestHasResetHeader(t *testing.T) {
	tests := []struct {
		name   string
		header string
		key    string
		want   bool
	}{
		{"absent", "", "activityFeed", false},
		{"single match", "activityFeed", "activityFeed", true},
		{"csv match", "foo,activityFeed,bar", "activityFeed", true},
		{"csv no match", "foo,bar,baz", "activityFeed", false},
		{"whitespace match", "foo, activityFeed , bar", "activityFeed", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/x", nil)

			if tt.header != "" {
				r.Header.Set(httpx.HeaderReset, tt.header)
			}

			if got := hasResetHeader(r, tt.key); got != tt.want {
				t.Errorf("hasResetHeader(%q, %q) = %v, want %v", tt.header, tt.key, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// B. Flash cookie functions
// ---------------------------------------------------------------------------

func TestFlashRoundTrip(t *testing.T) {
	t.Run("set and consume", func(t *testing.T) {
		w := httptest.NewRecorder()
		setFlash(w, flashPayload{Kind: "success", Title: "Done", Message: "It worked"})

		resp := w.Result()

		var flashCookie *http.Cookie

		for _, c := range resp.Cookies() {
			if c.Name == flashCookieName {
				flashCookie = c

				break
			}
		}

		if flashCookie == nil {
			t.Fatal("flash cookie not set")
		}

		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		r.AddCookie(flashCookie)
		w2 := httptest.NewRecorder()

		flash := consumeFlash(w2, r)

		if flash == nil {
			t.Fatal("consumeFlash returned nil")
		}

		if flash["kind"] != "success" {
			t.Errorf("kind = %v, want %q", flash["kind"], "success")
		}

		if flash["title"] != "Done" {
			t.Errorf("title = %v, want %q", flash["title"], "Done")
		}

		// verify cookie is cleared
		for _, c := range w2.Result().Cookies() {
			if c.Name == flashCookieName && c.MaxAge != -1 {
				t.Errorf("flash cookie MaxAge = %d, want -1", c.MaxAge)
			}
		}
	})

	t.Run("no cookie returns nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		w := httptest.NewRecorder()

		if got := consumeFlash(w, r); got != nil {
			t.Errorf("consumeFlash = %v, want nil", got)
		}
	})

	t.Run("invalid JSON returns nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		r.AddCookie(&http.Cookie{Name: flashCookieName, Value: url.QueryEscape("not-json")})
		w := httptest.NewRecorder()

		if got := consumeFlash(w, r); got != nil {
			t.Errorf("consumeFlash = %v, want nil", got)
		}
	})

	t.Run("invalid URL encoding returns nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		r.AddCookie(&http.Cookie{Name: flashCookieName, Value: "%zz"})
		w := httptest.NewRecorder()

		if got := consumeFlash(w, r); got != nil {
			t.Errorf("consumeFlash = %v, want nil", got)
		}
	})
}

func TestRenderDashboardWithFlash(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	// Set a flash cookie, then render a page — verify flash is consumed
	w := httptest.NewRecorder()
	setFlash(w, flashPayload{Kind: "info", Title: "Test", Message: "flash message"})
	resp := w.Result()

	var flashCookie *http.Cookie

	for _, c := range resp.Cookies() {
		if c.Name == flashCookieName {
			flashCookie = c

			break
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboard/navigation", nil)
	req.RequestURI = "/dashboard/navigation"
	req.AddCookie(flashCookie)

	page := assert.AssertFromHandler(t, testInertia, navigationHandler, req)
	page.AssertHasProp(t, "flash")
}

// ---------------------------------------------------------------------------
// C. Data generator functions
// ---------------------------------------------------------------------------

func TestLedgerRows(t *testing.T) {
	rows := ledgerRows("keep")

	if len(rows) != 18 {
		t.Fatalf("len = %d, want 18", len(rows))
	}

	if rows[0]["id"] != "inv_001" {
		t.Errorf("first id = %v, want %q", rows[0]["id"], "inv_001")
	}

	note, _ := rows[0]["note"].(string)

	if !strings.Contains(note, "keep") {
		t.Errorf("note = %q, want to contain %q", note, "keep")
	}
}

func TestSignalChunk(t *testing.T) {
	tests := []struct {
		name    string
		page    int
		wantLen int
		firstID string
	}{
		{"page 1", 1, 3, "sig_01"},
		{"page 2", 2, 3, "sig_04"},
		{"page 3", 3, 3, "sig_07"},
		{"page 0", 0, 0, ""},
		{"page 4", 4, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := signalChunk(tt.page, "30d")

			if len(chunk) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(chunk), tt.wantLen)
			}

			if tt.firstID != "" && chunk[0]["id"] != tt.firstID {
				t.Errorf("first id = %v, want %q", chunk[0]["id"], tt.firstID)
			}
		})
	}

	t.Run("timeframe affects detail", func(t *testing.T) {
		c7 := signalChunk(1, "7d")
		c90 := signalChunk(1, "90d")

		if c7[0]["detail"] == c90[0]["detail"] {
			t.Error("expected different detail text for different timeframes")
		}
	})
}

func TestCapacityModel(t *testing.T) {
	baseline := capacityModel("baseline")
	burst := capacityModel("burst")

	routing := baseline["routing"].(map[string]any)

	if routing["coverage"] != "96.8%" {
		t.Errorf("baseline coverage = %v, want %q", routing["coverage"], "96.8%")
	}

	burstRouting := burst["routing"].(map[string]any)

	if burstRouting["coverage"] != "99.1%" {
		t.Errorf("burst coverage = %v, want %q", burstRouting["coverage"], "99.1%")
	}
}

func TestFeedChunk(t *testing.T) {
	t.Run("page 1 all", func(t *testing.T) {
		items, next, label := feedChunk(1, "all", "all")

		if len(items) != 4 {
			t.Fatalf("len = %d, want 4", len(items))
		}

		if next != 2 {
			t.Errorf("nextPage = %v, want 2", next)
		}

		if !strings.Contains(label, "12") {
			t.Errorf("label = %q, want to contain %q", label, "12")
		}
	})

	t.Run("page 3 is last", func(t *testing.T) {
		items, next, _ := feedChunk(3, "all", "all")

		if len(items) != 4 {
			t.Fatalf("len = %d, want 4", len(items))
		}

		if next != nil {
			t.Errorf("nextPage = %v, want nil", next)
		}
	})

	t.Run("team filter", func(t *testing.T) {
		items, _, _ := feedChunk(1, "sales", "all")

		for _, item := range items {
			if item["team"] != "sales" {
				t.Errorf("item team = %v, want %q", item["team"], "sales")
			}
		}
	})

	t.Run("kind filter", func(t *testing.T) {
		items, _, _ := feedChunk(1, "all", "invoice")

		for _, item := range items {
			if item["kind"] != "invoice" {
				t.Errorf("item kind = %v, want %q", item["kind"], "invoice")
			}
		}
	})

	t.Run("out of range", func(t *testing.T) {
		items, next, _ := feedChunk(99, "all", "all")

		if items != nil {
			t.Errorf("items = %v, want nil", items)
		}

		if next != nil {
			t.Errorf("nextPage = %v, want nil", next)
		}
	})

	t.Run("nonexistent team", func(t *testing.T) {
		items, next, _ := feedChunk(1, "nonexistent", "all")

		if items != nil {
			t.Errorf("items = %v, want nil", items)
		}

		if next != nil {
			t.Errorf("nextPage = %v, want nil", next)
		}
	})
}

func TestStateTasks(t *testing.T) {
	tasks := stateTasks()

	if len(tasks) != 24 {
		t.Fatalf("len = %d, want 24", len(tasks))
	}

	if tasks[0]["id"] != "task_01" {
		t.Errorf("first id = %v, want %q", tasks[0]["id"], "task_01")
	}
}

func TestFormsProps(t *testing.T) {
	setupTestDB(t)

	fp := formsProps()

	for _, key := range []string{"pageTitle", "pageSubtitle", "inviteRoles", "recentInvites", "uploadedFiles", "approvalSummary", "recentApprovals"} {
		if _, ok := fp[key]; !ok {
			t.Errorf("missing key %q", key)
		}
	}
}

// ---------------------------------------------------------------------------
// D. Inertia handler tests
// ---------------------------------------------------------------------------

func TestRootHandler(t *testing.T) {
	testMux := newDashboardTestMux(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusConflict {
		// Inertia redirect may use 302 or 409 depending on X-Inertia header
		if loc := w.Header().Get("Location"); loc == "" {
			t.Errorf("expected redirect, got status %d with no Location header", w.Code)
		}
	}
}

func TestNavigationHandler(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("default pane", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/navigation", nil)
		req.RequestURI = "/dashboard/navigation"

		page := assert.AssertFromHandler(t, testInertia, navigationHandler, req)
		page.AssertComponent(t, "Dashboard/Navigation")
		page.AssertPropEquals(t, "pane", "links")
		page.AssertHasProp(t, "visitMatrix")
		page.AssertHasProp(t, "visitTargets")
	})

	t.Run("custom pane", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/navigation?pane=manual", nil)
		req.RequestURI = "/dashboard/navigation?pane=manual"

		page := assert.AssertFromHandler(t, testInertia, navigationHandler, req)
		page.AssertPropEquals(t, "pane", "manual")
	})
}

func TestScrollHandler(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("default preserve", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/scroll", nil)
		req.RequestURI = "/dashboard/scroll"

		page := assert.AssertFromHandler(t, testInertia, scrollHandler, req)
		page.AssertComponent(t, "Dashboard/Scroll")
		page.AssertPropEquals(t, "preserve", "keep")
		page.AssertHasProp(t, "layoutNotes")

		ledger, ok := page.Props["scrollLedger"].([]any)

		if !ok {
			t.Fatal("scrollLedger not a slice")
		}

		if len(ledger) != 18 {
			t.Errorf("scrollLedger len = %d, want 18", len(ledger))
		}
	})

	t.Run("custom preserve", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/scroll?preserve=reset", nil)
		req.RequestURI = "/dashboard/scroll?preserve=reset"

		page := assert.AssertFromHandler(t, testInertia, scrollHandler, req)
		page.AssertPropEquals(t, "preserve", "reset")
	})
}

func TestRedirectsHandler(t *testing.T) {
	testInertia := newDashboardTestInertia(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard/redirects", nil)
	req.RequestURI = "/dashboard/redirects"

	page := assert.AssertFromHandler(t, testInertia, redirectsHandler, req)
	page.AssertComponent(t, "Dashboard/Redirects")
	page.AssertHasProp(t, "redirectNotes")
	page.AssertHasProp(t, "locationNotes")
	page.AssertPropEquals(t, "pageTitle", "Redirects and location visits")
}

func TestFormsHandler(t *testing.T) {
	testInertia := newDashboardTestInertia(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard/forms", nil)
	req.RequestURI = "/dashboard/forms"

	page := assert.AssertFromHandler(t, testInertia, formsHandler, req)
	page.AssertComponent(t, "Dashboard/Forms")
	page.AssertHasProp(t, "pageTitle")
	page.AssertHasProp(t, "inviteRoles")
	page.AssertHasProp(t, "recentInvites")
	page.AssertHasProp(t, "uploadedFiles")
	page.AssertHasProp(t, "approvalSummary")
}

func TestDataHandler(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("defaults", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data", nil)
		req.RequestURI = "/dashboard/data"

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)
		page.AssertComponent(t, "Dashboard/Data")
		page.AssertPropEquals(t, "timeframe", "30d")
		page.AssertPropEquals(t, "capacityMode", "baseline")
		page.AssertPropEquals(t, "signalsPage", float64(1))
		page.AssertPropEquals(t, "signalsHasMore", true)
		page.AssertHasProp(t, "accountSummary")
		page.AssertHasProp(t, "liveMetrics")
		page.AssertHasProp(t, "prefetchTargets")
	})

	t.Run("custom timeframe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data?timeframe=7d", nil)
		req.RequestURI = "/dashboard/data?timeframe=7d"

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)
		page.AssertPropEquals(t, "timeframe", "7d")
	})

	t.Run("last signal page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data?signalPage=3", nil)
		req.RequestURI = "/dashboard/data?signalPage=3"

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)
		page.AssertPropEquals(t, "signalsHasMore", false)
	})

	t.Run("deferred props", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data", nil)
		req.RequestURI = "/dashboard/data"

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)

		insights, ok := page.DeferredProps["insights"]

		if !ok {
			t.Fatal("missing deferred group 'insights'")
		}

		hasNarrative := false
		hasBreakdown := false

		for _, prop := range insights {
			if prop == "trendNarrative" {
				hasNarrative = true
			}

			if prop == "trendBreakdown" {
				hasBreakdown = true
			}
		}

		if !hasNarrative {
			t.Error("deferred group missing 'trendNarrative'")
		}

		if !hasBreakdown {
			t.Error("deferred group missing 'trendBreakdown'")
		}
	})

	t.Run("once props", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data", nil)
		req.RequestURI = "/dashboard/data"

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)

		if _, ok := page.OnceProps["releaseNotes"]; !ok {
			t.Error("missing onceProps 'releaseNotes'")
		}
	})
}

func TestDataHandlerMergeProps(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("signals merge on page 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data?signalPage=2", nil)
		req.RequestURI = "/dashboard/data?signalPage=2"
		req.Header.Set(httpx.HeaderInertia, "true")

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)

		found := false

		for _, p := range page.MergeProps {
			if p == "signals" {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("MergeProps = %v, want to contain %q", page.MergeProps, "signals")
		}
	})

	t.Run("capacity deep merge on burst", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/data?capacityMode=burst", nil)
		req.RequestURI = "/dashboard/data?capacityMode=burst"
		req.Header.Set(httpx.HeaderInertia, "true")

		page := assert.AssertFromHandler(t, testInertia, dataHandler, req)

		found := false

		for _, p := range page.DeepMergeProps {
			if p == "capacityModel" {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("DeepMergeProps = %v, want to contain %q", page.DeepMergeProps, "capacityModel")
		}
	})
}

func TestDataHandlerPartialReload(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/data", nil)
	req.RequestURI = "/dashboard/data"
	req.Header.Set(httpx.HeaderInertia, "true")
	req.Header.Set(httpx.HeaderPartialComponent, "Dashboard/Data")
	req.Header.Set(httpx.HeaderPartialData, "trendNarrative,trendBreakdown,auditTrail")

	page := assert.AssertFromHandler(t, testInertia, dataHandler, req)
	page.AssertComponent(t, "Dashboard/Data")
	page.AssertHasProp(t, "trendNarrative")
	page.AssertHasProp(t, "trendBreakdown")
}

func TestStateHandler(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("default mode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/state", nil)
		req.RequestURI = "/dashboard/state"

		page := assert.AssertFromHandler(t, testInertia, stateHandler, req)
		page.AssertComponent(t, "Dashboard/State")
		page.AssertPropEquals(t, "historyMode", "default")

		if page.EncryptHistory {
			t.Error("EncryptHistory should be false")
		}

		if page.ClearHistory {
			t.Error("ClearHistory should be false")
		}

		page.AssertHasProp(t, "playbooks")
		page.AssertHasProp(t, "errorLinks")

		tasks, ok := page.Props["longTasks"].([]any)

		if !ok {
			t.Fatal("longTasks not a slice")
		}

		if len(tasks) != 24 {
			t.Errorf("longTasks len = %d, want 24", len(tasks))
		}
	})

	t.Run("encrypted mode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/state?mode=encrypted", nil)
		req.RequestURI = "/dashboard/state?mode=encrypted"

		page := assert.AssertFromHandler(t, testInertia, stateHandler, req)

		if !page.EncryptHistory {
			t.Error("EncryptHistory should be true")
		}
	})

	t.Run("clear mode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/state?mode=clear", nil)
		req.RequestURI = "/dashboard/state?mode=clear"

		page := assert.AssertFromHandler(t, testInertia, stateHandler, req)

		if !page.ClearHistory {
			t.Error("ClearHistory should be true")
		}
	})

	t.Run("once props", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/state", nil)
		req.RequestURI = "/dashboard/state"

		page := assert.AssertFromHandler(t, testInertia, stateHandler, req)

		if _, ok := page.OnceProps["releaseTrack"]; !ok {
			t.Error("missing onceProps 'releaseTrack'")
		}
	})
}

func TestFeedHandlerFilters(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("team filter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/feed?team=sales", nil)
		req.RequestURI = "/dashboard/feed?team=sales"

		page := assert.AssertFromHandler(t, testInertia, feedHandler, req)
		summary, ok := page.Props["feedSummary"].(map[string]any)

		if !ok {
			t.Fatal("feedSummary missing")
		}

		if summary["teamLabel"] != "Sales" {
			t.Errorf("teamLabel = %v, want %q", summary["teamLabel"], "Sales")
		}
	})

	t.Run("kind filter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/feed?kind=invoice", nil)
		req.RequestURI = "/dashboard/feed?kind=invoice"

		page := assert.AssertFromHandler(t, testInertia, feedHandler, req)
		summary := page.Props["feedSummary"].(map[string]any)

		if summary["kindLabel"] != "Invoice" {
			t.Errorf("kindLabel = %v, want %q", summary["kindLabel"], "Invoice")
		}
	})

	t.Run("page 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/feed?feedPage=2", nil)
		req.RequestURI = "/dashboard/feed?feedPage=2"

		page := assert.AssertFromHandler(t, testInertia, feedHandler, req)
		summary := page.Props["feedSummary"].(map[string]any)

		if summary["page"] != float64(2) {
			t.Errorf("page = %v, want 2", summary["page"])
		}
	})
}

func TestFeedHandlerMergeAndReset(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("merge on page 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/feed?feedPage=2", nil)
		req.RequestURI = "/dashboard/feed?feedPage=2"
		req.Header.Set(httpx.HeaderInertia, "true")

		page := assert.AssertFromHandler(t, testInertia, feedHandler, req)

		found := false

		for _, p := range page.MergeProps {
			if p == "activityFeed" {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("MergeProps = %v, want to contain %q", page.MergeProps, "activityFeed")
		}
	})

	t.Run("reset header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/feed", nil)
		req.RequestURI = "/dashboard/feed"
		req.Header.Set(httpx.HeaderInertia, "true")
		req.Header.Set(httpx.HeaderReset, "activityFeed")

		page := assert.AssertFromHandler(t, testInertia, feedHandler, req)
		scroll, ok := page.ScrollProps["activityFeed"]

		if !ok {
			t.Fatal("missing scrollProps.activityFeed")
		}

		if !scroll.Reset {
			t.Error("reset should be true when X-Inertia-Reset includes activityFeed")
		}
	})
}

// ---------------------------------------------------------------------------
// E. POST / non-Inertia handler tests
// ---------------------------------------------------------------------------

func TestUploadHandlerPrecognition(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	before := database.UploadCount(db)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("label", "Test upload")
	part, _ := writer.CreateFormFile("file", "report.csv")
	part.Write([]byte("data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(httpx.HeaderPrecognition, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	if after := database.UploadCount(db); after != before {
		t.Fatalf("upload count = %d, want %d (no mutation on precognition)", after, before)
	}
}

func TestEscalateHandlerPrecognition(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	before := database.GetCounter(db, "priority_escalations")

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/escalate", nil)
	req.Header.Set(httpx.HeaderPrecognition, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	if after := database.GetCounter(db, "priority_escalations"); after != before {
		t.Fatalf("counter = %d, want %d", after, before)
	}
}

func TestRedirectsRedirectHandlerPrecognition(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	req := httptest.NewRequest(http.MethodPost, "/dashboard/redirects/redirect", nil)
	req.Header.Set(httpx.HeaderPrecognition, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRedirectsLocationHandlerPrecognition(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	req := httptest.NewRequest(http.MethodPost, "/dashboard/redirects/location", nil)
	req.Header.Set(httpx.HeaderPrecognition, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestUploadHandlerValidation(t *testing.T) {
	testInertia := newDashboardTestInertia(t)

	t.Run("missing both", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/upload", body)
		req.RequestURI = "/dashboard/forms/upload"
		req.Header.Set("Content-Type", writer.FormDataContentType())

		page := assert.AssertFromHandler(t, testInertia, uploadHandler, req)
		page.AssertComponent(t, "Dashboard/Forms")

		errors, ok := page.Props["errors"].(map[string]any)

		if !ok {
			t.Fatal("errors prop missing")
		}

		if errors["label"] == nil {
			t.Error("expected label validation error")
		}

		if errors["file"] == nil {
			t.Error("expected file validation error")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("label", "Test upload")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/upload", body)
		req.RequestURI = "/dashboard/forms/upload"
		req.Header.Set("Content-Type", writer.FormDataContentType())

		page := assert.AssertFromHandler(t, testInertia, uploadHandler, req)
		errors := page.Props["errors"].(map[string]any)

		if errors["file"] == nil {
			t.Error("expected file validation error")
		}

		if errors["label"] != nil {
			t.Error("unexpected label validation error")
		}
	})

	t.Run("missing label", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("hello"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/upload", body)
		req.RequestURI = "/dashboard/forms/upload"
		req.Header.Set("Content-Type", writer.FormDataContentType())

		page := assert.AssertFromHandler(t, testInertia, uploadHandler, req)
		errors := page.Props["errors"].(map[string]any)

		if errors["label"] == nil {
			t.Error("expected label validation error")
		}

		if errors["file"] != nil {
			t.Error("unexpected file validation error")
		}
	})
}

func TestUploadHandlerSuccess(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	before := database.UploadCount(db)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("label", "Test upload")
	part, _ := writer.CreateFormFile("file", "report.csv")
	part.Write([]byte("col1,col2\nval1,val2"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(httpx.HeaderInertia, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard/forms" {
		t.Errorf("location = %q, want %q", got, "/dashboard/forms")
	}

	if after := database.UploadCount(db); after != before+1 {
		t.Errorf("upload count = %d, want %d", after, before+1)
	}
}

func TestEscalateHandler(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	beforeCounter := database.GetCounter(db, "priority_escalations")
	beforeApprovals := database.ApprovalCount(db)

	req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/escalate", nil)
	req.Header.Set(httpx.HeaderInertia, "true")
	req.Header.Set("X-CSRF-TOKEN", rawToken)
	req.AddCookie(csrfCookie)
	w := httptest.NewRecorder()

	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if got := w.Header().Get("Location"); got != "/dashboard/forms" {
		t.Errorf("location = %q, want %q", got, "/dashboard/forms")
	}

	if after := database.GetCounter(db, "priority_escalations"); after != beforeCounter+1 {
		t.Errorf("counter = %d, want %d", after, beforeCounter+1)
	}

	if after := database.ApprovalCount(db); after != beforeApprovals+1 {
		t.Errorf("approval count = %d, want %d", after, beforeApprovals+1)
	}
}

func TestHttpPreviewHandler(t *testing.T) {
	tests := []struct {
		name            string
		seats, tier     string
		wantMonthly     float64
		wantRecommended bool
	}{
		{"default growth", "12", "growth", 336, false},
		{"enterprise", "10", "enterprise", 420, false},
		{"starter", "5", "starter", 80, false},
		{"recommended", "25", "growth", 700, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(url.Values{
				"seats": {tt.seats},
				"tier":  {tt.tier},
			}.Encode())

			req := httptest.NewRequest(http.MethodPost, "/dashboard/forms/http-preview", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			httpPreviewHandler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", w.Code)
			}

			var payload map[string]any

			if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
				t.Fatalf("decode: %v", err)
			}

			if payload["monthlyEstimate"] != tt.wantMonthly {
				t.Errorf("monthlyEstimate = %v, want %v", payload["monthlyEstimate"], tt.wantMonthly)
			}

			if payload["recommended"] != tt.wantRecommended {
				t.Errorf("recommended = %v, want %v", payload["recommended"], tt.wantRecommended)
			}
		})
	}

	t.Run("GET not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/forms/http-preview", nil)
		w := httptest.NewRecorder()

		httpPreviewHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestRedirectsRedirectHandler(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	t.Run("POST redirects", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dashboard/redirects/redirect", nil)
		req.Header.Set(httpx.HeaderInertia, "true")
		req.Header.Set("X-CSRF-TOKEN", rawToken)
		req.AddCookie(csrfCookie)
		w := httptest.NewRecorder()

		testMux.ServeHTTP(w, req)

		if w.Code != http.StatusFound && w.Code != http.StatusSeeOther {
			t.Fatalf("status = %d, want 302 or 303", w.Code)
		}

		if got := w.Header().Get("Location"); got != "/dashboard/redirects" {
			t.Errorf("location = %q, want %q", got, "/dashboard/redirects")
		}
	})

	t.Run("GET not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/redirects/redirect", nil)
		w := httptest.NewRecorder()

		testMux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestRedirectsLocationHandler(t *testing.T) {
	testMux := newDashboardTestMux(t)
	csrfCookie, rawToken := issueDashboardCSRFCookie(t, testMux)

	t.Run("POST sets location header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dashboard/redirects/location", nil)
		req.Header.Set(httpx.HeaderInertia, "true")
		req.Header.Set("X-CSRF-TOKEN", rawToken)
		req.AddCookie(csrfCookie)
		w := httptest.NewRecorder()

		testMux.ServeHTTP(w, req)

		if loc := w.Header().Get(httpx.HeaderLocation); loc == "" {
			t.Error("expected X-Inertia-Location header to be set")
		}
	})

	t.Run("GET not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard/redirects/location", nil)
		w := httptest.NewRecorder()

		testMux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestStateErrorHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{"404", "/dashboard/state/error?code=404", http.StatusNotFound},
		{"500", "/dashboard/state/error?code=500", http.StatusInternalServerError},
		{"default", "/dashboard/state/error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			stateErrorHandler(w, req)

			if w.Code != tt.want {
				t.Errorf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}
