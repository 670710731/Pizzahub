package auth

import (
	"strings"
	"testing"
	"time"
)

func TestTokenIssueAndVerify(t *testing.T) {
	service, err := NewTokenService(strings.Repeat("x", 32), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	token, expiresAt, err := service.Issue(42, "customer")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := service.Verify(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 42 || claims.Role != "customer" || claims.Expiry != expiresAt.Unix() {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestTokenRejectsTampering(t *testing.T) {
	service, err := NewTokenService(strings.Repeat("x", 32), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := service.Issue(42, "customer")
	if err != nil {
		t.Fatal(err)
	}
	token += "tampered"
	if _, err := service.Verify(token); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}

func TestTokenServiceRequiresStrongSecret(t *testing.T) {
	if _, err := NewTokenService("short", time.Hour); err == nil {
		t.Fatal("expected short secret to be rejected")
	}
}
