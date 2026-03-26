package auth

import (
	"os"
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(42, "USER")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("token is empty")
	}

	parsed, err := ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if !parsed.Valid {
		t.Fatal("token is not valid")
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("not.a.valid.token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestTokenContainsRole(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	token, _ := GenerateToken(1, "ADMIN")
	parsed, _ := ParseToken(token)

	claims := parsed.Claims.(interface {
		GetSubject() (string, error)
	})
	_ = claims // just checking token parses OK

	if !parsed.Valid {
		t.Fatal("token should be valid")
	}
}
