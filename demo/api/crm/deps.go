package crm

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

// Flash carries the shared flash payload written by the host app.
type Flash struct {
	Kind    string
	Title   string
	Message string
}

// Deps contains the host application integrations required by the CRM package.
type Deps struct {
	Repository        Repository
	RequireAuth       func(http.Handler) http.Handler
	Render            func(http.ResponseWriter, *http.Request, string, httpx.Props)
	RenderWithContext func(http.ResponseWriter, *http.Request, string, httpx.Props)
	Redirect          func(http.ResponseWriter, *http.Request, string)
	RouteURL          func(string, map[string]string) string
	SetFlash          func(http.ResponseWriter, Flash)
	CurrentUser       func(*http.Request) *database.User
}
