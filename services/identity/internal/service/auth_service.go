package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/domain"
	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmptyFields        = errors.New("required fields cannot be empty")
)

type AuthService struct {
	repo       repository.UserRepository
	jwtManager *JWTManager
}

func NewAuthService(repo repository.UserRepository, jwtManager *JWTManager) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, ErrEmptyFields
	}

	// Validate password strength
	if err := domain.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Calculate age and assign correct RBAC role
	role := domain.RoleUserStandard
	if req.DateOfBirth != "" {
		age, err := domain.CalculateAge(req.DateOfBirth, time.Now().UTC())
		if err != nil {
			return nil, err
		}
		if age < 14 {
			return nil, domain.ErrAgeUnderMinimum
		}
		if age < 18 {
			role = domain.RoleUserMinor
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := generateUUID()
	now := time.Now().UTC()

	user := &domain.User{
		ID:            userID,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		Status:        domain.StatusActive,
		EmailVerified: false,
		Role:          role,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	username := req.Username
	if username == "" {
		username = strings.Split(email, "@")[0] + "_" + userID[:6]
	}

	profile := &domain.Profile{
		UserID:             userID,
		Username:           username,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		DateOfBirth:        req.DateOfBirth,
		Gender:             req.Gender,
		PrivacyLevel:       "PUBLIC",
		IsProfessionalMode: false,
		CreatedAt:          now,
	}

	if err := s.repo.CreateUser(ctx, user, profile); err != nil {
		return nil, err
	}

	tokens, err := s.jwtManager.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:    user,
		Profile: profile,
		Tokens:  tokens,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		return nil, ErrEmptyFields
	}

	user, profile, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.jwtManager.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:    user,
		Profile: profile,
		Tokens:  tokens,
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*domain.User, *domain.Profile, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *AuthService) VerifyAccessToken(tokenStr string) (*JWTClaims, error) {
	return s.jwtManager.VerifyToken(tokenStr)
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC4122
	return hex.EncodeToString(b[0:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:16])
}
