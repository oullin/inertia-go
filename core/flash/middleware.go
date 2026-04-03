package flash

import (
	"net/http"

	"github.com/oullin/inertia-go/core/inertia"
)

type middlewareConfig struct {
	propKey string
}

// MiddlewareOption configures the flash middleware.
type MiddlewareOption func(*middlewareConfig)

// WithPropKey changes the Inertia prop key used for flash messages.
// Default: "flash".
func WithPropKey(key string) MiddlewareOption {
	return func(c *middlewareConfig) {
		c.propKey = key
	}
}

// Middleware returns an HTTP middleware that automatically consumes
// flash messages from the store and places them into the Inertia
// request context as a shared prop.
func Middleware(store Store, opts ...MiddlewareOption) func(http.Handler) http.Handler {
	cfg := &middlewareConfig{propKey: "flash"}

	for _, opt := range opts {
		opt(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if msg := store.Consume(w, r); msg != nil {
				ctx := inertia.SetProp(r.Context(), cfg.propKey, msg)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
