package i18n

import (
	"fmt"
	"os"

	"github.com/oullin/inertia-go/core/httpx"
	"gopkg.in/yaml.v3"
)

// Config holds the multilanguage configuration loaded from a YAML file.
type Config struct {
	DefaultLocale string                   `yaml:"default_locale"`
	URLPrefix     bool                     `yaml:"url_prefix"`
	Locales       map[string]*httpx.Locale `yaml:"locales"`
}

// LoadConfig reads a YAML i18n config file and applies env var overrides.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("i18n: read config: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("i18n: parse config: %w", err)
	}

	// Backfill locale codes from map keys.
	for code, locale := range cfg.Locales {
		locale.Code = code
	}

	cfg.applyEnv()

	return &cfg, nil
}

// Lookup returns the Locale for the given code, or nil if not found.
func (cfg *Config) Lookup(code string) *httpx.Locale {
	return cfg.Locales[code]
}

// Default returns the default Locale.
func (cfg *Config) Default() *httpx.Locale {
	return cfg.Locales[cfg.DefaultLocale]
}

// Codes returns all configured locale codes.
func (cfg *Config) Codes() []string {
	codes := make([]string, 0, len(cfg.Locales))

	for code := range cfg.Locales {
		codes = append(codes, code)
	}

	return codes
}

func (cfg *Config) applyEnv() {
	if v := os.Getenv("INERTIA_I18N_DEFAULT_LOCALE"); v != "" {
		cfg.DefaultLocale = v
	}
}
