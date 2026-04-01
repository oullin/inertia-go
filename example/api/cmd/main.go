package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/props"
)

//go:embed resources/views/app.html
var rootTemplateFS embed.FS

var distFS = os.DirFS("../../app/dist")

var i *inertia.Inertia

func main() {
	tmpl, err := rootTemplateFS.ReadFile("resources/views/app.html")

	if err != nil {
		log.Fatal("failed to read template:", err)
	}

	version := "dev"

	if v := os.Getenv("APP_VERSION"); v != "" {
		version = v
	}

	i, err = inertia.New(string(tmpl), inertia.WithVersion(version))

	if err != nil {
		log.Fatal(err)
	}

	i.ShareProp("app", map[string]string{
		"name": "Inertia Go + Vue",
	})

	mux := http.NewServeMux()

	// Serve Vite-built assets.
	mux.Handle("/assets/", http.FileServer(http.FS(distFS)))

	// Pages.
	mux.Handle("/", i.Middleware(http.HandlerFunc(homeHandler)))
	mux.Handle("/about", i.Middleware(http.HandlerFunc(aboutHandler)))
	mux.Handle("/users", i.Middleware(http.HandlerFunc(usersHandler)))
	mux.Handle("/redirect-test", i.Middleware(http.HandlerFunc(redirectHandler)))

	addr := ":8080"

	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	if url := os.Getenv("PORTLESS_URL"); url != "" {
		fmt.Printf("Server running at %s\n", url)
	} else {
		fmt.Printf("Server running at http://localhost%s\n", addr)
	}

	log.Fatal(http.ListenAndServe(addr, mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := i.Render(w, r, "Home", httpx.Props{
		"title":   "Welcome",
		"message": "This is the Inertia.js Go adapter with Vue and Vite.",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := i.Render(w, r, "About", httpx.Props{
		"title":       "About",
		"description": "A Go server-side adapter for the Inertia.js protocol.",
		"features": []string{
			"Zero external dependencies",
			"net/http compatible middleware",
			"Partial reloads",
			"Deferred props",
			"Asset versioning",
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	err := i.Render(w, r, "Users", httpx.Props{
		"title": "Users",
		"users": []map[string]any{
			{"id": 1, "name": "Alice Johnson", "email": "alice@example.com", "role": "Admin"},
			{"id": 2, "name": "Bob Smith", "email": "bob@example.com", "role": "Editor"},
			{"id": 3, "name": "Charlie Brown", "email": "charlie@example.com", "role": "Viewer"},
		},
		"stats": props.Defer(func() any {
			return map[string]int{
				"total":  3,
				"active": 2,
			}
		}, "sidebar"),
		"metadata": props.Once(map[string]string{
			"last_sync": "2026-04-01T12:00:00Z",
		}),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	i.Redirect(w, r, "/")
}
