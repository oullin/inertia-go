package validation

import (
	"testing"
)

type testForm struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=255"`
	Age   int    `json:"age" validate:"omitempty,gte=0,lte=150"`
}

//nolint:goconst

// Should use json tag "name" not struct field "Name"

type formTagForm struct {
	Title string `form:"title" validate:"required"`
}

func TestValidate_PassesWithValidData(t *testing.T) {
	t.Parallel()

	form := testForm{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "555-1234",
		Age:   30,
	}

	errors := Validate(form)

	if errors != nil {
		t.Errorf("expected no errors, got %v", errors)
	}
}

func TestValidate_RequiredFields(t *testing.T) {
	t.Parallel()

	form := testForm{}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected error for 'name'")
	}

	if _, ok := errors["email"]; !ok {
		t.Error("expected error for 'email'")
	}

	if _, ok := errors["phone"]; ok {
		t.Error("did not expect error for optional 'phone'")
	}
}

func TestValidate_InvalidEmail(t *testing.T) {
	t.Parallel()

	form := testForm{
		Name:  "John",
		Email: "not-an-email",
	}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["email"].(string)

	if !ok {
		t.Fatalf("expected string error for email, got %T", errors["email"])
	}

	if msg != "The email field must be a valid email address." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_MaxLength(t *testing.T) {
	t.Parallel()

	long := make([]byte, 256)

	for i := range long {
		long[i] = 'a'
	}

	form := testForm{
		Name:  string(long),
		Email: "valid@example.com",
	}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected max-length error for 'name'")
	}
}

func TestValidate_UsesJSONTagForFieldName(t *testing.T) {
	t.Parallel()

	form := testForm{}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected field key to be json tag 'name'")
	}

	if _, ok := errors["Name"]; ok {
		t.Error("field key should not be struct field 'Name'")
	}
}

func TestValidate_FallsBackToFormTag(t *testing.T) {
	t.Parallel()

	form := formTagForm{}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors")
	}

	if _, ok := errors["title"]; !ok {
		t.Error("expected field key to be form tag 'title'")
	}
}

func TestValidate_MaxLength_AtBoundary(t *testing.T) {
	t.Parallel()

	boundary := make([]byte, 255)

	for i := range boundary {
		boundary[i] = 'a'
	}

	form := testForm{
		Name:  string(boundary),
		Email: "valid@example.com",
	}

	errors := Validate(form)

	if errors != nil {
		t.Errorf("expected no errors for 255-char name, got %v", errors)
	}
}

func TestValidate_MaxLength_Phone(t *testing.T) {
	t.Parallel()

	long := make([]byte, 256)

	for i := range long {
		long[i] = '5'
	}

	form := testForm{
		Name:  "John",
		Email: "john@example.com",
		Phone: string(long),
	}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors for 256-char phone, got nil")
	}

	msg, ok := errors["phone"].(string)

	if !ok {
		t.Fatalf("expected string error for phone, got %T", errors["phone"])
	}

	if msg != "The phone field must not exceed 255 characters." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestHumanize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"first_name", "first name"},
		{"firstName", "first name"},
		{"FirstName", "first name"},
		{"email", "email"},
		{"already lower", "already lower"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := humanize(tt.input)

			if got != tt.want {
				t.Errorf("humanize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidate_ReturnsNilForValidStruct(t *testing.T) {
	t.Parallel()

	form := testForm{
		Name:  "Valid",
		Email: "valid@example.com",
	}

	errors := Validate(form)

	if errors != nil {
		t.Errorf("expected nil, got %v", errors)
	}
}
