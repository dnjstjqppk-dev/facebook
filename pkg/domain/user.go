package domain

import "time"

type UserStatus string

const (
	StatusActive              UserStatus = "ACTIVE"
	StatusSuspended           UserStatus = "SUSPENDED"
	StatusPendingVerification UserStatus = "PENDING_VERIFICATION"
)

type Role string

const (
	RoleUserStandard   Role = "ROLE_USER_STANDARD"
	RoleUserMinor      Role = "ROLE_USER_MINOR"
	RoleCreatorPro     Role = "ROLE_CREATOR_PRO"
	RolePlatformAdmin  Role = "ROLE_PLATFORM_ADMIN"
	RoleModeratorStaff Role = "ROLE_MODERATOR_SAFETY"
)

type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone,omitempty"`
	PasswordHash  string     `json:"-"`
	Status        UserStatus `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	PhoneVerified bool       `json:"phone_verified"`
	MFAEnabled    bool       `json:"mfa_enabled"`
	Role          Role       `json:"role"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Profile struct {
	UserID             string    `json:"user_id"`
	Username           string    `json:"username"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Bio                string    `json:"bio,omitempty"`
	AvatarURL          string    `json:"avatar_url,omitempty"`
	CoverURL           string    `json:"cover_url,omitempty"`
	DateOfBirth        string    `json:"date_of_birth"`
	Gender             string    `json:"gender"`
	PrivacyLevel       string    `json:"privacy_level"`
	IsProfessionalMode bool      `json:"is_professional_mode"`
	CreatedAt          time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

type AuthResponse struct {
	User    *User       `json:"user"`
	Profile *Profile    `json:"profile"`
	Tokens  *AuthTokens `json:"tokens"`
}
