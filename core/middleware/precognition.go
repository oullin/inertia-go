package middleware

import (
	"net/http"

	"github.com/oullin/inertia-go/core/httpx"
)

// Precognition returns an HTTP middleware that detects precognition
// requests (Precognition: true header) and marks them in the request
// context. The Inertia Render method uses this flag to return
// validation-only responses (204 or 422) instead of full pages.
func Precognition() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", httpx.HeaderPrecognition)

			if !httpx.IsPrecognitionRequest(r) {
				next.ServeHTTP(w, r)

				return
			}

			ctx := httpx.SetPrecognition(r.Context())
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
