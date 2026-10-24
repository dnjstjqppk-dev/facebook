package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/domain"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrResetTokenExpired = errors.New("password reset token is expired or invalid")
	ErrAccountDeactivated = errors.New("account is deactivated; please recover within the 30-day grace period")
	ErrAccountPurged      = errors.New("account has been permanently deleted")
)

type ResetToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Used      bool
}

type LifecycleService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	mu          sync.RWMutex
	resetTokens map[string]*ResetToken // token -> details
}

func NewLifecycleService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *LifecycleService {
	return &LifecycleService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		resetTokens: make(map[string]*ResetToken),
	}
}

// RequestPasswordReset generates a single-use 15-minute token
func (s *LifecycleService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	user, _, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Silent success pattern for security / anti-enumeration
		return "", nil
	}

	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(b)

	s.mu.Lock()
	s.resetTokens[tokenStr] = &ResetToken{
		Token:     tokenStr,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
		Used:      false,
	}
	s.mu.Unlock()

	return tokenStr, nil
}

// ResetPassword verifies the token, updates password hash, and revokes all active sessions
func (s *LifecycleService) ResetPassword(ctx context.Context, tokenStr, newPassword string) error {
	if err := domain.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	s.mu.Lock()
	record, exists := s.resetTokens[tokenStr]
	if !exists || record.Used || time.Now().UTC().After(record.ExpiresAt) {
		s.mu.Unlock()
		return ErrResetTokenExpired
	}
	record.Used = true
	userID := record.UserID
	s.mu.Unlock()

	user, _, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now().UTC()

	// Revoke all existing sessions
	if s.sessionRepo != nil {
		_ = s.sessionRepo.RevokeAllUserSessionsExcept(ctx, userID, "", "Password reset initiated")
	}

	return nil
}

// DeactivateAccount initiates the 30-day grace period soft deletion
func (s *LifecycleService) DeactivateAccount(ctx context.Context, userID string) error {
	user, _, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Status = domain.UserStatus("DEACTIVATED")
	user.UpdatedAt = time.Now().UTC()

	// Revoke all sessions immediately
	if s.sessionRepo != nil {
		_ = s.sessionRepo.RevokeAllUserSessionsExcept(ctx, userID, "", "Account deactivation requested")
	}
	return nil
}

// ExportUserData generates the complete GDPR / PIPA JSON compliance archive
func (s *LifecycleService) ExportUserData(ctx context.Context, userID string) ([]byte, error) {
	user, profile, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	exportPayload := map[string]interface{}{
		"export_metadata": map[string]interface{}{
			"version":       "1.0.0",
			"generated_at":  time.Now().UTC().Format(time.RFC3339),
			"format":        "GDPR_PIPA_COMPLIANT_JSON",
		},
		"user": map[string]interface{}{
			"id":             user.ID,
			"email":          user.Email,
			"status":         user.Status,
			"mfa_enabled":    user.MFAEnabled,
			"role":           user.Role,
			"created_at":     user.CreatedAt,
		},
		"profile": profile,
	}

	return json.MarshalIndent(exportPayload, "", "  ")
}

