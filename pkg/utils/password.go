package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// PasswordUtil validates, hashes and verifies passwords.
type PasswordUtil interface {
	Validate(password string) bool
	Hash(password string) (string, error)
	Verify(password string, hash string) bool
}

type passwordUtil struct{}

// NewPasswordUtil creates a bcrypt-based PasswordUtil.
func NewPasswordUtil() PasswordUtil {
	return &passwordUtil{}
}

// Validate reports whether the password is at least 8 characters and has an upper case letter,
// a lower case letter, a digit and a special character.
func (vp *passwordUtil) Validate(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Hash returns the bcrypt hash of the password.
func (vp *passwordUtil) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify reports whether the password matches the bcrypt hash.
func (vp *passwordUtil) Verify(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
