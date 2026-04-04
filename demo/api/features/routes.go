package features

import "net/http"

// RegisterRoutes mounts all feature showcase HTTP routes onto the provided mux.
func RegisterRoutes(mux *http.ServeMux, deps Deps) {
	a := newApp(deps)

	auth := func(h http.HandlerFunc) http.Handler {
		return deps.RequireAuth(h)
	}

	// Forms
	mux.Handle("/features/forms/use-form", auth(a.useFormHandler))
	mux.Handle("/features/forms/form-component", auth(a.formComponentHandler))
	mux.Handle("/features/forms/file-uploads", auth(a.fileUploadsHandler))
	mux.Handle("/features/forms/validation", auth(a.validationHandler))
	mux.Handle("/features/forms/validation/secondary", auth(a.validationSecondaryHandler))
	mux.Handle("/features/forms/precognition", auth(a.precognitionHandler))
	mux.Handle("/features/forms/optimistic-updates", auth(a.optimisticUpdatesHandler))
	mux.Handle("POST /features/forms/optimistic-toggle/{id}", auth(a.optimisticToggleHandler))
	mux.Handle("/features/forms/use-form-context", auth(a.formContextHandler))
	mux.Handle("/features/forms/dotted-keys", auth(a.dottedKeysHandler))
	mux.Handle("/features/forms/wayfinder", auth(a.wayfinderHandler))

	// Navigation
	mux.Handle("/features/navigation/links", auth(a.linksHandler))
	mux.Handle("/features/navigation/preserve-state", auth(a.preserveStateHandler))
	mux.Handle("/features/navigation/preserve-scroll", auth(a.preserveScrollHandler))
	mux.Handle("/features/navigation/view-transitions", auth(a.viewTransitionsHandler))
	mux.Handle("/features/navigation/history-management", auth(a.historyManagementHandler))
	mux.Handle("/features/navigation/async-requests", auth(a.asyncRequestsHandler))
	mux.Handle("/features/navigation/async-slow", auth(a.asyncSlowHandler))
	mux.Handle("/features/navigation/manual-visits", auth(a.manualVisitsHandler))
	mux.Handle("/features/navigation/redirects", auth(a.redirectsHandler))
	mux.Handle("POST /features/navigation/redirects/{action}", auth(a.redirectsActionHandler))
	mux.Handle("/features/navigation/scroll-management", auth(a.scrollManagementHandler))
	mux.Handle("/features/navigation/instant-visits", auth(a.instantVisitsHandler))
	mux.Handle("/features/navigation/instant-visit-target", auth(a.instantVisitTargetHandler))
	mux.Handle("/features/navigation/url-fragments", auth(a.urlFragmentsHandler))
	mux.Handle("POST /features/navigation/url-fragments/{action}", auth(a.urlFragmentsActionHandler))

	// Data Loading
	mux.Handle("/features/data-loading/deferred-props", auth(a.deferredPropsHandler))
	mux.Handle("/features/data-loading/partial-reloads", auth(a.partialReloadsHandler))
	mux.Handle("/features/data-loading/infinite-scroll", auth(a.infiniteScrollHandler))
	mux.Handle("/features/data-loading/when-visible", auth(a.whenVisibleHandler))
	mux.Handle("/features/data-loading/polling", auth(a.pollingHandler))
	mux.Handle("/features/data-loading/prop-merging", auth(a.propMergingHandler))
	mux.Handle("/features/data-loading/optional-props", auth(a.optionalPropsHandler))
	mux.Handle("/features/data-loading/once-props/{page}", auth(a.oncePropsHandler))

	// Prefetching
	mux.Handle("/features/prefetching/link-prefetch", auth(a.linkPrefetchHandler))
	mux.Handle("/features/prefetching/stale-while-revalidate", auth(a.staleWhileRevalidateHandler))
	mux.Handle("/features/prefetching/manual-prefetch", auth(a.manualPrefetchHandler))
	mux.Handle("/features/prefetching/cache-management", auth(a.cacheManagementHandler))

	// State
	mux.Handle("/features/state/remember", auth(a.rememberHandler))
	mux.Handle("/features/state/flash-data", auth(a.flashDataHandler))
	mux.Handle("POST /features/state/flash-data/{action}", auth(a.flashDataActionHandler))
	mux.Handle("/features/state/shared-props", auth(a.sharedPropsHandler))

	// Layouts
	mux.Handle("/features/layouts/persistent-layouts", auth(a.persistentLayoutsHandler))
	mux.Handle("/features/layouts/persistent-layouts/page-2", auth(a.persistentLayoutsPage2Handler))
	mux.Handle("/features/layouts/nested-layouts", auth(a.nestedLayoutsHandler))
	mux.Handle("/features/layouts/head", auth(a.headHandler))
	mux.Handle("/features/layouts/layout-props", auth(a.layoutPropsHandler))

	// Events
	mux.Handle("/features/events/global-events", auth(a.globalEventsHandler))
	mux.Handle("/features/events/visit-callbacks", auth(a.visitCallbacksHandler))
	mux.Handle("/features/events/progress", auth(a.progressHandler))
	mux.Handle("/features/events/progress/slow", auth(a.progressSlowHandler))

	// Errors
	mux.Handle("/features/errors/http-exceptions", auth(a.httpExceptionsHandler))
	mux.Handle("/features/errors/http-exceptions/{code}", auth(a.httpExceptionsTriggerHandler))
	mux.Handle("/features/errors/network-errors", auth(a.networkErrorsHandler))

	// HTTP
	mux.Handle("/features/http/use-http", auth(a.useHttpHandler))
	mux.Handle("/features/http/use-http/api", auth(a.useHttpApiHandler))
}
