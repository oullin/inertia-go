package flash

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Message carries a flash notification across requests.
type Message struct {
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Store defines the interface for flash message persistence.
// CookieStore is the default implementation.
type Store interface {
	Set(w http.ResponseWriter, msg Message)
	Consume(w http.ResponseWriter, r *http.Request) *Message
}

// CookieStore persists flash messages in an HTTP cookie.
type CookieStore struct {
	cookieName string
	path       string
	httpOnly   bool
	secure     bool
	sameSite   http.SameSite
}

// Option configures a CookieStore.
type Option func(*CookieStore)

// NewCookieStore creates a CookieStore with sensible defaults.
func NewCookieStore(opts ...Option) *CookieStore {
	s := &CookieStore{
		cookieName: "inertia_flash",
		path:       "/",
		httpOnly:   true,
		sameSite:   http.SameSiteLaxMode,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithCookieName sets the cookie name.
func WithCookieName(name string) Option {
	return func(s *CookieStore) {
		s.cookieName = name
	}
}

// WithPath sets the cookie path.
func WithPath(path string) Option {
	return func(s *CookieStore) {
		s.path = path
	}
}

// WithSecure enables the Secure flag on the cookie.
func WithSecure(secure bool) Option {
	return func(s *CookieStore) {
		s.secure = secure
	}
}

// WithHTTPOnly sets the HttpOnly flag.
func WithHTTPOnly(httpOnly bool) Option {
	return func(s *CookieStore) {
		s.httpOnly = httpOnly
	}
}

// WithSameSite sets the SameSite mode.
func WithSameSite(mode http.SameSite) Option {
	return func(s *CookieStore) {
		s.sameSite = mode
	}
}

// Set serializes msg as JSON, URL-escapes it, and writes it as a cookie.
func (s *CookieStore) Set(w http.ResponseWriter, msg Message) {
	data, err := json.Marshal(msg)

	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    url.QueryEscape(string(data)),
		Path:     s.path,
		HttpOnly: s.httpOnly,
		Secure:   s.secure,
		SameSite: s.sameSite,
	})
}

// Consume reads the flash cookie, deletes it, and returns the decoded
// Message. Returns nil if no flash cookie exists or decoding fails.
func (s *CookieStore) Consume(w http.ResponseWriter, r *http.Request) *Message {
	cookie, err := r.Cookie(s.cookieName)

	if err != nil {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:   s.cookieName,
		Value:  "",
		Path:   s.path,
		MaxAge: -1,
	})

	value, err := url.QueryUnescape(cookie.Value)

	if err != nil {
		return nil
	}

	var payload Message

	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return nil
	}

	return &payload
}
