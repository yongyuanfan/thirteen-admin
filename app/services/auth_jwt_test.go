package services

import (
	"testing"
	"time"

	"thirteen-admin/internal/config"
)

func TestIssueAndParseTokenPair(t *testing.T) {
	m := NewJWTManager(config.JWT{
		Secret:           "test-secret",
		ExpiresIn:        60,
		RefreshExpiresIn: 120,
	})

	pair, err := m.IssueTokenPair(1, "Super")
	if err != nil {
		t.Fatalf("IssueTokenPair: %v", err)
	}
	if pair.Token == "" || pair.RefreshToken == "" {
		t.Fatal("empty tokens")
	}

	access, err := m.ParseAccess(pair.Token)
	if err != nil {
		t.Fatalf("ParseAccess: %v", err)
	}
	if access.UserID != "1" || access.UserName != "Super" || access.Typ != tokenTypeAccess {
		t.Fatalf("unexpected access claims: %+v", access)
	}

	refresh, err := m.ParseRefresh(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ParseRefresh: %v", err)
	}
	if refresh.ID != pair.RefreshJTI {
		t.Fatalf("refresh jti mismatch: %s vs %s", refresh.ID, pair.RefreshJTI)
	}
}

func TestParseAccessExpired(t *testing.T) {
	m := &JWTManager{
		secret:     []byte("test-secret"),
		accessTTL:  time.Millisecond,
		refreshTTL: time.Hour,
	}
	pair, err := m.IssueTokenPair(2, "User")
	if err != nil {
		t.Fatalf("IssueTokenPair: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	_, err = m.ParseAccess(pair.Token)
	if err != ErrExpiredToken {
		t.Fatalf("want ErrExpiredToken, got %v", err)
	}
}
