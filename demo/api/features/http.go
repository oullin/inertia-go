package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oullin/inertia-go/core/httpx"
)

func (a app) useHttpHandler(w http.ResponseWriter, r *http.Request) {
	a.deps.Render(w, r, "Features/Http/UseHttp", httpx.Props{})
}

func (a app) useHttpApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	httpx.ParseForm(r)
	name := strings.TrimSpace(r.FormValue("name"))

	if name == "" {
		name = "World"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"greeting":  fmt.Sprintf("Hello, %s!", name),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
