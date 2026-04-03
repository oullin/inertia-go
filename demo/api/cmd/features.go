package main

import "net/http"

func registerFeatureRoutes(mux *http.ServeMux) {
	mux.Handle("/features/forms/use-form", requireDemoAuth(http.HandlerFunc(formsHandler)))
	mux.Handle("/features/navigation/links", requireDemoAuth(http.HandlerFunc(navigationHandler)))
	mux.Handle("/features/data-loading/deferred-props", requireDemoAuth(http.HandlerFunc(dataHandler)))
	mux.Handle("/features/state/remember", requireDemoAuth(http.HandlerFunc(stateHandler)))
}
