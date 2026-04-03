package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oullin/inertia-go/core/config"
	corei18n "github.com/oullin/inertia-go/core/i18n"
	"github.com/oullin/inertia-go/core/inertia"
	"github.com/oullin/inertia-go/core/middleware"
	"github.com/oullin/inertia-go/demo/api/internal/database"
	"github.com/oullin/inertia-go/demo/api/internal/seed"
)

//go:embed resources/views/app.html
var rootTemplateFS embed.FS

var i *inertia.Inertia
var localeCfg *config.I18nConfig

func main() {
	tmpl, err := rootTemplateFS.ReadFile("resources/views/app.html")

	if err != nil {
		log.Fatal("failed to read template:", err)
	}

	version := "dev"

	if v := os.Getenv("APP_VERSION"); v != "" {
		version = v
	}

	seoPath, err := resolveResourcePath("seo.yml")

	if err != nil {
		log.Fatal(err)
	}

	i, err = inertia.New(string(tmpl),
		inertia.WithVersion(version),
		inertia.WithHeadFromFile(seoPath),
	)

	if err != nil {
		log.Fatal(err)
	}

	localeCfg, err = corei18n.LoadConfig(mustResolveResourcePath("i18n.yml"))

	if err != nil {
		log.Fatal(err)
	}

	// The demo keeps canonical, non-prefixed routes in the frontend while
	// still consuming locale-driven head defaults from config.
	localeCfg.URLPrefix = false

	csrfMiddleware, err := middleware.CSRFFromFile(
		mustResolveResourcePath("csrf.yml"),
		mustResolveResourcePath("crypto.yml"),
	)

	if err != nil {
		log.Fatal(err)
	}

	db, err = database.Open("beacon.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := seed.Run(db); err != nil {
		log.Fatal("failed to seed database:", err)
	}

	distPath, err := resolveDistPath()

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir(distPath))),
	)

	appMux := http.NewServeMux()
	registerAuthRoutes(appMux)
	registerCRMRoutes(appMux)
	registerFeatureRoutes(appMux)
	registerLegacyDashboardRoutes(appMux)
	mux.Handle("/", dashboardAppHandler(withDemoProps(appMux), csrfMiddleware, localeCfg))

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
		"storage/dist/app",
		"demo/app/dist",
		"../app/dist",
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return filepath.Clean(candidate), nil
		}
	}

	return "", fmt.Errorf("failed to locate demo app dist directory")
}

func resolveResourcePath(name string) (string, error) {
	candidates := []string{
		filepath.Join("cmd", "resources", name),
		filepath.Join("resources", name),
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return filepath.Clean(candidate), nil
		}
	}

	return "", fmt.Errorf("failed to locate demo resource %q", name)
}

func mustResolveResourcePath(name string) string {
	path, err := resolveResourcePath(name)

	if err != nil {
		log.Fatal(err)
	}

	return path
}

func dashboardAppHandler(base http.Handler, csrfMiddleware func(http.Handler) http.Handler, cfg *config.I18nConfig) http.Handler {
	handler := base

	if cfg != nil {
		handler = corei18n.Middleware(cfg, handler)
	}

	if csrfMiddleware != nil {
		handler = csrfMiddleware(handler)
	}

	handler = middleware.Precognition()(handler)

	return i.Middleware(handler)
}
