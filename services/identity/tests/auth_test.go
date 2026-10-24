package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/domain"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/service"
)

func setupTestAuthService() (*service.AuthService, *service.JWTManager) {
	repo := repository.NewInMemoryUserRepository()
	jwtManager := service.NewJWTManager("test-jwt-secret-must-be-very-secure-32b", time.Hour, 24*time.Hour)
	authSvc := service.NewAuthService(repo, jwtManager)
	return authSvc, jwtManager
}

func TestRegistrationAndLogin(t *testing.T) {
	authSvc, _ := setupTestAuthService()
	ctx := context.Background()

	regReq := &domain.RegisterRequest{
		Email:       "alice@example.com",
		Password:    "SuperSecret123!",
		FirstName:   "Alice",
		LastName:    "Smith",
		Username:    "alicesmith",
		DateOfBirth: "1998-05-14",
		Gender:      "FEMALE",
	}

	res, err := authSvc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("expected registration to succeed, got: %v", err)
	}

	if res.User.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got: %s", res.User.Email)
	}
	if res.Profile.Username != "alicesmith" {
		t.Errorf("expected username alicesmith, got: %s", res.Profile.Username)
	}
	if res.Tokens.AccessToken == "" || res.Tokens.RefreshToken == "" {
		t.Errorf("expected valid tokens, got empty")
	}

	// Test Successful Login
	loginReq := &domain.LoginRequest{
		Email:    "alice@example.com",
		Password: "SuperSecret123!",
	}
	loginRes, err := authSvc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("expected login to succeed, got: %v", err)
	}
	if loginRes.User.ID != res.User.ID {
		t.Errorf("expected user ID %s, got %s", res.User.ID, loginRes.User.ID)
	}

	// Test Invalid Password Login
	badLoginReq := &domain.LoginRequest{
		Email:    "alice@example.com",
		Password: "WrongPassword999!",
	}
	_, err = authSvc.Login(ctx, badLoginReq)
	if err != service.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestDuplicateRegistration(t *testing.T) {
	authSvc, _ := setupTestAuthService()
	ctx := context.Background()

	req := &domain.RegisterRequest{
		Email:     "bob@example.com",
		Password:  "Password123!",
		FirstName: "Bob",
		LastName:  "Jones",
	}

	_, err := authSvc.Register(ctx, req)
	if err != nil {
		t.Fatalf("first registration should succeed, got: %v", err)
	}

	_, err = authSvc.Register(ctx, req)
	if err != repository.ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got: %v", err)
	}
}

func TestJWTSignatureVerification(t *testing.T) {
	authSvc, jwtMgr := setupTestAuthService()
	ctx := context.Background()

	regReq := &domain.RegisterRequest{
		Email:     "carol@example.com",
		Password:  "CarolPassword1!",
		FirstName: "Carol",
		LastName:  "Danvers",
	}
	res, err := authSvc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	claims, err := jwtMgr.VerifyToken(res.Tokens.AccessToken)
	if err != nil {
		t.Fatalf("token verification failed: %v", err)
	}
	if claims.Email != "carol@example.com" {
		t.Errorf("expected claim email carol@example.com, got: %s", claims.Email)
	}

	// Tampered Token
	tampered := res.Tokens.AccessToken + "tampered"
	_, err = jwtMgr.VerifyToken(tampered)
	if err != service.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken on tampered token, got: %v", err)
	}
}
