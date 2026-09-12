package jwt

import "testing"

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(42, "test-secret", 1)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	claims, err := ParseToken(token, "test-secret")
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("UserID = %d, want 42", claims.UserID)
	}
}

func TestAdminTokenIncludesVerification(t *testing.T) {
	token, err := GenerateAdminToken(7, "test-secret", 1)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseToken(token, "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	if !claims.AdminVerified {
		t.Fatal("admin verification claim was not set")
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	token, err := GenerateToken(42, "test-secret", 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(token, "wrong-secret"); err == nil {
		t.Fatal("ParseToken() accepted a token signed with another secret")
	}
}
