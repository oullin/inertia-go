package features

import (
	"net/http"

	"github.com/oullin/inertia-go/core/wayfinder"
)

// RegisterRoutes mounts all feature showcase HTTP routes onto the provided mux.
func RegisterRoutes(routes *wayfinder.Registry, mux *http.ServeMux, container Container) {
	a := newApp(container)

	auth := func(h http.HandlerFunc) http.Handler {
		return container.RequireAuth(h)
	}

	// Forms
	routes.Handle("features.forms.use-form", auth(a.useFormHandler), mux)
	routes.Handle("features.forms.form-component", auth(a.formComponentHandler), mux)
	routes.Handle("features.forms.file-uploads", auth(a.fileUploadsHandler), mux)
	routes.Handle("features.forms.validation", auth(a.validationHandler), mux)
	mux.Handle("/features/forms/validation/secondary", auth(a.validationSecondaryHandler))
	routes.Handle("features.forms.precognition", auth(a.precognitionHandler), mux)
	routes.Handle("features.forms.optimistic-updates", auth(a.optimisticUpdatesHandler), mux)
	mux.Handle("POST /features/forms/optimistic-toggle/{id}", auth(a.optimisticToggleHandler))
	routes.Handle("features.forms.use-form-context", auth(a.formContextHandler), mux)
	routes.Handle("features.forms.dotted-keys", auth(a.dottedKeysHandler), mux)
	routes.Handle("features.forms.wayfinder", auth(a.wayfinderHandler), mux)

	// Navigation
	routes.Handle("features.navigation.links", auth(a.linksHandler), mux)
	routes.Handle("features.navigation.preserve-state", auth(a.preserveStateHandler), mux)
	routes.Handle("features.navigation.preserve-scroll", auth(a.preserveScrollHandler), mux)
	routes.Handle("features.navigation.view-transitions", auth(a.viewTransitionsHandler), mux)
	routes.Handle("features.navigation.history-management", auth(a.historyManagementHandler), mux)
	routes.Handle("features.navigation.async-requests", auth(a.asyncRequestsHandler), mux)
	routes.Handle("features.navigation.async-slow", auth(a.asyncSlowHandler), mux)
	routes.Handle("features.navigation.manual-visits", auth(a.manualVisitsHandler), mux)
	routes.Handle("features.navigation.redirects", auth(a.redirectsHandler), mux)
	mux.Handle("POST /features/navigation/redirects/{action}", auth(a.redirectsActionHandler))
	routes.Handle("features.navigation.scroll-management", auth(a.scrollManagementHandler), mux)
	routes.Handle("features.navigation.instant-visits", auth(a.instantVisitsHandler), mux)
	routes.Handle("features.navigation.instant-visit-target", auth(a.instantVisitTargetHandler), mux)
	routes.Handle("features.navigation.url-fragments", auth(a.urlFragmentsHandler), mux)
	mux.Handle("POST /features/navigation/url-fragments/{action}", auth(a.urlFragmentsActionHandler))

	// Data Loading
	routes.Handle("features.data-loading.deferred-props", auth(a.deferredPropsHandler), mux)
	routes.Handle("features.data-loading.partial-reloads", auth(a.partialReloadsHandler), mux)
	routes.Handle("features.data-loading.infinite-scroll", auth(a.infiniteScrollHandler), mux)
	routes.Handle("features.data-loading.when-visible", auth(a.whenVisibleHandler), mux)
	routes.Handle("features.data-loading.polling", auth(a.pollingHandler), mux)
	routes.Handle("features.data-loading.prop-merging", auth(a.propMergingHandler), mux)
	routes.Handle("features.data-loading.optional-props", auth(a.optionalPropsHandler), mux)
	routes.Handle("features.data-loading.once-props", auth(a.oncePropsHandler), mux)

	// Prefetching
	routes.Handle("features.prefetching.link-prefetch", auth(a.linkPrefetchHandler), mux)
	routes.Handle("features.prefetching.stale-while-revalidate", auth(a.staleWhileRevalidateHandler), mux)
	routes.Handle("features.prefetching.manual-prefetch", auth(a.manualPrefetchHandler), mux)
	routes.Handle("features.prefetching.cache-management", auth(a.cacheManagementHandler), mux)

	// State
	routes.Handle("features.state.remember", auth(a.rememberHandler), mux)
	routes.Handle("features.state.flash-data", auth(a.flashDataHandler), mux)
	mux.Handle("POST /features/state/flash-data/{action}", auth(a.flashDataActionHandler))
	routes.Handle("features.state.shared-props", auth(a.sharedPropsHandler), mux)

	// Layouts
	routes.Handle("features.layouts.persistent-layouts", auth(a.persistentLayoutsHandler), mux)
	routes.Handle("features.layouts.persistent-layouts-page-2", auth(a.persistentLayoutsPage2Handler), mux)
	routes.Handle("features.layouts.nested-layouts", auth(a.nestedLayoutsHandler), mux)
	routes.Handle("features.layouts.head", auth(a.headHandler), mux)
	routes.Handle("features.layouts.layout-props", auth(a.layoutPropsHandler), mux)

	// Events
	routes.Handle("features.events.global-events", auth(a.globalEventsHandler), mux)
	routes.Handle("features.events.visit-callbacks", auth(a.visitCallbacksHandler), mux)
	routes.Handle("features.events.progress", auth(a.progressHandler), mux)
	routes.Handle("features.events.progress-slow", auth(a.progressSlowHandler), mux)

	// HTTP
	routes.Handle("features.http.use-http", auth(a.useHttpHandler), mux)
	mux.Handle("/features/http/use-http/api", auth(a.useHttpApiHandler))
}
