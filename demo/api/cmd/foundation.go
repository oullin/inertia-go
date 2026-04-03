package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/demo/api/internal/database"
)

type ctxKey string

const demoSessionCookieName = "inertia_go_demo_session"

const currentUserKey ctxKey = "current_user"

var demoRouteManifest = map[string]string{
	"home":                           "/",
	"login":                          "/login",
	"logout":                         "/logout",
	"dashboard":                      "/dashboard",
	"contacts.index":                 "/contacts",
	"contacts.create":                "/contacts/create",
	"contacts.store":                 "/contacts",
	"contacts.show":                  "/contacts/{contact}",
	"contacts.edit":                  "/contacts/{contact}/edit",
	"contacts.update":                "/contacts/{contact}",
	"contacts.favorite":              "/contacts/{contact}/favorite",
	"contacts.notes.store":           "/contacts/{contact}/notes",
	"organizations.index":            "/organizations",
	"organizations.show":             "/organizations/{organization}",
	"organizations.update":           "/organizations/{organization}",
	"features.forms.use-form":        "/features/forms/use-form",
	"features.navigation.links":      "/features/navigation/links",
	"features.data-loading.deferred": "/features/data-loading/deferred-props",
	"features.state.remember":        "/features/state/remember",
}

func withDemoProps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := loadCurrentUser(r)
		ctx := context.WithValue(r.Context(), currentUserKey, user)
		ctx = inertia.SetProps(ctx, httpx.Props{
			"app": map[string]any{
				"name":        "Inertia.js Kitchen Sink",
				"productLine": "Go Demo Port",
				"environment": "Demo",
			},
			"auth": map[string]any{
				"user": publicUser(user),
			},
			"workspace": map[string]any{
				"name": "Inertia Go",
				"plan": "Porting",
			},
			"routes": manifestProps(),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requireDemoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if currentUser(r) == nil {
			i.Redirect(w, r, routeURL("login", nil))

			return
		}

		next.ServeHTTP(w, r)
	})
}

func guestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if currentUser(r) != nil {
			i.Redirect(w, r, routeURL("dashboard", nil))

			return
		}

		next.ServeHTTP(w, r)
	})
}

func currentUser(r *http.Request) *database.User {
	if user, ok := r.Context().Value(currentUserKey).(*database.User); ok {
		return user
	}

	return nil
}

func loadCurrentUser(r *http.Request) *database.User {
	cookie, err := r.Cookie(demoSessionCookieName)

	if err != nil || cookie.Value == "" || db == nil {
		return nil
	}

	id, err := strconv.ParseInt(cookie.Value, 10, 64)

	if err != nil {
		return nil
	}

	user, err := database.FindUserByID(db, id)

	if err != nil {
		return nil
	}

	return user
}

func setDemoSession(w http.ResponseWriter, userID int64, remember bool) {
	cookie := &http.Cookie{
		Name:     demoSessionCookieName,
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

func clearDemoSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     demoSessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func renderPage(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	renderPageWithContext(w, r, component, pageProps)
}

func renderPageWithContext(w http.ResponseWriter, r *http.Request, component string, pageProps httpx.Props) {
	ctx := r.Context()

	if flash := consumeFlash(w, r); flash != nil {
		ctx = inertia.SetProp(ctx, "flash", flash)
	}

	if err := i.Render(w, r.WithContext(ctx), component, pageProps); err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			http.Error(w, "page not found", http.StatusNotFound)
		default:
			http.Error(w, "demo internal error", http.StatusInternalServerError)
		}
	}
}

func publicUser(user *database.User) any {
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

func manifestProps() map[string]any {
	props := make(map[string]any, len(demoRouteManifest))

	for name, pattern := range demoRouteManifest {
		props[name] = pattern
	}

	return props
}

func routeURL(name string, params map[string]string) string {
	pattern, ok := demoRouteManifest[name]

	if !ok {
		return "/"
	}

	for key, value := range params {
		pattern = strings.ReplaceAll(pattern, "{"+key+"}", value)
	}

	return pattern
}
