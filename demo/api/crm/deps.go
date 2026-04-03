package crm

import (
	"database/sql"
	"net/http"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

// Deps contains the host application integrations required by the CRM package.
type Deps struct {
	DB                *sql.DB
	RequireAuth       func(http.Handler) http.Handler
	Render            func(http.ResponseWriter, *http.Request, string, httpx.Props)
	RenderWithContext func(http.ResponseWriter, *http.Request, string, httpx.Props)
	Redirect          func(http.ResponseWriter, *http.Request, string)
	RouteURL          func(string, map[string]string) string
	SetFlash          func(http.ResponseWriter, flash.Message)
	CurrentUser       func(*http.Request) *database.User
}
