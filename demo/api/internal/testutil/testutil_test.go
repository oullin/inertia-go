package testutil

import (
	"testing"
)

func TestFindCSRFMetaToken(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "name before content",
			body: `<meta name="csrf-token" content="abc123">`,
			want: "abc123",
		},
		{
			name: "content before name",
			body: `<meta content="xyz789" name="csrf-token">`,
			want: "xyz789",
		},
		{
			name: "extra attributes between",
			body: `<meta name="csrf-token" id="tok" content="mid456">`,
			want: "mid456",
		},
		{
			name: "self-closing tag",
			body: `<meta name="csrf-token" content="sc1" />`,
			want: "sc1",
		},
		{
			name: "embedded in full page",
			body: `<!DOCTYPE html><html><head><meta name="csrf-token" content="page1"></head><body></body></html>`,
			want: "page1",
		},
		{
			name: "content before name with extra attrs",
			body: `<meta class="x" content="rev2" name="csrf-token">`,
			want: "rev2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindCSRFMetaToken(t, tt.body)

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
