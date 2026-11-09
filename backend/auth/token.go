package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired access token")

type Claims struct {
	UserID int64  `json:"sub"`
	Role   string `json:"role"`
	Expiry int64  `json:"exp"`
}

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttl time.Duration) (*TokenService, error) {
	if len(secret) < 32 {
		return nil, errors.New("token secret must be at least 32 bytes")
	}
	if ttl <= 0 {
		return nil, errors.New("token lifetime must be positive")
	}
	return &TokenService{secret: []byte(secret), ttl: ttl}, nil
}

func (s *TokenService) Issue(userID int64, role string) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(s.ttl)
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("encode token header: %w", err)
	}
	payload, err := json.Marshal(Claims{UserID: userID, Role: role, Expiry: expiresAt.Unix()})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("encode token claims: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	return encoded + "." + s.sign(encoded), expiresAt, nil
}

func (s *TokenService) Verify(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || !hmac.Equal([]byte(s.sign(parts[0]+"."+parts[1])), []byte(parts[2])) {
		return Claims{}, ErrInvalidToken
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if json.Unmarshal(headerBytes, &header) != nil || header.Algorithm != "HS256" || header.Type != "JWT" {
		return Claims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var claims Claims
	if json.Unmarshal(payload, &claims) != nil || claims.UserID <= 0 || claims.Expiry <= time.Now().Unix() {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}

func (s *TokenService) sign(value string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
