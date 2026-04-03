package features

import "net/http"

// RegisterRoutes mounts all feature showcase HTTP routes onto the provided mux.
func RegisterRoutes(mux *http.ServeMux, deps Deps) {
	a := newApp(deps)

	// Forms
	mux.Handle("/features/forms/use-form", deps.RequireAuth(http.HandlerFunc(a.useFormHandler)))
	mux.Handle("/features/forms/form-component", deps.RequireAuth(http.HandlerFunc(a.formComponentHandler)))
	mux.Handle("/features/forms/file-uploads", deps.RequireAuth(http.HandlerFunc(a.fileUploadsHandler)))
	mux.Handle("/features/forms/validation", deps.RequireAuth(http.HandlerFunc(a.validationHandler)))
	mux.Handle("/features/forms/validation/secondary", deps.RequireAuth(http.HandlerFunc(a.validationSecondaryHandler)))
	mux.Handle("/features/forms/precognition", deps.RequireAuth(http.HandlerFunc(a.precognitionHandler)))
	mux.Handle("/features/forms/optimistic-updates", deps.RequireAuth(http.HandlerFunc(a.optimisticUpdatesHandler)))
	mux.Handle("/features/forms/optimistic-toggle/", deps.RequireAuth(http.HandlerFunc(a.optimisticToggleHandler)))
	mux.Handle("/features/forms/use-form-context", deps.RequireAuth(http.HandlerFunc(a.formContextHandler)))
	mux.Handle("/features/forms/dotted-keys", deps.RequireAuth(http.HandlerFunc(a.dottedKeysHandler)))
	mux.Handle("/features/forms/wayfinder", deps.RequireAuth(http.HandlerFunc(a.wayfinderHandler)))

	// Navigation
	mux.Handle("/features/navigation/links", deps.RequireAuth(http.HandlerFunc(a.linksHandler)))
	mux.Handle("/features/navigation/preserve-state", deps.RequireAuth(http.HandlerFunc(a.preserveStateHandler)))
	mux.Handle("/features/navigation/preserve-scroll", deps.RequireAuth(http.HandlerFunc(a.preserveScrollHandler)))
	mux.Handle("/features/navigation/view-transitions", deps.RequireAuth(http.HandlerFunc(a.viewTransitionsHandler)))
	mux.Handle("/features/navigation/history-management", deps.RequireAuth(http.HandlerFunc(a.historyManagementHandler)))
	mux.Handle("/features/navigation/async-requests", deps.RequireAuth(http.HandlerFunc(a.asyncRequestsHandler)))
	mux.Handle("/features/navigation/async-slow", deps.RequireAuth(http.HandlerFunc(a.asyncSlowHandler)))
	mux.Handle("/features/navigation/manual-visits", deps.RequireAuth(http.HandlerFunc(a.manualVisitsHandler)))
	mux.Handle("/features/navigation/redirects", deps.RequireAuth(http.HandlerFunc(a.redirectsHandler)))
	mux.Handle("/features/navigation/redirects/", deps.RequireAuth(http.HandlerFunc(a.redirectsActionHandler)))
	mux.Handle("/features/navigation/scroll-management", deps.RequireAuth(http.HandlerFunc(a.scrollManagementHandler)))
	mux.Handle("/features/navigation/instant-visits", deps.RequireAuth(http.HandlerFunc(a.instantVisitsHandler)))
	mux.Handle("/features/navigation/instant-visit-target", deps.RequireAuth(http.HandlerFunc(a.instantVisitTargetHandler)))
	mux.Handle("/features/navigation/url-fragments", deps.RequireAuth(http.HandlerFunc(a.urlFragmentsHandler)))
	mux.Handle("/features/navigation/url-fragments/", deps.RequireAuth(http.HandlerFunc(a.urlFragmentsActionHandler)))

	// Data Loading
	mux.Handle("/features/data-loading/deferred-props", deps.RequireAuth(http.HandlerFunc(a.deferredPropsHandler)))
	mux.Handle("/features/data-loading/partial-reloads", deps.RequireAuth(http.HandlerFunc(a.partialReloadsHandler)))
	mux.Handle("/features/data-loading/infinite-scroll", deps.RequireAuth(http.HandlerFunc(a.infiniteScrollHandler)))
	mux.Handle("/features/data-loading/when-visible", deps.RequireAuth(http.HandlerFunc(a.whenVisibleHandler)))
	mux.Handle("/features/data-loading/polling", deps.RequireAuth(http.HandlerFunc(a.pollingHandler)))
	mux.Handle("/features/data-loading/prop-merging", deps.RequireAuth(http.HandlerFunc(a.propMergingHandler)))
	mux.Handle("/features/data-loading/optional-props", deps.RequireAuth(http.HandlerFunc(a.optionalPropsHandler)))
	mux.Handle("/features/data-loading/once-props/", deps.RequireAuth(http.HandlerFunc(a.oncePropsHandler)))

	// Prefetching
	mux.Handle("/features/prefetching/link-prefetch", deps.RequireAuth(http.HandlerFunc(a.linkPrefetchHandler)))
	mux.Handle("/features/prefetching/stale-while-revalidate", deps.RequireAuth(http.HandlerFunc(a.staleWhileRevalidateHandler)))
	mux.Handle("/features/prefetching/manual-prefetch", deps.RequireAuth(http.HandlerFunc(a.manualPrefetchHandler)))
	mux.Handle("/features/prefetching/cache-management", deps.RequireAuth(http.HandlerFunc(a.cacheManagementHandler)))

	// State
	mux.Handle("/features/state/remember", deps.RequireAuth(http.HandlerFunc(a.rememberHandler)))
	mux.Handle("/features/state/flash-data", deps.RequireAuth(http.HandlerFunc(a.flashDataHandler)))
	mux.Handle("/features/state/flash-data/", deps.RequireAuth(http.HandlerFunc(a.flashDataActionHandler)))
	mux.Handle("/features/state/shared-props", deps.RequireAuth(http.HandlerFunc(a.sharedPropsHandler)))

	// Layouts
	mux.Handle("/features/layouts/persistent-layouts", deps.RequireAuth(http.HandlerFunc(a.persistentLayoutsHandler)))
	mux.Handle("/features/layouts/persistent-layouts/page-2", deps.RequireAuth(http.HandlerFunc(a.persistentLayoutsPage2Handler)))
	mux.Handle("/features/layouts/nested-layouts", deps.RequireAuth(http.HandlerFunc(a.nestedLayoutsHandler)))
	mux.Handle("/features/layouts/head", deps.RequireAuth(http.HandlerFunc(a.headHandler)))
	mux.Handle("/features/layouts/layout-props", deps.RequireAuth(http.HandlerFunc(a.layoutPropsHandler)))

	// Events
	mux.Handle("/features/events/global-events", deps.RequireAuth(http.HandlerFunc(a.globalEventsHandler)))
	mux.Handle("/features/events/visit-callbacks", deps.RequireAuth(http.HandlerFunc(a.visitCallbacksHandler)))
	mux.Handle("/features/events/progress", deps.RequireAuth(http.HandlerFunc(a.progressHandler)))
	mux.Handle("/features/events/progress/slow", deps.RequireAuth(http.HandlerFunc(a.progressSlowHandler)))

	// Errors
	mux.Handle("/features/errors/http-exceptions", deps.RequireAuth(http.HandlerFunc(a.httpExceptionsHandler)))
	mux.Handle("/features/errors/http-exceptions/", deps.RequireAuth(http.HandlerFunc(a.httpExceptionsTriggerHandler)))
	mux.Handle("/features/errors/network-errors", deps.RequireAuth(http.HandlerFunc(a.networkErrorsHandler)))

	// HTTP
	mux.Handle("/features/http/use-http", deps.RequireAuth(http.HandlerFunc(a.useHttpHandler)))
	mux.Handle("/features/http/use-http/api", deps.RequireAuth(http.HandlerFunc(a.useHttpApiHandler)))
}
