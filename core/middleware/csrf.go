package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/spf13/viper"
)

// CSRFConfig holds configuration for the CSRF middleware.
type CSRFConfig struct {
	Secret     string `json:"secret"      yaml:"secret"      mapstructure:"secret"`
	CookieName string `json:"cookie_name" yaml:"cookie_name" mapstructure:"cookie_name"`
	HeaderName string `json:"header_name" yaml:"header_name" mapstructure:"header_name"`
	Secure     bool   `json:"secure"      yaml:"secure"      mapstructure:"secure"`
	SameSite   string `json:"same_site"   yaml:"same_site"   mapstructure:"same_site"`
}

func (c *CSRFConfig) defaults() {
	if c.CookieName == "" {
		c.CookieName = "_csrf_token"
	}

	if c.HeaderName == "" {
		c.HeaderName = "X-CSRF-TOKEN"
	}

	if c.SameSite == "" {
		c.SameSite = "lax"
	}
}

func (c *CSRFConfig) sameSite() http.SameSite {
	switch strings.ToLower(c.SameSite) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

// CSRF returns an HTTP middleware that provides CSRF protection using
// the double-submit cookie pattern. A random token is generated and
// stored in an HTTP-only cookie; the same token is placed in the request
// context so the Inertia Render method can emit it as a <meta> tag.
// Mutation requests must include the token in the X-CSRF-TOKEN header.
func CSRF(cfg CSRFConfig) func(http.Handler) http.Handler {
	cfg.defaults()

	secret := []byte(cfg.Secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read existing token from cookie, or generate a new one.
			token, err := tokenFromCookie(r, cfg.CookieName, secret)

			if err != nil {
				token, err = generateToken()

				if err != nil {
					http.Error(w, "csrf: failed to generate token", http.StatusInternalServerError)

					return
				}

				setTokenCookie(w, cfg.CookieName, token, secret, cfg.Secure, cfg.sameSite())
			}

			// Store the raw token in context so Render auto-appends
			// <meta name="csrf-token" content="TOKEN">.
			ctx := httpx.SetCSRFToken(r.Context(), token)
			r = r.WithContext(ctx)

			// Safe methods pass through without validation.
			if isSafeMethod(r.Method) {
				next.ServeHTTP(w, r)

				return
			}

			// Mutation requests must present the token in the header.
			headerToken := r.Header.Get(cfg.HeaderName)

			if subtle.ConstantTimeCompare([]byte(token), []byte(headerToken)) != 1 {
				http.Error(w, "csrf: token mismatch", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRFFromFile reads a YAML config file and returns the CSRF middleware.
// After parsing, env var overrides are applied via Viper's AutomaticEnv.
func CSRFFromFile(path string) (func(http.Handler) http.Handler, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("csrf: read config: %w", err)
	}

	v.SetEnvPrefix("INERTIA_CSRF")
	v.AutomaticEnv()

	var cfg CSRFConfig

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("csrf: parse config: %w", err)
	}

	return CSRF(cfg), nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func signToken(token string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(token))

	return hex.EncodeToString(mac.Sum(nil))
}

func setTokenCookie(w http.ResponseWriter, name, token string, secret []byte, secure bool, sameSite http.SameSite) {
	signed := token + "." + signToken(token, secret)

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}

func tokenFromCookie(r *http.Request, name string, secret []byte) (string, error) {
	cookie, err := r.Cookie(name)

	if err != nil {
		return "", err
	}

	parts := strings.SplitN(cookie.Value, ".", 2)

	if len(parts) != 2 {
		return "", fmt.Errorf("csrf: invalid cookie format")
	}

	token, sig := parts[0], parts[1]
	expected := signToken(token, secret)

	if subtle.ConstantTimeCompare([]byte(sig), []byte(expected)) != 1 {
		return "", fmt.Errorf("csrf: invalid signature")
	}

	return token, nil
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
