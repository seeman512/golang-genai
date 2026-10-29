// Package auth implements the small JWT example used by Lesson 13.
//
// It intentionally supports only HS256 and uses the standard library so the
// authentication flow stays visible. Production applications should use a
// maintained JWT library and an external identity provider.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	RolePlayer    = "player"
	RoleModerator = "moderator"
	// DefaultDemoSecret is only for local demonstrations.
	DefaultDemoSecret = "mystic-square-demo-secret-change-me"
)

var (
	ErrMalformed = errors.New("malformed token")
	ErrSignature = errors.New("bad signature")
	ErrExpired   = errors.New("token expired")
)

type Claims struct {
	Subject   string `json:"sub"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
}

var encoding = base64.RawURLEncoding

// Sign creates an HS256 JWT for the supplied claims.
func Sign(secret []byte, claims Claims) (string, error) {
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	input := encoding.EncodeToString(header) + "." + encoding.EncodeToString(payload)
	return input + "." + encoding.EncodeToString(mac(secret, input)), nil
}

// Verify validates the token signature, algorithm, and expiry time.
func Verify(secret []byte, token string, now time.Time) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrMalformed
	}
	headerBytes, err := encoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrMalformed
	}
	var header struct {
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Alg != "HS256" {
		return nil, ErrMalformed
	}
	signature, err := encoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(signature, mac(secret, parts[0]+"."+parts[1])) {
		return nil, ErrSignature
	}
	payload, err := encoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrMalformed
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrMalformed
	}
	if claims.ExpiresAt == 0 || now.Unix() >= claims.ExpiresAt {
		return nil, ErrExpired
	}
	if claims.Subject == "" || claims.Role == "" {
		return nil, ErrMalformed
	}
	return &claims, nil
}

func mac(secret []byte, input string) []byte {
	h := hmac.New(sha256.New, secret)
	_, _ = h.Write([]byte(input))
	return h.Sum(nil)
}
