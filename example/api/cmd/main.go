package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oullin/inertia-go/core/httpx"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/example/api/internal/database"
)

//go:embed resources/views/app.html
var rootTemplateFS embed.FS

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

	db, err = database.Open("beacon.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	distPath, err := resolveDistPath()

	if err != nil {
		log.Fatal(err)
	}

	i.ShareProps(httpx.Props{
		"app": map[string]any{
			"name":        "Progressive Oullin",
			"productLine": "Documents",
			"environment": "Production",
		},
		"auth": map[string]any{
			"user": map[string]any{
				"name":     "shadcn",
				"email":    "m@example.com",
				"initials": "CN",
			},
		},
		"workspace": map[string]any{
			"name": "Oullin.io",
			"plan": "Growth",
		},
	})

	mux := http.NewServeMux()
	mux.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir(distPath))),
	)

	registerDashboardRoutes(mux)

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

func resolveDistPath() (string, error) {
	candidates := []string{
		"example/app/dist",
		"../app/dist",
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return filepath.Clean(candidate), nil
		}
	}

	return "", fmt.Errorf("failed to locate example app dist directory")
}
