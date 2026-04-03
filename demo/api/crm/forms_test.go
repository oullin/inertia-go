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
				"first_name": "The first name field is required.",
				"last_name":  "The last name field is required.",
				"email":      "The email field is required.",
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
				"email": "The email field must be a valid email address.",
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
				t.Fatalf("len(validate()) = %d, want %d; got=%v", len(got), len(tt.want), got)
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

	if got := (organizationForm{}).validate()["name"]; got != "The name field is required." {
		t.Fatalf("validate()[name] = %q", got)
	}

	if got := (organizationForm{Name: "Acme"}).validate(); len(got) != 0 {
		t.Fatalf("validate() = %#v, want no errors", got)
	}
}
