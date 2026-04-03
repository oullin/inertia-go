package crm

import "testing"

func TestContactFormValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		form contactForm
		want map[string]string
	}{
		{
			name: "missing required fields",
			form: contactForm{},
			want: map[string]string{
				"first_name": "First name is required.",
				"last_name":  "Last name is required.",
				"email":      "A valid email address is required.",
			},
		},
		{
			name: "invalid email",
			form: contactForm{
				FirstName: "Mina",
				LastName:  "Cole",
				Email:     "mina.example.test",
			},
			want: map[string]string{
				"email": "A valid email address is required.",
			},
		},
		{
			name: "valid form",
			form: contactForm{
				FirstName: "Mina",
				LastName:  "Cole",
				Email:     "mina@example.test",
			},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.form.validate()

			if len(got) != len(tt.want) {
				t.Fatalf("len(validate()) = %d, want %d", len(got), len(tt.want))
			}

			for key, want := range tt.want {
				if got[key] != want {
					t.Fatalf("validate()[%q] = %q, want %q", key, got[key], want)
				}
			}
		})
	}
}

func TestOrganizationFormValidate(t *testing.T) {
	t.Parallel()

	if got := (organizationForm{}).validate()["name"]; got != "Organization name is required." {
		t.Fatalf("validate()[name] = %q", got)
	}

	if got := (organizationForm{Name: "Acme"}).validate(); len(got) != 0 {
		t.Fatalf("validate() = %#v, want no errors", got)
	}
}
