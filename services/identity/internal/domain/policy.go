package domain

import (
	"errors"
	"time"
	"unicode"
)

var (
	ErrPasswordTooShort        = errors.New("password must be at least 8 characters long")
	ErrPasswordMissingUpper    = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLower    = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingDigit    = errors.New("password must contain at least one digit")
	ErrPasswordMissingSpecial  = errors.New("password must contain at least one special character")
	ErrAgeUnderMinimum         = errors.New("user must be at least 14 years old to register")
	ErrInvalidDateFormat       = errors.New("invalid date of birth format, expected YYYY-MM-DD")
)

// ValidatePasswordStrength checks strict password security requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordMissingUpper
	}
	if !hasLower {
		return ErrPasswordMissingLower
	}
	if !hasDigit {
		return ErrPasswordMissingDigit
	}
	if !hasSpecial {
		return ErrPasswordMissingSpecial
	}

	return nil
}

// CalculateAge determines the user's age in years from YYYY-MM-DD
func CalculateAge(dobStr string, now time.Time) (int, error) {
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return 0, ErrInvalidDateFormat
	}

	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		years--
	}
	return years, nil
}
