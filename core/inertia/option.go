package inertia

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"os"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/spf13/viper"
)

// Option configures an Inertia instance during construction.
type Option func(*Inertia) error

// WithVersion sets a static asset version string.
func WithVersion(version string) Option {
	return func(i *Inertia) error {
		i.version = version

		return nil
	}
}

// WithVersionFromFile hashes the contents of path to produce the
// asset version string. This is useful for cache-busting when a
// manifest file (e.g. build/manifest.json) changes on deploy.
func WithVersionFromFile(path string) Option {
	return func(i *Inertia) error {
		data, err := os.ReadFile(path)

		if err != nil {
			return fmt.Errorf("inertia: version file: %w", err)
		}

		i.version = fmt.Sprintf("%x", md5.Sum(data))

		return nil
	}
}

// WithContainerID sets the HTML element ID used for the root container
// div. Defaults to "app".
func WithContainerID(id string) Option {
	return func(i *Inertia) error {
		i.containerID = id

		return nil
	}
}

// WithJSONMarshaler replaces the default encoding/json marshaler with
// a custom implementation.
func WithJSONMarshaler(m httpx.JSONMarshaler) Option {
	return func(i *Inertia) error {
		i.jsonMarshaler = m

		return nil
	}
}

// WithLogger sets a logger for diagnostic messages.
func WithLogger(l httpx.Logger) Option {
	return func(i *Inertia) error {
		i.logger = l

		return nil
	}
}

// WithTemplateFuncs registers additional template functions available
// in the root HTML template.
func WithTemplateFuncs(funcMap template.FuncMap) Option {
	return func(i *Inertia) error {
		i.templateFuncs = funcMap

		return nil
	}
}

// WithEncryptHistory enables encrypted browser history by default for
// all responses.
func WithEncryptHistory() Option {
	return func(i *Inertia) error {
		i.encryptHistory = true

		return nil
	}
}

// WithHead sets default head elements rendered into {{ .inertiaHead }} on
// every initial page load. Per-request head elements (set via SetHead,
// SetTitle, or SetMeta) override these defaults.
func WithHead(head httpx.Head) Option {
	return func(i *Inertia) error {
		i.head = head

		return nil
	}
}

// WithHeadFromFile reads a YAML file at path and sets the default head
// elements. After parsing the YAML, environment variable overrides are
// applied (see Head.ApplyEnv). Meta tags with empty Content serve as
// placeholders and are excluded from rendering.
func WithHeadFromFile(path string) Option {
	return func(i *Inertia) error {
		v := viper.New()
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("inertia: head file: %w", err)
		}

		var head httpx.Head

		if err := v.Unmarshal(&head); err != nil {
			return fmt.Errorf("inertia: parse head yaml: %w", err)
		}

		head.ApplyEnv()
		i.head = head

		return nil
	}
}
