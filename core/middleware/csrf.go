package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/oullin/inertia-go/core/config"
	"github.com/oullin/inertia-go/core/cryptox"
	"github.com/oullin/inertia-go/core/httpx"
)

// statusPageExpired is the 419 "Page Expired" status code used when
// CSRF validation fails.
const statusPageExpired = 419

// CSRF returns an HTTP middleware that provides CSRF protection. A random
// token is encrypted into an XSRF-TOKEN
// cookie (not HTTP-only, so JavaScript can read it). The same raw token
// is placed in the request context for the Inertia Render method to emit
// as a <meta> tag. Mutation requests must present the token via one of
// three sources: _token form field, X-CSRF-TOKEN header, or X-XSRF-TOKEN
// header (the encrypted cookie value, auto-sent by Axios).
func CSRF(cfg config.CSRFConfig, key []byte) func(http.Handler) http.Handler {
	cfg.Defaults()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read existing token from cookie, or generate a new one.
			token, err := tokenFromCookie(r, cfg.CookieName, key)

			if err != nil {
				token, err = generateToken()

				if err != nil {
					http.Error(w, "csrf: failed to generate token", http.StatusInternalServerError)

					return
				}

				setTokenCookie(w, cfg.CookieName, token, key, cfg.Secure, cfg.SameSiteMode())
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

			// Mutation requests must present the token.
			submittedToken := extractToken(r, key)

			if subtle.ConstantTimeCompare([]byte(token), []byte(submittedToken)) != 1 {
				http.Error(w, "csrf: token mismatch", statusPageExpired)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRFFromFile reads YAML config files and returns the CSRF middleware.
// It loads both the CSRF config and the crypto config (for the encryption key).
func CSRFFromFile(csrfPath, cryptoPath string) (func(http.Handler) http.Handler, error) {
	cfg, err := config.LoadCSRF(csrfPath)

	if err != nil {
		return nil, err
	}

	cryptoCfg, err := config.LoadCrypto(cryptoPath)

	if err != nil {
		return nil, err
	}

	key, err := cryptoCfg.DecodedKey()

	if err != nil {
		return nil, err
	}

	return CSRF(cfg, key), nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func setTokenCookie(w http.ResponseWriter, name, token string, key []byte, secure bool, sameSite http.SameSite) {
	encrypted, err := cryptox.Encrypt(token, key)

	if err != nil {
		http.Error(w, "csrf: failed to encrypt token", http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encrypted,
		Path:     "/",
		HttpOnly: false,
		Secure:   secure,
		SameSite: sameSite,
	})
}

func tokenFromCookie(r *http.Request, name string, key []byte) (string, error) {
	cookie, err := r.Cookie(name)

	if err != nil {
		return "", err
	}

	return cryptox.Decrypt(cookie.Value, key)
}

// extractToken retrieves the submitted CSRF token from the request,
// checking sources in order: _token form field, X-CSRF-TOKEN
// header, then X-XSRF-TOKEN header. The X-XSRF-TOKEN value is the
// encrypted cookie value (auto-sent by Axios) and must be decrypted.
func extractToken(r *http.Request, key []byte) string {
	// 1. _token POST form field (raw token from HTML forms).
	if token := r.PostFormValue("_token"); token != "" {
		return token
	}

	// 2. X-CSRF-TOKEN header (raw token, typically from meta tag).
	if token := r.Header.Get("X-CSRF-TOKEN"); token != "" {
		return token
	}

	// 3. X-XSRF-TOKEN header (encrypted value from cookie, auto-sent by Axios).
	if encrypted := r.Header.Get("X-XSRF-TOKEN"); encrypted != "" {
		token, err := cryptox.Decrypt(encrypted, key)

		if err != nil {
			return ""
		}

		return token
	}

	return ""
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
