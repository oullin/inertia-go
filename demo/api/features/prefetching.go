package features

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) linkPrefetchHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Prefetching/LinkPrefetch", httpx.Props{})
}

func (a app) staleWhileRevalidateHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Prefetching/StaleWhileRevalidate", httpx.Props{})
}

func (a app) manualPrefetchHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Prefetching/ManualPrefetch", httpx.Props{})
}

func (a app) cacheManagementHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Prefetching/CacheManagement", httpx.Props{})
}
