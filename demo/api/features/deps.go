package features

import (
	"database/sql"
	"net/http"

	"github.com/oullin/inertia-go/core/flash"
	"github.com/oullin/inertia-go/core/httpx"
)

// Deps contains the host application integrations required by the features package.
type Deps struct {
	DB          *sql.DB
	RequireAuth func(http.Handler) http.Handler
	Render      func(http.ResponseWriter, *http.Request, string, httpx.Props)
	Redirect    func(http.ResponseWriter, *http.Request, string)
	Location    func(http.ResponseWriter, *http.Request, string)
	RouteURL    func(string, map[string]string) string
	SetFlash    func(http.ResponseWriter, flash.Message)
}
