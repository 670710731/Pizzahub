package auth

import "testing"

func TestPasswordHashVerification(t *testing.T) {
	hash, err := HashPassword("correct horse battery")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword(hash, "correct horse battery") {
		t.Fatal("expected password to verify")
	}
	if VerifyPassword(hash, "incorrect password") {
		t.Fatal("expected wrong password not to verify")
	}
}

func TestPasswordLengthValidation(t *testing.T) {
	for _, password := range []string{"short", string(make([]byte, 73))} {
		if _, err := HashPassword(password); err != ErrInvalidPassword {
			t.Fatalf("expected ErrInvalidPassword for length %d, got %v", len(password), err)
		}
	}
}
