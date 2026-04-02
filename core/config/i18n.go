package config

import (
	"fmt"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/spf13/viper"
)

// I18nConfig holds the multilanguage configuration.
type I18nConfig struct {
	DefaultLocale string                   `mapstructure:"default_locale"`
	URLPrefix     bool                     `mapstructure:"url_prefix"`
	Locales       map[string]*httpx.Locale `mapstructure:"locales"`
}

// DefaultI18n returns an I18nConfig with a single English locale.
func DefaultI18n() *I18nConfig {
	cfg := &I18nConfig{
		DefaultLocale: "en",
		URLPrefix:     false,
		Locales: map[string]*httpx.Locale{
			"en": {
				Code:      "en",
				Name:      "English",
				Direction: "ltr",
			},
		},
	}

	return cfg
}

// LoadI18n reads a YAML i18n config file. Defaults are applied first,
// then the file values are merged on top, and finally environment
// variable overrides (INERTIA_I18N_*) are applied.
func LoadI18n(path string) (*I18nConfig, error) {
	v := viper.New()
	v.SetDefault("default_locale", "en")
	v.SetDefault("url_prefix", false)

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("i18n: read config: %w", err)
	}

	v.SetEnvPrefix("INERTIA_I18N")
	v.AutomaticEnv()

	var cfg I18nConfig

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("i18n: parse config: %w", err)
	}

	// Backfill locale codes from map keys.
	for code, locale := range cfg.Locales {
		locale.Code = code
	}

	return &cfg, nil
}

// Lookup returns the Locale for the given code, or nil if not found.
func (cfg *I18nConfig) Lookup(code string) *httpx.Locale {
	return cfg.Locales[code]
}

// Default returns the default Locale.
func (cfg *I18nConfig) Default() *httpx.Locale {
	return cfg.Locales[cfg.DefaultLocale]
}

// Codes returns all configured locale codes.
func (cfg *I18nConfig) Codes() []string {
	codes := make([]string, 0, len(cfg.Locales))

	for code := range cfg.Locales {
		codes = append(codes, code)
	}

	return codes
}
