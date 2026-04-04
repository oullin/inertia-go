package httpx

import (
	"fmt"
	"net/url"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	t.Run("flattens nested map", func(t *testing.T) {
		data := map[string]any{
			"user": map[string]any{
				"name": "Alice",
				"address": map[string]any{
					"city": "Berlin",
				},
			},
			"active": true,
		}

		out := make(url.Values)

		if err := flattenJSON("", data, out, 0); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tests := map[string]string{
			"user.name":         "Alice",
			"user.address.city": "Berlin",
			"active":            "1",
		}

		for key, want := range tests {
			if got := out.Get(key); got != want {
				t.Errorf("key %q = %q, want %q", key, got, want)
			}
		}
	})

	t.Run("returns error when depth exceeds limit", func(t *testing.T) {
		data := buildDeepMap(maxJSONDepth + 2)

		out := make(url.Values)
		err := flattenJSON("", data, out, 0)

		if err == nil {
			t.Fatal("expected error for deeply nested JSON, got nil")
		}

		want := fmt.Sprintf("JSON nesting exceeds maximum depth of %d", maxJSONDepth)

		if err.Error() != want {
			t.Errorf("error = %q, want %q", err.Error(), want)
		}
	})

	t.Run("succeeds at exactly max depth", func(t *testing.T) {
		data := buildDeepMap(maxJSONDepth)

		out := make(url.Values)

		if err := flattenJSON("", data, out, 0); err != nil {
			t.Fatalf("unexpected error at max depth: %v", err)
		}
	})
}

func buildDeepMap(depth int) map[string]any {
	current := map[string]any{"leaf": "value"}

	for i := depth - 1; i > 0; i-- {
		current = map[string]any{fmt.Sprintf("level%d", i): current}
	}

	return current
}
