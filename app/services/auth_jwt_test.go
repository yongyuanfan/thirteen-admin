package services

import (
	"testing"
	"time"

	"thirteen-admin/internal/config"
)

func TestIssueAndParseToken(t *testing.T) {
	m := NewJWTManager(config.JWT{
		Secret:    "test-secret",
		ExpiresIn: 60,
	})

	token, err := m.IssueToken(1, "admin")
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	if token == "" {
		t.Fatal("empty token")
	}

	access, err := m.ParseAccess(token)
	if err != nil {
		t.Fatalf("ParseAccess: %v", err)
	}
	if access.UserID != "1" || access.UserName != "admin" || access.Typ != tokenTypeAccess {
		t.Fatalf("unexpected access claims: %+v", access)
	}
}

func TestParseAccessExpired(t *testing.T) {
	m := &JWTManager{
		secret:    []byte("test-secret"),
		accessTTL: time.Millisecond,
	}
	token, err := m.IssueToken(2, "User")
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	_, err = m.ParseAccess(token)
	if err != ErrExpiredToken {
		t.Fatalf("want ErrExpiredToken, got %v", err)
	}
}
