package response

import (
	"fmt"
	"html"
	"html/template"
	"net/http"

	ihttp "github.com/oullin/inertia-go/core/http"
)

// WriteJSON writes an Inertia JSON response for XHR visits.
// It sets the Content-Type and X-Inertia headers.
func WriteJSON(w http.ResponseWriter, page *Page, marshaler ihttp.JSONMarshaler) error {
	data, err := marshaler.Marshal(page)

	if err != nil {
		return fmt.Errorf("inertia: marshal page: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(ihttp.HeaderInertia, "true")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)

	return err
}

// WriteHTML executes the root HTML template for initial (non-XHR) page
// visits. The page object is embedded as a JSON string in a data-page
// attribute on the container div.
func WriteHTML(
	w http.ResponseWriter,
	tmpl *template.Template,
	page *Page,
	containerID string,
	marshaler ihttp.JSONMarshaler,
	extraData ihttp.TemplateData,
) error {
	data, err := marshaler.Marshal(page)

	if err != nil {
		return fmt.Errorf("inertia: marshal page: %w", err)
	}

	// Build the container div with the page data escaped for HTML attributes.
	inertiaHTML := fmt.Sprintf(
		`<div id="%s" data-page="%s"></div>`,
		containerID,
		html.EscapeString(string(data)),
	)

	templateData := map[string]any{
		"inertia":     template.HTML(inertiaHTML),
		"inertiaHead": template.HTML(""),
	}

	for k, v := range extraData {
		templateData[k] = v
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return tmpl.Execute(w, templateData)
}
