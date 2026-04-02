package i18n

import (
	"net/http"
	"sort"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
)

// Middleware returns an HTTP middleware that detects the locale from the
// URL prefix (e.g. /es/dashboard), strips the prefix, sets the locale
// in context, and auto-generates hreflang alternate links.
func (cfg *Config) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale, cleanPath := cfg.resolve(r.URL.Path)

		// Rewrite the URL path with the prefix stripped so downstream
		// handlers see clean paths (e.g. /dashboard, not /es/dashboard).
		r.URL.Path = cleanPath
		r.RequestURI = cleanPath

		if q := r.URL.RawQuery; q != "" {
			r.RequestURI = cleanPath + "?" + q
		}

		// Build locale head with auto-generated hreflang links.
		localeHead := locale.Head
		localeHead.Links = append(localeHead.Links, cfg.hreflangLinks(cleanPath)...)

		ctx := r.Context()
		ctx = httpx.SetLocale(ctx, &httpx.Locale{
			Code:      locale.Code,
			Name:      locale.Name,
			Direction: locale.Direction,
			Head:      localeHead,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// resolve extracts the locale code from a URL-prefix path. Returns the
// matching Locale and the path with the prefix stripped. Falls back to
// the default locale when no prefix matches.
func (cfg *Config) resolve(path string) (*Locale, string) {
	if !cfg.URLPrefix {
		return cfg.Default(), path
	}

	// Path format: /es/dashboard or /es or /es/
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)

	if len(parts) > 0 {
		code := parts[0]

		if locale := cfg.Lookup(code); locale != nil {
			clean := "/"

			if len(parts) > 1 {
				clean = "/" + parts[1]
			}

			return locale, clean
		}
	}

	return cfg.Default(), path
}

// hreflangLinks builds <link rel="alternate" hreflang> tags for all
// configured locales, using the given clean path.
func (cfg *Config) hreflangLinks(cleanPath string) []httpx.LinkTag {
	if len(cfg.Locales) <= 1 {
		return nil
	}

	// Sort codes for deterministic output.
	codes := cfg.Codes()

	sort.Strings(codes)

	links := make([]httpx.LinkTag, 0, len(codes))

	for _, code := range codes {
		href := cleanPath

		if cfg.URLPrefix {
			href = "/" + code + cleanPath

			if href != "/" && strings.HasSuffix(href, "/") {
				href = strings.TrimSuffix(href, "/")
			}
		}

		links = append(links, httpx.LinkTag{
			Rel:      "alternate",
			Href:     href,
			HrefLang: code,
		})
	}

	return links
}
