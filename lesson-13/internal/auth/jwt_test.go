package auth

import (
	"strings"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	now := time.Now()
	token, err := Sign([]byte("secret"), Claims{Subject: "ada", Role: RolePlayer, ExpiresAt: now.Add(time.Hour).Unix()})
	if err != nil {
		t.Fatal(err)
	}
	claims, err := Verify([]byte("secret"), token, now)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != "ada" || claims.Role != RolePlayer {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestRejectsExpiredWrongSecretAndAlgNone(t *testing.T) {
	now := time.Now()
	token, err := Sign([]byte("secret"), Claims{Subject: "ada", Role: RolePlayer, ExpiresAt: now.Add(-time.Minute).Unix()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify([]byte("secret"), token, now); err != ErrExpired {
		t.Fatalf("expired error = %v", err)
	}
	token, err = Sign([]byte("secret"), Claims{Subject: "ada", Role: RolePlayer, ExpiresAt: now.Add(time.Hour).Unix()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Verify([]byte("other"), token, now); err != ErrSignature {
		t.Fatalf("wrong secret error = %v", err)
	}
	none := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJhZGEiLCJyb2xlIjoibW9kZXJhdG9yIiwiZXhwIjo5OTk5OTk5OTk5fQ."
	if _, err := Verify([]byte("secret"), none, now); err == nil {
		t.Fatal("alg=none accepted")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatal("unexpected token")
	}
}
