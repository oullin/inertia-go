package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ParseForm parses the request body based on its Content-Type.
// It handles application/json (sent by Inertia.js v3 for non-file forms),
// multipart/form-data, and application/x-www-form-urlencoded.
// After calling ParseForm, r.FormValue() works regardless of content type.
func ParseForm(r *http.Request) error {
	ct := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return parseJSONForm(r)
	case strings.HasPrefix(ct, "multipart/"):
		return r.ParseMultipartForm(32 << 20)
	default:
		return r.ParseForm()
	}
}

func parseJSONForm(r *http.Request) error {
	var raw map[string]any

	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return err
	}

	values := make(url.Values)
	flattenJSON("", raw, values)

	r.PostForm = values
	r.Form = make(url.Values)

	if r.URL != nil {
		qs, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			return err
		}

		for k, v := range qs {
			r.Form[k] = v
		}
	}

	for k, v := range values {
		r.Form[k] = append(r.Form[k], v...)
	}

	return nil
}

func flattenJSON(prefix string, data map[string]any, out url.Values) {
	for key, val := range data {
		fullKey := key

		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := val.(type) {
		case map[string]any:
			flattenJSON(fullKey, v, out)
		case []any:
			for i, item := range v {
				arrKey := fmt.Sprintf("%s[%d]", fullKey, i)

				if nested, ok := item.(map[string]any); ok {
					flattenJSON(arrKey, nested, out)
				} else {
					out.Set(arrKey, toFormValue(item))
				}
			}
		default:
			out.Set(fullKey, toFormValue(val))
		}
	}
}

func toFormValue(v any) string {
	switch val := v.(type) {
	case bool:
		if val {
			return "1"
		}

		return "0"
	case nil:
		return ""
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}

		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
