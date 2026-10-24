package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrInvalidToken = errors.New("invalid token format or signature")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTClaims struct {
	UserID    string      `json:"sub"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	TokenType string      `json:"type"` // "access" or "refresh"
	IssuedAt  int64       `json:"iat"`
	ExpiresAt int64       `json:"exp"`
}

type JWTManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (m *JWTManager) GenerateTokens(user *domain.User) (*domain.AuthTokens, error) {
	now := time.Now()

	accessClaims := JWTClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "access",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.accessExpiry).Unix(),
	}

	accessToken, err := m.signToken(accessClaims)
	if err != nil {
		return nil, err
	}

	refreshClaims := JWTClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "refresh",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.refreshExpiry).Unix(),
	}

	refreshToken, err := m.signToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (m *JWTManager) signToken(claims JWTClaims) (string, error) {
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	unsignedToken := headerB64 + "." + payloadB64
	sig := m.computeHmac(unsignedToken)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return unsignedToken + "." + sigB64, nil
}

func (m *JWTManager) VerifyToken(tokenStr string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSig := m.computeHmac(unsignedToken)
	actualSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(expectedSig, actualSig) {
		return nil, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func (m *JWTManager) computeHmac(data string) []byte {
	h := hmac.New(sha256.New, m.secretKey)
	h.Write([]byte(data))
	return h.Sum(nil)
}
