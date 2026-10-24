package domain

import (
	"time"
)

type UserSession struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	RefreshTokenHash  string    `json:"-"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	DeviceID          string    `json:"device_id"`
	DeviceName        string    `json:"device_name"`
	IsRevoked         bool      `json:"is_revoked"`
	RevocationReason  string    `json:"revocation_reason,omitempty"`
	ExpiresAt         time.Time `json:"expires_at"`
	CreatedAt         time.Time `json:"created_at"`
	LastActiveAt      time.Time `json:"last_active_at"`
}

