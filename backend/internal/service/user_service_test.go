package service

import (
	"regexp"
	"testing"
)

func TestNewAdminChallenge(t *testing.T) {
	id, code, err := newAdminChallenge()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 48 {
		t.Fatalf("challenge id length = %d", len(id))
	}
	if !regexp.MustCompile(`^[0-9]{6}$`).MatchString(code) {
		t.Fatalf("invalid code format")
	}
	if adminCodeHash("secret", id, code) == adminCodeHash("secret", id, "000000") && code != "000000" {
		t.Fatal("different codes produced the same hash")
	}
}

func TestMaskEmail(t *testing.T) {
	if got := maskEmail("admin@example.com"); got != "a***@example.com" {
		t.Fatalf("maskEmail() = %q", got)
	}
	if got := maskEmail("invalid"); got != "***" {
		t.Fatalf("invalid email mask = %q", got)
	}
}
