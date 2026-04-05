package config

import (
	"fmt"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/spf13/viper"
)

// DefaultHead returns a Head with sensible defaults: lang "en", robots
// "index, follow", and placeholder slots for common meta and link tags.
// Empty Content values are skipped during rendering.
func DefaultHead() httpx.Head {
	head := httpx.Head{
		Lang: "en",
		Meta: []httpx.MetaTag{
			{Name: "description", Content: ""},
			{Name: "keywords", Content: ""},
			{Name: "robots", Content: "index, follow"},
			{Property: "og:title", Content: ""},
			{Property: "og:description", Content: ""},
			{Property: "og:image", Content: ""},
			{Property: "og:url", Content: ""},
			{Property: "og:type", Content: "website"},
			{Property: "og:site_name", Content: ""},
			{Property: "og:locale", Content: "en_US"},
			{Name: "twitter:card", Content: ""},
			{Name: "twitter:title", Content: ""},
			{Name: "twitter:description", Content: ""},
			{Name: "twitter:image", Content: ""},
			{Name: "twitter:site", Content: ""},
		},
		Links: []httpx.LinkTag{
			{Rel: "canonical", Href: ""},
			{Rel: "alternate", Href: "", HrefLang: ""},
		},
	}

	head.ApplyEnv()

	return head
}

// LoadHead reads a YAML head/SEO config file. Defaults are applied first,
// then the file values are merged on top, and finally environment variable
// overrides (INERTIA_SEO_*) are applied.
func LoadHead(path string) (httpx.Head, error) {
	v := viper.New()

	v.SetDefault("lang", "en")

	v.SetEnvPrefix("INERTIA_SEO")

	v.AutomaticEnv()

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return httpx.Head{}, fmt.Errorf("head: read config: %w", err)
	}

	var override httpx.Head

	if err := v.Unmarshal(&override); err != nil {
		return httpx.Head{}, fmt.Errorf("head: parse config: %w", err)
	}

	head := httpx.MergeHead(DefaultHead(), override)

	head.ApplyEnv()

	return head, nil
}
