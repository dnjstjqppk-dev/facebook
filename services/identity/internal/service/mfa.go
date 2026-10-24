package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

var (
	ErrInvalidMFACode = errors.New("invalid or expired multi-factor authentication code")
)

type MFAService struct{}

func NewMFAService() *MFAService {
	return &MFAService{}
}

// GenerateSecret creates a 20-byte random Base32 encoded secret key
func (s *MFAService) GenerateSecret() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}

// GenerateBackupCodes produces 8 unique cryptographically random recovery codes
func (s *MFAService) GenerateBackupCodes(count int) ([]string, error) {
	if count <= 0 {
		count = 8
	}
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 5)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		codes[i] = strings.ToUpper(hex.EncodeToString(b))
	}
	return codes, nil
}

// GenerateTOTPUri produces standard otpauth:// URL for Google Authenticator / Authy
func (s *MFAService) GenerateTOTPUri(issuer, accountName, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		issuer, accountName, secret, issuer)
}

// ValidateTOTP checks code against current time, plus or minus 1 step tolerance (window of ±30s)
func (s *MFAService) ValidateTOTP(secret, passcode string, t time.Time) bool {
	secret = strings.ToUpper(strings.TrimSpace(secret))
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return false
	}

	currentTimeStep := t.Unix() / 30
	// Test current, -1 step, and +1 step
	for delta := int64(-1); delta <= 1; delta++ {
		step := currentTimeStep + delta
		expectedCode := generateCode(key, step)
		if expectedCode == passcode {
			return true
		}
	}
	return false
}

func generateCode(key []byte, step int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(step))

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	code := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	code = code % uint32(math.Pow10(6))

	return fmt.Sprintf("%06d", code)
}
