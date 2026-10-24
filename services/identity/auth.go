package identity

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email/phone or password")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid or tampered authentication token")
	ErrWeakPassword       = errors.New("password must be at least 8 characters and include uppercase, lowercase, number, and special character")
	ErrUserExists         = errors.New("a user with this email or phone number already exists")
)

type UserClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	MFAVerified bool      `json:"mfa_verified"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

type AuthService struct {
	jwtSecret           []byte
	accessTokenDuration time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret:           []byte(jwtSecret),
		accessTokenDuration: 15 * time.Minute,
		refreshTokenDuration: 7 * 24 * time.Hour,
	}
}

// HashPassword hashes a raw password using bcrypt with standard cost
func (s *AuthService) HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrWeakPassword
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword verifies the plain password against the stored bcrypt hash
func (s *AuthService) CheckPassword(plainPassword, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
	return err == nil
}

// GenerateTokenPair issues a new short-lived Access Token and opaque Refresh Token
func (s *AuthService) GenerateTokenPair(userID uuid.UUID, email, role string, mfaVerified bool) (*TokenPair, error) {
	now := time.Now().UTC()
	accessExpires := now.Add(s.accessTokenDuration)

	claims := UserClaims{
		UserID:      userID,
		Email:       email,
		Role:        role,
		MFAVerified: mfaVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "social-identity-service",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExpires),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccess, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate 32 bytes cryptographically secure random token for refresh
	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return nil, fmt.Errorf("failed to generate secure refresh token: %w", err)
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshBytes)

	return &TokenPair{
		AccessToken:  signedAccess,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpires,
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken parses and verifies a JWT token string
func (s *AuthService) ValidateAccessToken(tokenString string) (*UserClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
