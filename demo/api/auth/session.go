package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/oullin/inertia-go/demo/api/internal/database"
)

type ctxKey string

// SessionCookieName is the cookie used by the demo auth flow.
const SessionCookieName = "inertia_go_demo_session"

const currentUserKey ctxKey = "current_user"

// WithCurrentUser resolves the demo user from the session cookie into request context.
func (a App) WithCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), currentUserKey, a.loadCurrentUser(r))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth redirects unauthenticated requests to the login page.
func (a App) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.CurrentUser(r) == nil {
			a.deps.Redirect(w, r, a.deps.RouteURL("login", nil))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// GuestOnly redirects authenticated requests away from guest-only pages.
func (a App) GuestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.CurrentUser(r) != nil {
			a.deps.Redirect(w, r, a.deps.RouteURL("dashboard", nil))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// CurrentUser returns the authenticated demo user from the request context.
func (a App) CurrentUser(r *http.Request) *database.User {
	if user, ok := r.Context().Value(currentUserKey).(*database.User); ok {
		return user
	}

	return nil
}

func (a App) loadCurrentUser(r *http.Request) *database.User {
	cookie, err := r.Cookie(SessionCookieName)

	if err != nil || cookie.Value == "" || a.deps.DB == nil {
		return nil
	}

	id, err := strconv.ParseInt(cookie.Value, 10, 64)

	if err != nil {
		return nil
	}

	user, err := database.FindUserByID(a.deps.DB, id)

	if err != nil {
		return nil
	}

	return user
}

func (a App) setSession(w http.ResponseWriter, userID int64, remember bool) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    strconv.FormatInt(userID, 10),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if remember {
		cookie.MaxAge = 60 * 60 * 24 * 30
	}

	http.SetCookie(w, cookie)
}

func (a App) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// PublicUser returns the frontend-safe user payload shared in page props.
func (a App) PublicUser(user *database.User) any {
	if user == nil {
		return nil
	}

	initials := ""

	for _, part := range strings.Fields(user.Name) {
		if part != "" {
			initials += strings.ToUpper(part[:1])
		}
	}

	return map[string]any{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"initials": initials,
	}
}
