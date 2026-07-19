package services

import (
	"testing"
	"time"

	"thirteen-admin/internal/config"
)

func boolPtr(v bool) *bool { return &v }

func TestIssueAndParsePair(t *testing.T) {
	m := NewJWTManager(config.JWT{
		Secret:           "test-secret",
		ExpiresIn:        60,
		RefreshExpiresIn: 3600,
		AutoRefreshToken: boolPtr(true),
		AutoRevokeToken:  boolPtr(true),
	})

	pair, err := m.IssuePair(1, "admin", "", time.Time{})
	if err != nil {
		t.Fatalf("IssuePair: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" || pair.Sid == "" {
		t.Fatal("empty token pair")
	}
	if pair.ExpiresIn != 60 {
		t.Fatalf("expiresIn=%d", pair.ExpiresIn)
	}

	access, err := m.ParseAccess(pair.AccessToken)
	if err != nil {
		t.Fatalf("ParseAccess: %v", err)
	}
	if access.UserID != "1" || access.UserName != "admin" || access.Typ != tokenTypeAccess || access.Sid != pair.Sid {
		t.Fatalf("unexpected access claims: %+v", access)
	}

	refresh, err := m.ParseRefresh(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ParseRefresh: %v", err)
	}
	if refresh.Typ != tokenTypeRefresh || refresh.Sid != pair.Sid {
		t.Fatalf("unexpected refresh claims: %+v", refresh)
	}

	if _, err := m.ParseAccess(pair.RefreshToken); err != ErrInvalidToken {
		t.Fatalf("refresh as access want ErrInvalidToken, got %v", err)
	}
	if _, err := m.ParseRefresh(pair.AccessToken); err != ErrInvalidToken {
		t.Fatalf("access as refresh want ErrInvalidToken, got %v", err)
	}
}

func TestParseAccessExpired(t *testing.T) {
	m := &JWTManager{
		secret:           []byte("test-secret"),
		accessTTL:        time.Millisecond,
		refreshTTL:       time.Hour,
		autoRefreshToken: true,
		autoRevokeToken:  true,
	}
	pair, err := m.IssuePair(2, "User", "", time.Time{})
	if err != nil {
		t.Fatalf("IssuePair: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	_, err = m.ParseAccess(pair.AccessToken)
	if err != ErrExpiredToken {
		t.Fatalf("want ErrExpiredToken, got %v", err)
	}
}

func TestAbsoluteRefreshKeepsExpiry(t *testing.T) {
	m := NewJWTManager(config.JWT{
		Secret:           "test-secret",
		ExpiresIn:        60,
		RefreshExpiresIn: 3600,
		AutoRefreshToken: boolPtr(false),
		AutoRevokeToken:  boolPtr(false),
	})
	pair, err := m.IssuePair(1, "admin", "", time.Time{})
	if err != nil {
		t.Fatalf("IssuePair: %v", err)
	}
	old, err := m.ParseRefresh(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ParseRefresh: %v", err)
	}
	next, err := m.IssuePair(1, "admin", pair.Sid, old.ExpiresAt.Time)
	if err != nil {
		t.Fatalf("IssuePair rotate: %v", err)
	}
	neu, err := m.ParseRefresh(next.RefreshToken)
	if err != nil {
		t.Fatalf("ParseRefresh new: %v", err)
	}
	if !neu.ExpiresAt.Time.Equal(old.ExpiresAt.Time) {
		t.Fatalf("refresh exp changed: old=%v new=%v", old.ExpiresAt.Time, neu.ExpiresAt.Time)
	}
	if next.Sid != pair.Sid {
		t.Fatalf("sid changed")
	}
}

func TestSlidingRefreshExtendsExpiry(t *testing.T) {
	m := NewJWTManager(config.JWT{
		Secret:           "test-secret",
		ExpiresIn:        60,
		RefreshExpiresIn: 2,
		AutoRefreshToken: boolPtr(true),
		AutoRevokeToken:  boolPtr(true),
	})
	pair, err := m.IssuePair(1, "admin", "", time.Time{})
	if err != nil {
		t.Fatalf("IssuePair: %v", err)
	}
	old, err := m.ParseRefresh(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ParseRefresh: %v", err)
	}
	time.Sleep(1100 * time.Millisecond)
	next, err := m.IssuePair(1, "admin", pair.Sid, time.Time{})
	if err != nil {
		t.Fatalf("IssuePair slide: %v", err)
	}
	neu, err := m.ParseRefresh(next.RefreshToken)
	if err != nil {
		t.Fatalf("ParseRefresh new: %v", err)
	}
	if !neu.ExpiresAt.Time.After(old.ExpiresAt.Time) {
		t.Fatalf("refresh exp not extended: old=%v new=%v", old.ExpiresAt.Time, neu.ExpiresAt.Time)
	}
}

func TestDefaultSwitchesTrue(t *testing.T) {
	m := NewJWTManager(config.JWT{Secret: "x", ExpiresIn: 1, RefreshExpiresIn: 2})
	if !m.AutoRefreshToken() || !m.AutoRevokeToken() {
		t.Fatal("switches should default true")
	}
}
