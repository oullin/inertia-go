package database

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserHashesPassword(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if _, err := CreateUser(db, "Demo User", "test@example.com", "password", nil); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	user, err := FindUserByEmail(db, "test@example.com")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	if user == nil {
		t.Fatal("FindUserByEmail() returned nil user")
	}

	if user.PasswordHash == "password" {
		t.Fatal("CreateUser() stored plaintext password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password")); err != nil {
		t.Fatalf("stored password hash did not verify: %v", err)
	}
}

func TestCreateUserNormalisesEmail(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")

	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if _, err := CreateUser(db, "Demo User", "  User@Example.COM  ", "password", nil); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	user, err := FindUserByEmail(db, "USER@EXAMPLE.COM")

	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}

	if user == nil {
		t.Fatal("FindUserByEmail() returned nil for differently-cased email")
	}

	if user.Email != "user@example.com" {
		t.Fatalf("stored email = %q, want %q", user.Email, "user@example.com")
	}
}
