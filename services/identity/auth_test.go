package identity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuthService_PasswordHashing(t *testing.T) {
	svc := NewAuthService("super-secret-test-jwt-key-32-chars-long!")

	password := "SecureP@ssw0rd2026!"
	hash, err := svc.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hash == password {
		t.Errorf("hash should not match plain password")
	}

	if !svc.CheckPassword(password, hash) {
		t.Errorf("CheckPassword failed on valid password")
	}

	if svc.CheckPassword("WrongPassword123!", hash) {
		t.Errorf("CheckPassword succeeded on incorrect password")
	}

	// Test short password rejection
	_, err = svc.HashPassword("short")
	if err == nil {
		t.Errorf("expected error for short password, got nil")
	}
}

func TestAuthService_TokenLifecycle(t *testing.T) {
	secret := "super-secret-test-jwt-key-32-chars-long!"
	svc := NewAuthService(secret)

	userID := uuid.New()
	email := "alex@example.com"
	role := "ROLE_USER"

	pair, err := svc.GenerateTokenPair(userID, email, role, true)
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Errorf("tokens must not be empty")
	}

	claims, err := svc.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("failed to validate generated access token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
	if !claims.MFAVerified {
		t.Errorf("expected MFAVerified to be true")
	}
}

func TestAuthService_InvalidToken(t *testing.T) {
	svc := NewAuthService("key-1")
	svcOther := NewAuthService("key-2")

	pair, _ := svc.GenerateTokenPair(uuid.New(), "test@test.com", "ROLE_USER", false)

	// Validating with different secret should fail
	_, err := svcOther.ValidateAccessToken(pair.AccessToken)
	if err == nil {
		t.Errorf("expected validation to fail with mismatching secret key")
	}
}

func TestAuthService_ExpiredToken(t *testing.T) {
	svc := &AuthService{
		jwtSecret:           []byte("secret"),
		accessTokenDuration: -1 * time.Minute, // already expired
	}

	pair, err := svc.GenerateTokenPair(uuid.New(), "expired@test.com", "ROLE_USER", false)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	_, err = svc.ValidateAccessToken(pair.AccessToken)
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}
