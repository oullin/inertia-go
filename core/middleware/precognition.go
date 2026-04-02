package middleware

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
)

// Precognition returns an HTTP middleware that detects precognition
// requests (Precognition: true header) and marks them in the request
// context. The Inertia Render method uses this flag to return
// validation-only responses (204 or 422) instead of full pages.
func Precognition() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appendVary(w.Header(), httpx.HeaderPrecognition)

			if !httpx.IsPrecognitionRequest(r) {
				next.ServeHTTP(w, r)

				return
			}

			ctx := httpx.SetPrecognition(r.Context())
			r = r.WithContext(ctx)
			w.Header().Set(httpx.HeaderPrecognition, "true")

			next.ServeHTTP(w, r)
		})
	}
}

func appendVary(h http.Header, value string) {
	if value == "" {
		return
	}

	existing := h.Get("Vary")

	if existing == "" {
		h.Set("Vary", value)

		return
	}

	for _, part := range strings.Split(existing, ",") {
		if strings.TrimSpace(part) == value {
			return
		}
	}

	h.Set("Vary", existing+", "+value)
}
