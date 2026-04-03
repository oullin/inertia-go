package auth

import (
	"database/sql"
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
)

// Flash carries the shared flash payload written by the host app.
type Flash struct {
	Kind    string
	Title   string
	Message string
}

// Deps contains the host application integrations required by the auth package.
type Deps struct {
	DB                *sql.DB
	Render            func(http.ResponseWriter, *http.Request, string, httpx.Props)
	RenderWithContext func(http.ResponseWriter, *http.Request, string, httpx.Props)
	Redirect          func(http.ResponseWriter, *http.Request, string)
	RouteURL          func(string, map[string]string) string
	SetFlash          func(http.ResponseWriter, Flash)
}
