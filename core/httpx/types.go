package httpx

import (
	"context"
	"html"
	"net/http"
	"os"
	"strings"
)

// Props holds page properties passed to the frontend component.
type Props map[string]any

// TemplateData holds additional data passed to the root HTML template
// during initial (non-XHR) page visits.
type TemplateData map[string]any

// ValidationErrors holds field-level validation errors to be shared
// with the frontend via the "errors" prop.
type ValidationErrors map[string]any

// JSONMarshaler abstracts JSON encoding/decoding so callers can swap
// in a faster implementation (e.g. github.com/goccy/go-json) without
// changing any Inertia code.
type JSONMarshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// Logger is a minimal logging interface compatible with the standard
// library's log.Logger and most structured logging packages.
type Logger interface {
	Printf(format string, v ...any)
}

// MetaTag represents a single <meta> element for server-side head rendering.
// Exactly one of Name or Property should be set to identify the tag.
type MetaTag struct {
	Name     string `json:"name,omitempty"     yaml:"name,omitempty"`
	Property string `json:"property,omitempty" yaml:"property,omitempty"`
	Content  string `json:"content"            yaml:"content"`
}

// LinkTag represents a single <link> element for server-side head rendering.
type LinkTag struct {
	Rel      string `json:"rel"                yaml:"rel"`
	Href     string `json:"href"               yaml:"href"`
	HrefLang string `json:"hreflang,omitempty"  yaml:"hreflang,omitempty"`
	Type     string `json:"type,omitempty"      yaml:"type,omitempty"`
}

// Head holds the server-side head elements rendered into {{ .inertiaHead }}
// on initial page loads, and the lang/dir rendered into {{ .inertiaLang }}
// and {{ .inertiaDir }}.
type Head struct {
	Title     string    `json:"title,omitempty"     yaml:"title,omitempty"`
	Lang      string    `json:"lang,omitempty"      yaml:"lang,omitempty"`
	Direction string    `json:"direction,omitempty"  yaml:"direction,omitempty"`
	Meta      []MetaTag `json:"meta,omitempty"      yaml:"meta,omitempty"`
	Links     []LinkTag `json:"links,omitempty"     yaml:"links,omitempty"`
}

// metaKey returns the identifying key for a MetaTag. Tags are considered
// the same if they share the same Name or Property value.

// linkKey returns the identifying key for a LinkTag. Alternate links are
// keyed by rel+hreflang; other links are keyed by rel alone.

// RenderHTML produces the HTML string for embedding in {{ .inertiaHead }}.
// Meta tags with empty Content are skipped (they serve as YAML placeholders).
// Returns "" for a zero-value Head, preserving backward compatibility.

// MergeHead combines a base Head with an override Head. The override's
// Title, Lang, and Direction replace the base values if non-empty. Meta
// tags are merged by their identifying attribute (Name or Property);
// override tags replace base tags with the same key, and new tags are
// appended. Links are merged by Rel+HrefLang key.

// Merge meta tags.

// Merge link tags.

// envBindings maps env var suffixes to meta tag identifiers.

// MetaTag.Name match
// MetaTag.Property match

// ApplyEnv overrides Head fields with values from environment variables
// when present. Uses the INERTIA_SEO_ prefix convention.

// Locale holds locale information resolved by the i18n middleware.
type Locale struct {
	Code      string
	Name      string
	Direction string
	Head      Head
}

type ctxKey struct{ name string }

func metaKey(tag MetaTag) string {
	if tag.Name != "" {
		return "name:" + tag.Name
	}

	return "property:" + tag.Property
}

func linkKey(tag LinkTag) string {
	if tag.HrefLang != "" {
		return tag.Rel + ":" + tag.HrefLang
	}

	return tag.Rel
}

func (h Head) RenderHTML() string {
	if h.Title == "" && len(h.Meta) == 0 && len(h.Links) == 0 {
		return ""
	}

	var b strings.Builder

	if h.Title != "" {
		b.WriteString("<title>")
		b.WriteString(html.EscapeString(h.Title))
		b.WriteString("</title>\n")
	}

	for _, tag := range h.Meta {
		if tag.Content == "" {
			continue
		}

		b.WriteString("<meta ")

		switch {
		case tag.Name != "":
			b.WriteString(`name="`)
			b.WriteString(html.EscapeString(tag.Name))
			b.WriteString(`"`)
		case tag.Property != "":
			b.WriteString(`property="`)
			b.WriteString(html.EscapeString(tag.Property))
			b.WriteString(`"`)
		}

		b.WriteString(` content="`)
		b.WriteString(html.EscapeString(tag.Content))
		b.WriteString("\">\n")
	}

	for _, link := range h.Links {
		if link.Href == "" {
			continue
		}

		b.WriteString(`<link rel="`)
		b.WriteString(html.EscapeString(link.Rel))
		b.WriteString(`" href="`)
		b.WriteString(html.EscapeString(link.Href))
		b.WriteString(`"`)

		if link.HrefLang != "" {
			b.WriteString(` hreflang="`)
			b.WriteString(html.EscapeString(link.HrefLang))
			b.WriteString(`"`)
		}

		if link.Type != "" {
			b.WriteString(` type="`)
			b.WriteString(html.EscapeString(link.Type))
			b.WriteString(`"`)
		}

		b.WriteString(">\n")
	}

	return b.String()
}

func MergeHead(base, override Head) Head {
	result := Head{
		Title:     base.Title,
		Lang:      base.Lang,
		Direction: base.Direction,
	}

	if override.Title != "" {
		result.Title = override.Title
	}

	if override.Lang != "" {
		result.Lang = override.Lang
	}

	if override.Direction != "" {
		result.Direction = override.Direction
	}

	metaSeen := make(map[string]int, len(base.Meta))
	merged := make([]MetaTag, len(base.Meta))
	copy(merged, base.Meta)

	for i, tag := range merged {
		metaSeen[metaKey(tag)] = i
	}

	for _, tag := range override.Meta {
		key := metaKey(tag)

		if idx, ok := metaSeen[key]; ok {
			merged[idx] = tag
		} else {
			metaSeen[key] = len(merged)
			merged = append(merged, tag)
		}
	}

	result.Meta = merged

	linkSeen := make(map[string]int, len(base.Links))
	mergedLinks := make([]LinkTag, len(base.Links))
	copy(mergedLinks, base.Links)

	for i, link := range mergedLinks {
		linkSeen[linkKey(link)] = i
	}

	for _, link := range override.Links {
		key := linkKey(link)

		if idx, ok := linkSeen[key]; ok {
			mergedLinks[idx] = link
		} else {
			linkSeen[key] = len(mergedLinks)
			mergedLinks = append(mergedLinks, link)
		}
	}

	result.Links = mergedLinks

	return result
}

var envBindings = []struct {
	envSuffix string
	name      string
	property  string
}{
	{"DESCRIPTION", "description", ""},
	{"KEYWORDS", "keywords", ""},
	{"ROBOTS", "robots", ""},
	{"GOOGLEBOT", "googlebot", ""},
	{"OG_TITLE", "", "og:title"},
	{"OG_DESCRIPTION", "", "og:description"},
	{"OG_IMAGE", "", "og:image"},
	{"OG_URL", "", "og:url"},
	{"OG_TYPE", "", "og:type"},
	{"OG_SITE_NAME", "", "og:site_name"},
	{"OG_LOCALE", "", "og:locale"},
	{"TWITTER_CARD", "twitter:card", ""},
	{"TWITTER_TITLE", "twitter:title", ""},
	{"TWITTER_DESCRIPTION", "twitter:description", ""},
	{"TWITTER_IMAGE", "twitter:image", ""},
	{"TWITTER_SITE", "twitter:site", ""},
}

func (h *Head) ApplyEnv() {
	if v := os.Getenv("INERTIA_SEO_TITLE"); v != "" {
		h.Title = v
	}

	if v := os.Getenv("INERTIA_SEO_LANG"); v != "" {
		h.Lang = v
	}

	for _, binding := range envBindings {
		v := os.Getenv("INERTIA_SEO_" + binding.envSuffix)

		if v == "" {
			continue
		}

		found := false

		for i := range h.Meta {
			if binding.name != "" && h.Meta[i].Name == binding.name {
				h.Meta[i].Content = v
				found = true

				break
			}

			if binding.property != "" && h.Meta[i].Property == binding.property {
				h.Meta[i].Content = v
				found = true

				break
			}
		}

		if !found {
			tag := MetaTag{Content: v}

			if binding.name != "" {
				tag.Name = binding.name
			} else {
				tag.Property = binding.property
			}

			h.Meta = append(h.Meta, tag)
		}
	}
}

var (
	ctxKeyCSRFToken = &ctxKey{"csrfToken"}
	ctxKeyLocale    = &ctxKey{"locale"}
)

// SetCSRFToken stores a CSRF token in the request context. When present,
// Render automatically adds <meta name="csrf-token" content="TOKEN"> to
// the head on initial page loads.
func SetCSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxKeyCSRFToken, token)
}

// CSRFTokenFromContext returns the CSRF token stored in context, or "".
func CSRFTokenFromContext(ctx context.Context) string {
	s, _ := ctx.Value(ctxKeyCSRFToken).(string)

	return s
}

// SetLocale stores the resolved locale in the request context.
func SetLocale(ctx context.Context, locale *Locale) context.Context {
	return context.WithValue(ctx, ctxKeyLocale, locale)
}

// LocaleFromContext returns the locale stored in context, or nil.
func LocaleFromContext(ctx context.Context) *Locale {
	l, _ := ctx.Value(ctxKeyLocale).(*Locale)

	return l
}

// IsInertiaRequest reports whether r was initiated by the Inertia.js
// client (i.e. it carries the X-Inertia header).
func IsInertiaRequest(r *http.Request) bool {
	return strings.TrimSpace(r.Header.Get(HeaderInertia)) == "true"
}
