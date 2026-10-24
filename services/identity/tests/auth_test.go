package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/domain"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/service"
)

func setupTestAuthService() (*service.AuthService, *service.JWTManager, repository.UserRepository, repository.SessionRepository) {
	repo := repository.NewInMemoryUserRepository()
	sessionRepo := repository.NewInMemorySessionRepository()
	jwtManager := service.NewJWTManager("test-jwt-secret-must-be-very-secure-32b", time.Hour, 24*time.Hour)
	authSvc := service.NewAuthService(repo, jwtManager)
	return authSvc, jwtManager, repo, sessionRepo
}

func TestRegistrationAndLogin(t *testing.T) {
	authSvc, _, _, _ := setupTestAuthService()
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
	if res.User.Role != domain.RoleUserStandard {
		t.Errorf("expected ROLE_USER_STANDARD, got: %s", res.User.Role)
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
	authSvc, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	req := &domain.RegisterRequest{
		Email:       "bob@example.com",
		Password:    "Password123!",
		FirstName:   "Bob",
		LastName:    "Jones",
		DateOfBirth: "1995-01-01",
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

func TestPasswordStrengthEnforcement(t *testing.T) {
	cases := []struct {
		password string
		valid    bool
	}{
		{"short1!", false},
		{"alllowercase123!", false},
		{"ALLUPPERCASE123!", false},
		{"NoDigitsHere!@#", false},
		{"NoSpecial123456", false},
		{"ValidPassw0rd!", true},
	}

	for _, tc := range cases {
		err := domain.ValidatePasswordStrength(tc.password)
		if tc.valid && err != nil {
			t.Errorf("expected %s to be valid, got: %v", tc.password, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("expected %s to be invalid, got nil error", tc.password)
		}
	}
}

func TestAgePolicyEnforcement(t *testing.T) {
	authSvc, _, _, _ := setupTestAuthService()
	ctx := context.Background()

	// 1. Under 14: should be rejected
	childReq := &domain.RegisterRequest{
		Email:       "child@example.com",
		Password:    "P@ssword123!",
		FirstName:   "Child",
		LastName:    "Doe",
		DateOfBirth: "2018-01-01",
	}
	_, err := authSvc.Register(ctx, childReq)
	if err != domain.ErrAgeUnderMinimum {
		t.Fatalf("expected ErrAgeUnderMinimum, got: %v", err)
	}

	// 2. Minor (15 years old): should be registered with ROLE_USER_MINOR
	minorReq := &domain.RegisterRequest{
		Email:       "minor@example.com",
		Password:    "P@ssword123!",
		FirstName:   "Teen",
		LastName:    "Doe",
		DateOfBirth: "2010-06-15",
	}
	resMinor, err := authSvc.Register(ctx, minorReq)
	if err != nil {
		t.Fatalf("expected minor registration to succeed, got: %v", err)
	}
	if resMinor.User.Role != domain.RoleUserMinor {
		t.Errorf("expected ROLE_USER_MINOR, got: %s", resMinor.User.Role)
	}
}

func TestMFATOTPLifecycle(t *testing.T) {
	mfaSvc := service.NewMFAService()

	secret, err := mfaSvc.GenerateSecret()
	if err != nil {
		t.Fatalf("failed to generate secret: %v", err)
	}
	if len(secret) == 0 {
		t.Fatal("secret should not be empty")
	}

	uri := mfaSvc.GenerateTOTPUri("GlobalSocial", "user@example.com", secret)
	if len(uri) == 0 {
		t.Fatal("URI should not be empty")
	}

	backupCodes, err := mfaSvc.GenerateBackupCodes(8)
	if err != nil || len(backupCodes) != 8 {
		t.Fatalf("expected 8 backup codes, got %d, err: %v", len(backupCodes), err)
	}

	now := time.Now().UTC()
	if mfaSvc.ValidateTOTP(secret, "000000", now) {
		t.Errorf("random code should not validate against new secret")
	}
}

func TestSessionManagementAndRevocation(t *testing.T) {
	sessionRepo := repository.NewInMemorySessionRepository()
	ctx := context.Background()

	session1 := &domain.UserSession{
		ID:         "sess-1",
		UserID:     "user-100",
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		DeviceName: "Chrome Desktop",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		CreatedAt:  time.Now(),
	}
	session2 := &domain.UserSession{
		ID:         "sess-2",
		UserID:     "user-100",
		IPAddress:  "10.0.0.1",
		UserAgent:  "SocialApp/1.0",
		DeviceName: "iPhone 15",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		CreatedAt:  time.Now(),
	}

	_ = sessionRepo.SaveSession(ctx, session1)
	_ = sessionRepo.SaveSession(ctx, session2)

	list, err := sessionRepo.ListUserSessions(ctx, "user-100")
	if err != nil || len(list) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(list))
	}

	err = sessionRepo.RevokeSession(ctx, "sess-1", "User logged out")
	if err != nil {
		t.Fatalf("failed to revoke session: %v", err)
	}

	s1, _ := sessionRepo.GetSessionByID(ctx, "sess-1")
	if !s1.IsRevoked {
		t.Errorf("session1 should be marked as revoked")
	}

	_ = sessionRepo.RevokeAllUserSessionsExcept(ctx, "user-100", "sess-2", "Security alert")
	s2, _ := sessionRepo.GetSessionByID(ctx, "sess-2")
	if s2.IsRevoked {
		t.Errorf("session2 should remain active")
	}
}

func TestPasswordResetFlow(t *testing.T) {
	authSvc, _, userRepo, sessionRepo := setupTestAuthService()
	lifecycleSvc := service.NewLifecycleService(userRepo, sessionRepo)
	ctx := context.Background()

	regReq := &domain.RegisterRequest{
		Email:       "resetuser@example.com",
		Password:    "OldP@ssword123!",
		FirstName:   "Reset",
		LastName:    "User",
		DateOfBirth: "1990-01-01",
	}
	_, err := authSvc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	token, err := lifecycleSvc.RequestPasswordReset(ctx, "resetuser@example.com")
	if err != nil || token == "" {
		t.Fatalf("failed to request password reset: %v", err)
	}

	err = lifecycleSvc.ResetPassword(ctx, token, "NewP@ssword2026!")
	if err != nil {
		t.Fatalf("password reset failed: %v", err)
	}

	loginRes, err := authSvc.Login(ctx, &domain.LoginRequest{
		Email:    "resetuser@example.com",
		Password: "NewP@ssword2026!",
	})
	if err != nil {
		t.Fatalf("login with new password failed: %v", err)
	}
	if loginRes.User.Email != "resetuser@example.com" {
		t.Errorf("expected email match, got: %s", loginRes.User.Email)
	}

	err = lifecycleSvc.ResetPassword(ctx, token, "AnotherP@ssword123!")
	if err != service.ErrResetTokenExpired {
		t.Fatalf("expected ErrResetTokenExpired on token reuse, got: %v", err)
	}
}

func TestAccountDeactivationAndGDPRExport(t *testing.T) {
	authSvc, _, userRepo, sessionRepo := setupTestAuthService()
	lifecycleSvc := service.NewLifecycleService(userRepo, sessionRepo)
	ctx := context.Background()

	regReq := &domain.RegisterRequest{
		Email:       "gdpr@example.com",
		Password:    "GdprP@ssword123!",
		FirstName:   "GDPR",
		LastName:    "Subject",
		DateOfBirth: "1992-03-20",
	}
	regRes, _ := authSvc.Register(ctx, regReq)
	userID := regRes.User.ID

	exportBytes, err := lifecycleSvc.ExportUserData(ctx, userID)
	if err != nil {
		t.Fatalf("failed to export user data: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(exportBytes, &parsed); err != nil {
		t.Fatalf("exported data is not valid JSON: %v", err)
	}
	if parsed["export_metadata"] == nil || parsed["user"] == nil {
		t.Errorf("export missing critical metadata or user block")
	}

	err = lifecycleSvc.DeactivateAccount(ctx, userID)
	if err != nil {
		t.Fatalf("deactivation failed: %v", err)
	}

	user, _, _ := userRepo.GetUserByID(ctx, userID)
	if user.Status != "DEACTIVATED" {
		t.Errorf("expected status DEACTIVATED, got: %s", user.Status)
	}
}
